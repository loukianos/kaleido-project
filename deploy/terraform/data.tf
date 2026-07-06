resource "aws_ecr_repository" "api" {
  name         = "kaleido-project"
  force_delete = true

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "random_password" "db" {
  length  = 32
  special = false
}

resource "aws_db_subnet_group" "db" {
  name       = "${var.name}-db"
  subnet_ids = module.vpc.private_subnets
}

resource "aws_security_group" "db" {
  name_prefix = "${var.name}-db-"
  description = "Postgres access from the EKS nodes"
  vpc_id      = module.vpc.vpc_id

  ingress {
    description     = "Postgres from cluster nodes"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [module.eks.node_security_group_id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_db_instance" "postgres" {
  identifier     = "${var.name}-postgres"
  engine         = "postgres"
  engine_version = "16"
  instance_class = "db.t4g.micro"

  allocated_storage = 20
  storage_type      = "gp3"

  db_name  = var.db_name
  username = var.db_name
  password = random_password.db.result

  db_subnet_group_name   = aws_db_subnet_group.db.name
  vpc_security_group_ids = [aws_security_group.db.id]

  skip_final_snapshot = true
  apply_immediately   = true
}

# The signing-key master key: custodial key material is sealed by KMS, so plaintext keys never depend on an env-var secret in the cluster.
resource "aws_kms_key" "signing_keys" {
  description             = "Envelope encryption for kaleido-project custodial signing keys"
  deletion_window_in_days = 7
}

resource "aws_kms_alias" "signing_keys" {
  name          = "alias/${var.name}-signing-keys"
  target_key_id = aws_kms_key.signing_keys.key_id
}

# The API pod's IRSA role: scoped to Encrypt/Decrypt on the one key, bound to the kaleido/kaleido-api service account.
data "aws_iam_policy_document" "api_kms" {
  statement {
    actions   = ["kms:Encrypt", "kms:Decrypt"]
    resources = [aws_kms_key.signing_keys.arn]
  }
}

resource "aws_iam_policy" "api_kms" {
  name   = "${var.name}-api-kms"
  policy = data.aws_iam_policy_document.api_kms.json
}

module "api_irsa" {
  source  = "terraform-aws-modules/iam/aws//modules/iam-role-for-service-accounts-eks"
  version = "~> 5.48"

  role_name = "${var.name}-api"

  role_policy_arns = {
    kms = aws_iam_policy.api_kms.arn
  }

  oidc_providers = {
    main = {
      provider_arn               = module.eks.oidc_provider_arn
      namespace_service_accounts = ["kaleido:kaleido-api"]
    }
  }
}
