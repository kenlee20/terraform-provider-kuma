terraform {
  required_providers {
    kuma = {
      source = "registry.terraform.io/kenli/kuma"
    }
  }
}

provider "kuma" {
  host     = "http://127.0.0.1:8000"
  username = "admin"
  password = "admin"
}

# resource "kuma_http_monitor" "name" {
#   name        = "demo_monitor"
#   description = "demo1"
#   url         = "https://google.com"
# }

data "kuma_monitors" "name" {
}

output "name" {
  value = data.kuma_monitors.name
}