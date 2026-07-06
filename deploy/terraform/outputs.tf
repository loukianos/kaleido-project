output "region" {
  value = var.region
}

output "cluster_name" {
  value = module.eks.cluster_name
}

output "ecr_repository_url" {
  value = aws_ecr_repository.api.repository_url
}

output "database_url" {
  description = "Postgres connection string for the API and migrations"
  value       = "postgres://${var.db_name}:${random_password.db.result}@${aws_db_instance.postgres.endpoint}/${var.db_name}?sslmode=require"
  sensitive   = true
}

output "kms_key_id" {
  value = aws_kms_key.signing_keys.key_id
}

output "api_irsa_role_arn" {
  value = module.api_irsa.iam_role_arn
}
