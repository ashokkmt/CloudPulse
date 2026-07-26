# Deploying to Multiple Droplets via CI/CD

Now that you have successfully scaled horizontally by adding a second Droplet (`cloudpulse-app-02`), your CI/CD pipeline needs to be updated to deploy new code to both servers automatically.

Because `cloudpulse-app-02` was created from a snapshot of `cloudpulse-app-01`, it shares the exact same `deploy` user and SSH keys. This makes updating the GitHub Actions workflow incredibly simple!

## Step-by-Step Instructions

### Step 1: Add the New IP to GitHub Secrets
Your GitHub Actions runner needs to know the IP address of the second droplet.
1. Go to your repository on GitHub.
2. Navigate to **Settings** -> **Secrets and variables** -> **Actions**.
3. Click **New repository secret**.
4. Name: `DROPLET_02_IP`
5. Secret: *<Paste the public IPv4 address of cloudpulse-app-02>*
6. Click **Add secret**.

*(Note: You do not need to add a new SSH key secret because Droplet 2 inherited the authorized keys directly from Droplet 1's snapshot!)*

### Step 2: Update Your Deployment Workflow
Open your `.github/workflows/deploy.yml` file. Scroll to the very bottom where the `Deploy` step is defined.

You have two options for deploying to multiple servers using the `appleboy/ssh-action`:

#### Option A: The "Multiple Hosts" Method (Fastest)
The `appleboy/ssh-action` supports a comma-separated list of hosts. You can simply update the `host` line to include both secrets. It will connect and run the deployment script on both servers.

```yaml
      - name: Deploy to All Droplets
        uses: appleboy/ssh-action@v1.2.2
        with:
          host: ${{ secrets.DROPLET_IP }},${{ secrets.DROPLET_02_IP }}
          username: deploy
          key: ${{ secrets.DROPLET_SSH_KEY }}
          script: |
            cd cloudpulse
            export APP_VERSION=${{ github.sha }}
            chmod +x deploy.sh
            ./deploy.sh
```

#### Option B: The "Separate Steps" Method (Better for Debugging)
If you prefer to see clearly which droplet succeeded or failed in the GitHub Actions UI logs, you can duplicate the deploy step so they run sequentially as distinct UI items:

```yaml
      - name: Deploy to Droplet 1
        uses: appleboy/ssh-action@v1.2.2
        with:
          host: ${{ secrets.DROPLET_IP }}
          username: deploy
          key: ${{ secrets.DROPLET_SSH_KEY }}
          script: |
            cd cloudpulse
            export APP_VERSION=${{ github.sha }}
            chmod +x deploy.sh
            ./deploy.sh

      - name: Deploy to Droplet 2
        uses: appleboy/ssh-action@v1.2.2
        with:
          host: ${{ secrets.DROPLET_02_IP }}
          username: deploy
          key: ${{ secrets.DROPLET_SSH_KEY }}
          script: |
            cd cloudpulse
            export APP_VERSION=${{ github.sha }}
            chmod +x deploy.sh
            ./deploy.sh
```

### Step 3: Commit and Test
Once you've made the changes to your `deploy.yml`, commit and push them to the `main` branch. 

Your next push will trigger the pipeline, build the latest Docker images, push them to the DigitalOcean Container Registry, and then securely SSH into **both** droplets to pull and run the exact same `APP_VERSION`!
