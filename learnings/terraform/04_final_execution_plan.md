# 04: Final Execution Plan

Now that you understand the concepts, here is the exact step-by-step checklist to execute Day 9.

## Pre-requisites
1. Destroy the DOKS cluster manually in the DO Console (as instructed in the Day 9 roadmap) to save costs.
2. Ensure you have the `terraform` CLI installed.
3. Open a terminal in the `d:\projects\Digital Ocean Labs\terraform` directory.

---

## Step 1: Export your DigitalOcean Token
```powershell
$env:TF_VAR_do_token="<YOUR_DO_API_TOKEN>"
```

## Step 2: Initialize Terraform
```bash
terraform init
```
This downloads the DigitalOcean provider.

## Step 3: Gather your UUIDs
Go to the DigitalOcean web console and copy the UUIDs from the browser URLs for:
- VPC
- Firewall
- App Droplet (`cloudpulse-app-1`)
- Monitor Droplet (`cloudpulse-monitor`)
- PostgreSQL Database
- Redis Database
- CDN Endpoint ID

## Step 4: Run the Import Commands
Run these commands one by one, replacing the placeholders with your actual UUIDs.
*(Note: Since you deleted your Droplets, you will NOT import them. Terraform will create brand new ones!)*

```bash
terraform import digitalocean_vpc.main <VPC_UUID>
terraform import digitalocean_firewall.web <FIREWALL_UUID>

terraform import digitalocean_database_cluster.postgres <POSTGRES_UUID>
terraform import digitalocean_database_cluster.redis <REDIS_UUID>

terraform import digitalocean_spaces_bucket.cloudpulse sgp1,cloudpulsebucket
terraform import digitalocean_cdn.cloudpulse_cdn <CDN_UUID>

terraform import digitalocean_container_registry.cloudpulse cloudpulse-registry
```

## Step 5: Verify the State (The Dry Run)
Run a plan:
```bash
terraform plan
```
Read the output carefully. 
- If it says **"No changes. Your infrastructure matches the configuration."**, you have successfully imported everything perfectly!
- If it shows changes (e.g. changing the droplet size, or tags), it means your `.tf` files don't perfectly match what you created manually. Update your `.tf` files to match the real resources, then run `terraform plan` again until it is perfectly clean.

## Step 6: Configure Remote State (Spaces)
Once your local state is perfectly matching, move it to the cloud.
1. Create a `cloudpulse-tf-state` bucket in Spaces.
2. Generate Spaces Access Keys.
3. Export the keys:
   ```powershell
   $env:AWS_ACCESS_KEY_ID="<SPACES_ACCESS_KEY>"
   $env:AWS_SECRET_ACCESS_KEY="<SPACES_SECRET_KEY>"
   ```
4. Edit `terraform/main.tf` and uncomment the `backend "s3"` block.
5. Run:
   ```bash
   terraform init
   ```
6. Type `yes` to migrate the state to Spaces.

## Congratulations!
You have successfully transitioned your click-ops infrastructure into declarative Infrastructure as Code. In the future, to add a new server, you just copy-paste a droplet block in `droplet.tf` and run `terraform apply`.
