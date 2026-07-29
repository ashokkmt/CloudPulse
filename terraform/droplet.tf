resource "digitalocean_droplet" "app" {
  count    = 1
  image    = "ubuntu-22-04-x64"
  name     = "cloudpulse-app-${count.index + 1}"
  region   = var.region
  size     = "s-4vcpu-8gb" # Update if your droplet size is different
  vpc_uuid = digitalocean_vpc.main.id
  tags     = ["cloudpulse-app"]
  ssh_keys = [
    "58003540", # github-actions
    "57896375"  # windows-pc
  ]
}

resource "digitalocean_droplet" "monitor" {
  image    = "ubuntu-22-04-x64"
  name     = "cloudpulse-monitor"
  region   = var.region
  size     = "s-2vcpu-4gb" # Update if your droplet size is different
  vpc_uuid = digitalocean_vpc.main.id
  tags     = ["cloudpulse-monitoring"]
  ssh_keys = [
    "58003540", # github-actions
    "57896375"  # windows-pc
  ]
}

output "app_droplet_ips" {
  description = "The public IP addresses of the application droplets"
  value       = digitalocean_droplet.app[*].ipv4_address
}

output "monitor_droplet_ip" {
  description = "The public IP address of the monitoring droplet"
  value       = digitalocean_droplet.monitor.ipv4_address
}
