# Configuration-based authentication
provider "upkuapi" {
  username = "admin"
  password = "admin"
  host     = "http://localhost:8000"
}

terraform {
  required_providers {
    upkuapi = {
      source  = "registry.terraform.io/kenlee20/upkuapi"
      version = "1.0.0"
    }
  }
}
