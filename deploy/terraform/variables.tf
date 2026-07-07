variable "region" {
  description = "AWS region for the whole stack"
  type        = string
  default     = "us-east-1"
}

variable "aws_profile" {
  description = "AWS shared-config profile; empty uses the default credential chain (CI)"
  type        = string
  default     = "loukianos"
}

variable "name" {
  description = "Base name for the cluster and its resources"
  type        = string
  default     = "kaleido"
}

variable "kubernetes_version" {
  description = "EKS control-plane version"
  type        = string
  default     = "1.36"
}

variable "node_instance_type" {
  description = "Managed node group instance type; the Besu network plus Keycloak and the API need some headroom"
  type        = string
  default     = "t3.large"
}

variable "node_count" {
  description = "Managed node group desired size"
  type        = number
  default     = 2
}

variable "db_name" {
  description = "Application database name and username"
  type        = string
  default     = "loan_notes"
}
