resource "kuma_http_monitor" "example" {
  name        = "example"
  description = "example monitor"
  url         = "https://example.com"

  interval       = 60
  retry_interval = 60
  ignore_tls     = true

  accepted_statuscodes = ["200-299", "300-399"]
  notification_list    = [1, 2, 3]
  tags = {
    env      = "prod"
    createBy = "demo"
  }
}