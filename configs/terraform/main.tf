locals {
  app_name = "${terraform.workspace}-${var.app_name}"
}

resource "fly_app" "nba_stats_api" {
  name = local.app_name
  org  = "personal"
}

resource "fly_ip" "nba_stats_api_ip" {
  app        = local.app_name
  type       = "v4"
  depends_on = [fly_app.nba_stats_api]
}

resource "fly_ip" "nba_stats_api_ip_v6" {
  app        = local.app_name
  type       = "v6"
  depends_on = [fly_app.nba_stats_api]
}

resource "fly_cert" "nba_stats_api_cert" {
  app        = local.app_name
  hostname   = "nbastats.api.pedromealha.dev"
  depends_on = [fly_app.nba_stats_api]
}
