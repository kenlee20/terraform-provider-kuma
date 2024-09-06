terraform {
  required_providers {
    hashicups = {
      source = "gitlab.microfusion.cloud/kenli/kuma"
    }
  }
}

provider "hashicups" {}

data "hashicups_coffees" "example" {}