# 03: Remote State in DO Spaces

When you run `terraform import` or `terraform apply`, Terraform generates a file called `terraform.tfstate`.

This file contains the exact mapping of your code to the real-world Cloud resources, including sensitive information (like database connection strings, albeit not strictly in our basic setup, but often present).

**If you lose this file, Terraform loses track of your infrastructure.** You would have to re-import everything manually. If you are working on a team, you cannot just keep this file on your local laptop, because your coworkers won't have it.

## The Solution: Remote State
We can tell Terraform to store this `terraform.tfstate` file inside our DigitalOcean Spaces bucket! This way, it is safe in the cloud, backed up, and accessible to a team or a CI/CD pipeline.

### Step 1: Create a bucket for your state
It is best practice to have a dedicated bucket for Terraform state, separate from your application bucket.
1. Go to the DO Console -> Spaces.
2. Create a new bucket called `cloudpulse-tf-state` in the `sgp1` region.
3. Keep it **Private** (never make a state bucket public!).

### Step 2: Generate Spaces Access Keys
Terraform needs an S3-compatible Access Key and Secret Key to write to Spaces.
1. Go to DO Console -> API -> Spaces Keys.
2. Generate a new key.
3. Export them in your terminal:
   - Windows (PowerShell):
     ```powershell
     $env:AWS_ACCESS_KEY_ID="your_spaces_access_key"
     $env:AWS_SECRET_ACCESS_KEY="your_spaces_secret_key"
     ```
   - macOS / Linux:
     ```bash
     export AWS_ACCESS_KEY_ID="your_spaces_access_key"
     export AWS_SECRET_ACCESS_KEY="your_spaces_secret_key"
     ```
   *(Note: Terraform uses the AWS S3 backend plugin for DO Spaces, which is why the variables start with `AWS_`)*

### Step 3: Enable the Remote Backend
Open your `terraform/main.tf` file and uncomment the `backend "s3"` block:

```hcl
  backend "s3" {
    endpoint                    = "sgp1.digitaloceanspaces.com" 
    region                      = "us-east-1"                   
    bucket                      = "cloudpulse-tf-state"         
    key                         = "terraform.tfstate"
    skip_credentials_validation = true
    skip_metadata_api_check     = true
  }
```

### Step 4: Migrate the State
Run:
```bash
terraform init
```
Terraform will detect that you added a remote backend and ask:
> "Do you want to copy existing state to the new backend?"

Type **yes**.

Congratulations! Your infrastructure state is now safely backed up in the cloud. You can safely delete your local `terraform.tfstate` file.
