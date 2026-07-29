resource "digitalocean_database_cluster" "postgres" {
  name                 = "cloudpulse-pg-prod"
  engine               = "pg"
  version              = "15"
  size                 = "db-s-1vcpu-1gb"
  region               = var.region
  node_count           = 1
  private_network_uuid = digitalocean_vpc.main.id
}

resource "digitalocean_database_cluster" "redis" {
  name                 = "cloudpulse-redis-prod"
  engine               = "redis"
  version              = "7"
  size                 = "db-s-1vcpu-1gb"
  region               = var.region
  node_count           = 1
  private_network_uuid = digitalocean_vpc.main.id
}
