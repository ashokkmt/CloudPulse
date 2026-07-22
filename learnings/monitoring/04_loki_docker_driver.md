# Configuring Loki as the Default Docker Log Driver

By default, Docker stores logs locally in JSON format on the server. To centralize them, we can tell the Docker daemon to automatically send *all* logs from *every* container directly to our Loki server.

You need to perform these steps on **both** your Deploy Droplet and your Monitor Droplet.

## Step 1: Install the Loki Docker Driver Plugin
SSH into each of your Droplets and run this command to install the official Grafana Loki plugin for Docker:

```bash
docker plugin install grafana/loki-docker-driver:3.0.0 --alias loki --grant-all-permissions
```

## Step 2: Configure the Docker Daemon
We need to tell Docker to use this new plugin as its default log driver. 

Create or edit the Docker daemon configuration file on each droplet:
```bash
sudo nano /etc/docker/daemon.json
```

Paste the following JSON into the file. 
> [!IMPORTANT]
> Replace `<MONITOR_DROPLET_PRIVATE_IP>` with the **Private IP of your Monitor Droplet** (e.g., `10.106.x.x`). 
> For the Monitor Droplet itself, you can use `http://localhost:3100...` if you prefer, but using the Private IP on both servers is totally fine and keeps things consistent.

```json
{
  "log-driver": "loki",
  "log-opts": {
    "loki-url": "http://<MONITOR_DROPLET_PRIVATE_IP>:3100/loki/api/v1/push",
    "loki-retries": "5",
    "loki-batch-size": "400",
    "keep-file": "true"
  }
}
```
*Note: `keep-file: true` ensures that you can still use the `docker logs <container>` command locally on the server in an emergency.*

## Step 3: Restart Docker
Apply the changes by restarting the Docker daemon on both servers:
```bash
sudo systemctl restart docker
```

## Step 4: Recreate Existing Containers
The new default logging driver only applies to *newly created* containers. Any containers that are currently running are still using the old local JSON driver.

To force your existing containers to recreate and pick up the Loki driver, go to the directory containing your `docker-compose.prod.yml` and run:

**On Deploy Droplet:**
```bash
docker-compose -f docker-compose.prod.yml up -d --force-recreate
```

**On Monitor Droplet:**
```bash
docker-compose -f docker-compose.monitoring.prod.yml up -d --force-recreate
```

## Step 5: Update the Monitor Droplet Firewall
Since your Deploy Droplet is now sending logs to the Monitor Droplet over port `3100`, you must ensure the Monitor Droplet's firewall allows incoming traffic on port `3100`.

In the DigitalOcean Control Panel, add an Inbound Rule to your firewall:
- **Custom (TCP) on Port 3100** 
- **Source:** Select your Deploy Droplet or VPC (Do NOT select All IPv4).

## Step 6: Verify in Grafana
1. Open Grafana (`http://<Monitor-Droplet-IP>:3001`).
2. Go to **Explore**.
3. Select **Loki** from the datasource dropdown.
4. Run the query `{compose_project!=""}` to see logs streaming in from all your Docker Compose services across both droplets!
