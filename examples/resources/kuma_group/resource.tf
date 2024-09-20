resource "kuma_group" "example" {
  name        = "example"
  description = "example"
  tags = {
    env      = "prod"
    createBy = "demo"
  }
}