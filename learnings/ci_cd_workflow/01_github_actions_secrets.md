# Setting Up GitHub Actions Secrets for DigitalOcean

When configuring a CI/CD pipeline using GitHub Actions to deploy to a DigitalOcean Droplet and authenticate with the DigitalOcean Container Registry (DOCR), you need to provide specific credentials. 

To keep these credentials secure, they must be stored as **GitHub Secrets**.

## Where to Configure GitHub Secrets
1. Go to your repository on GitHub.
2. Click on the **Settings** tab at the top.
3. In the left sidebar, navigate to **Secrets and variables** -> **Actions**.
4. Click the green **New repository secret** button for each of the keys below.

---

## Required Secrets & Where to Find Them

### 1. `DIGITALOCEAN_ACCESS_TOKEN`

**What it is:** The official `digitalocean/action-doctl` GitHub Action requires a standard **Personal Access Token (PAT)**. *Do not use a Docker-specific registry token here.* The PAT is used to install and configure `doctl` inside the GitHub Action runner, which then seamlessly logs into DOCR behind the scenes.

**Where to get it:**
1. Log into your DigitalOcean Control Panel.
2. On the bottom left sidebar, click **API**.
3. Under the **Tokens/Custom Images** tab, click **Generate New Token**.
4. Name it (e.g., `GitHub Actions CI`), ensure it has **Write** permissions (required for pushing images), and set an expiration date (or never expire).
5. Copy the token immediately (you will not be able to see it again).

### 2. `REGISTRY_NAME`

**What it is:** The unique name of your DigitalOcean Container Registry where your Docker images are stored.

**Where to get it:**
1. In the DO Control Panel, click **Container Registry** in the left sidebar.
2. Your registry name is the name you chose when you created it (e.g., `cloudpulse`).
3. You can also see it in your registry endpoint URL: `registry.digitalocean.com/YOUR_REGISTRY_NAME`.

### 3. `DROPLET_SSH_KEY`

**What it is:** The **private** SSH key that corresponds to the public key authorized on your Droplet. GitHub Actions uses this to securely SSH into your Droplet to pull the new image and restart the containers.

**Where to get it:**
1. This is the private key on your local machine that you use to log into your droplet (e.g., the contents of `~/.ssh/id_rsa` or `~/.ssh/id_ed25519`).
2. Run `cat ~/.ssh/id_rsa` (or your specific key name) on your local machine.
3. Copy the *entire* output, including the `-----BEGIN OPENSSH PRIVATE KEY-----` and `-----END OPENSSH PRIVATE KEY-----` lines, and paste it into the GitHub secret.

### 4. `DROPLET_IP`

**What it is:** The public IPv4 address of your App Droplet.

**Where to get it:**
1. Go to the DO Control Panel -> **Droplets**.
2. Copy the IPv4 address of your target droplet (the one running the CloudPulse backend).

### 5. `DB_PASSWORD`

**What it is:** The password for your DigitalOcean Managed PostgreSQL database. This is often required during the CI/CD pipeline if you run integration tests or database migrations during the GitHub Actions workflow before deploying.

**Where to get it:**
1. Go to the DO Control Panel -> **Databases**.
2. Click on your PostgreSQL database cluster.
3. In the **Connection Details** section, click the "Show" button next to the password field, and copy it.

---

## Tagging Strategy Reminder

In your workflow YAML, avoid relying solely on the `latest` tag. When you build and push your images, use the Git Commit SHA as the tag:

```yaml
# Example tagging in GitHub Actions
IMAGE_TAG: ${{ github.sha }}
```

This guarantees that:
- Every image is uniquely identifiable.
- You know exactly which code commit corresponds to which image.
- If a deployment fails, you can instantly rollback by pulling the previous SHA tag.
