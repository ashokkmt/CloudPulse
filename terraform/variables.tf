variable "do_token" {
  description = "DigitalOcean API Token"
  type        = string
  sensitive   = true
}

variable "region" {
  description = "The region to deploy resources into"
  type        = string
  default     = "sgp1" # Change if using a different region (e.g. nyc3, fra1)
}

variable "vpc_name" {
  description = "The name of the VPC"
  type        = string
  default     = "cloudpulse-vpc"
}

variable "spaces_access_key" {
  description = "DO Spaces Access Key"
  type        = string
  sensitive   = true
}

variable "spaces_secret_key" {
  description = "DO Spaces Secret Key"
  type        = string
  sensitive   = true
}
