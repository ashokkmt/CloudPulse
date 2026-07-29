terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.67"
    }
  }

  # Once you configure Spaces for state, you will uncomment this:
  # backend "s3" {
  #   endpoint                    = "sgp1.digitaloceanspaces.com" # Change if using a different region
  #   region                      = "us-east-1"                   # S3 compatibility requirement, leave as is
  #   bucket                      = "cloudpulse-tf-state"         # Ensure this bucket is created first
  #   key                         = "terraform.tfstate"
  #   skip_credentials_validation = true
  #   skip_metadata_api_check     = true
  # }
}

provider "digitalocean" {
  token             = var.do_token
  spaces_access_id  = var.spaces_access_key
  spaces_secret_key = var.spaces_secret_key
}
