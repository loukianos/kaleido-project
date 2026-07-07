terraform {
  backend "s3" {
    bucket       = "kaleido-project-tfstate-433484250096-us-east-1"
    key          = "kaleido/terraform.tfstate"
    region       = "us-east-1"
    encrypt      = true
    use_lockfile = true
  }
}
