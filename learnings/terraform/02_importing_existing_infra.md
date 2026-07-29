    # 02: Importing Existing Infrastructure

Because you have already created your Droplets, Databases, VPC, and Spaces manually via the DO Console, Terraform doesn't know about them yet. 

If you just run `terraform apply` right now, Terraform will try to create *brand new* Droplets and Databases because its local "state file" is empty.

To fix this, we must **import** your existing resources into the state file. We do this by mapping the Terraform resource name (e.g., `digitalocean_droplet.app`) to the actual DigitalOcean UUID.

## How to find your UUIDs
1. **VPC:** Networking -> VPC -> Click your VPC -> Look at the URL for the UUID.
2. **Managed Databases:** Databases -> Click the cluster -> Look at the URL for the UUID.
3. **Spaces:** Simply use the name of the bucket (e.g. `cloudpulsebucket`).
4. **Registry:** Simply use the name of the registry (e.g. `cloudpulse-registry`).
5. **Firewall:** Networking -> Firewalls -> Click your firewall -> Look at the URL for the UUID.

## The Import Commands
Once you have the UUIDs, run these commands inside your `terraform/` directory.

### 1. Import VPC
```bash
terraform import digitalocean_vpc.main <YOUR_VPC_UUID>
```

### 2. Import Databases
```bash
terraform import digitalocean_database_cluster.postgres <YOUR_POSTGRES_UUID>
terraform import digitalocean_database_cluster.redis <YOUR_REDIS_UUID>
```

### 3. Import Spaces & CDN
```bash
terraform import digitalocean_spaces_bucket.cloudpulse sgp1,cloudpulsebucket
# Note: Spaces import requires "region,bucketname"

terraform import digitalocean_cdn.cloudpulse_cdn <YOUR_CDN_ENDPOINT_ID>
```

### 4. Import Registry
```bash
terraform import digitalocean_container_registry.cloudpulse cloudpulse-registry
```

### 5. Import Firewall
```bash
terraform import digitalocean_firewall.web <YOUR_FIREWALL_UUID>
```

## Verify the Import
Once you have imported everything, run:
```bash
terraform plan
```
If your `.tf` files perfectly match your actual configuration, the output will say:
**"No changes. Your infrastructure matches the configuration."**

If it says it wants to change or destroy something (like changing a tag or a droplet size), **do not run apply**. Instead, update your `.tf` files to match the reality of what's deployed, and run `terraform plan` again until it is perfectly clean.
