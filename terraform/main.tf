terraform {
  required_providers {
    hcp = {
      source  = "hashicorp/hcp"
      version = "0.95.1"
    }
  }
}

variable "hcp_client_id" {
  type = string
}
variable "hcp_client_secret" {
  type = string
}
variable "app_name" {
  type = string
}

provider "hcp" {
  client_id     = var.hcp_client_id
  client_secret = var.hcp_client_secret
}

data "hcp_vault_secrets_app" "demo_app" {
  app_name = var.app_name
}

output "demo_app_secrets" {
  value     = data.hcp_vault_secrets_app.demo_app.secrets
  sensitive = true
}
