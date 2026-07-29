resource "digitalocean_spaces_bucket" "cloudpulse" {
  name   = "cloudpulsebucket"
  region = var.region
  # acl    = "public-read"
}

# CDN for the Spaces bucket
resource "digitalocean_cdn" "cloudpulse_cdn" {
  origin = digitalocean_spaces_bucket.cloudpulse.bucket_domain_name
}
