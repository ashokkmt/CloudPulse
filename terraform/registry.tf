resource "digitalocean_container_registry" "cloudpulse" {
  name                   = "cloudpulse-registry"
  subscription_tier_slug = "basic" # Update if you are on a different tier (starter, basic, professional)
  region                 = var.region
}
