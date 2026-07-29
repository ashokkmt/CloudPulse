resource "digitalocean_vpc" "main" {
  name        = var.vpc_name
  region      = var.region
  ip_range    = "10.106.0.0/20"
  description = "All the other resources will be enclosed in this network"
}
