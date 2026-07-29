# 01: Terraform Fundamentals & Setup

## Why Infrastructure as Code (IaC)?
Until now, you've been clicking through the DigitalOcean web UI to create Droplets, Managed Databases, and Spaces. This is great for learning, but unacceptable in a production environment. 

With **Terraform**, you declare your infrastructure in code (`.tf` files). This provides:
1. **Reproducibility:** Spin up an identical staging environment in 5 minutes.
2. **Version Control:** Infrastructure changes are reviewed in Pull Requests (PRs), audited, and tracked in Git.
3. **Drift Detection:** If someone manually deletes a firewall rule, Terraform will detect it and change it back on the next run.
4. **Disaster Recovery:** You can rebuild your entire cloud environment instantly if something gets deleted.

## Core Concepts
- **Provider:** A plugin that tells Terraform how to talk to a specific Cloud's API (e.g., the DigitalOcean provider).
- **Resource:** A piece of infrastructure (e.g., a `digitalocean_droplet`, `digitalocean_vpc`).
- **State File (`terraform.tfstate`):** A JSON file where Terraform maps your code to the actual real-world resources (via their IDs). **Never lose this file.**
- **Plan:** A dry-run (`terraform plan`) showing what Terraform *intends* to create, modify, or destroy.
- **Apply:** The command (`terraform apply`) that actually makes the API calls to execute the plan.

## Step 1: Install Terraform
If you haven't already, download and install the Terraform CLI for your operating system:
- [Terraform Installation Guide](https://developer.hashicorp.com/terraform/downloads)

## Step 2: Generate a DigitalOcean API Token
Terraform needs permission to create and destroy resources on your behalf.
1. Go to your DigitalOcean Dashboard.
2. On the left sidebar, click **API**.
3. Under the **Tokens/Custom Keys** tab, click **Generate New Token**.
4. Name it `Terraform` and give it **Read** and **Write** scopes.
5. **Copy the token immediately.** You will not be able to see it again!

## Step 3: Export the Token Locally
Terraform will automatically look for an environment variable named `TF_VAR_do_token`.

### On Windows (PowerShell):
```powershell
$env:TF_VAR_do_token="your_long_api_token_here"
```

### On macOS / Linux:
```bash
export TF_VAR_do_token="your_long_api_token_here"
```

## Step 4: Initialize Terraform
Go to the root of your project where the `terraform/` folder is located:
```bash
cd terraform
terraform init
```
This command downloads the DigitalOcean provider plugin required to manage your resources.
