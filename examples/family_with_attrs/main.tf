terraform {
  required_providers {
    akeneo = {
      source  = "0xfrej/akeneo"
      version = "0.2.0"
    }
  }
}

provider "akeneo" {
  host              = "pim.example.com"
  api_username      = "connection username"
  api_password      = "connection password"
  api_client_id     = "connection api client id"
  api_client_secret = "connection api client secret"
}