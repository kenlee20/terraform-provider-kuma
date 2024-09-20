# Configuration-based authentication
provider "kuma" {
  username = "admin"
  password = "admin"
  host     = "http://localhost:8000"
}

terraform {
  required_providers {
    kuma = {
      source  = "registry.terraform.io/kenlee20/kuma"
      version = "0.1.0"
    }
  }
}
