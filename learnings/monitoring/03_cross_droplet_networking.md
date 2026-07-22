# Cross-Droplet Monitoring Setup

## The Problem
Locally, Docker places all containers on a virtual `bridge` network. `prometheus` can simply connect to `http://backend:8000` because Docker's internal DNS automatically resolves the name `backend` to the correct container IP. 
In production, your app is on the "Deploy" Droplet and Prometheus is on the "Monitor" Droplet. Docker DNS does not work across different servers.

## The Solution: VPC Private IPs
When you created your Droplets, you placed them in the same DigitalOcean Virtual Private Cloud (VPC). This means they are on a shared, highly secure private network (usually starting with `10.x.x.x`).

We will use the **Deploy Droplet's Private IP** to allow Prometheus (on the Monitor Droplet) to scrape metrics across the private network.

---

## Step 1: Install Exporters on the Deploy Droplet
Right now, `node-exporter` and `cadvisor` are configured in your monitoring compose file. This means they monitor the CPU and RAM of the *monitoring server itself*, not the app server! 

To fix this, you must run those two exporters on the Deploy Droplet. Add this to your `docker-compose.prod.yml` on the **Deploy Droplet**:
```yaml
  node-exporter:
    image: prom/node-exporter:latest
    container_name: node-exporter
    ports:
      - "9100:9100"
    restart: always

  cadvisor:
    image: gcr.io/cadvisor/cadvisor:latest
    container_name: cadvisor
    ports:
      - "8080:8080"
    volumes:
      - /:/rootfs:ro
      - /var/run:/var/run:rw
      - /sys:/sys:ro
      - /var/lib/docker/:/var/lib/docker:ro
    privileged: true
    restart: always
```
*(Note: You can leave them running on the Monitor droplet too if you want to monitor both servers!)*

## Step 2: Configure the DO Cloud Firewall
By default, your DO Cloud Firewall blocks all incoming traffic except ports 80, 443, and 22. 
You must allow Prometheus to reach the exporters on the Deploy Droplet.

In the DigitalOcean Control Panel, go to Networking -> Firewalls -> `cloudpulse-fw`.
Add **Inbound Rules**:
- **Custom (TCP) on Port 8000** (For your Go API `/metrics`)
- **Custom (TCP) on Port 9100** (For Node Exporter)
- **Custom (TCP) on Port 8080** (For cAdvisor)

**CRITICAL SECURITY STEP:** Under "Sources" for these three rules, **do NOT select "All IPv4"**. Instead, type the name of your Monitor Droplet or select your VPC. This ensures the public internet cannot access your internal metrics!

## Step 3: Find the Private IP
Go to the DO Control Panel -> Droplets. Click on your **Deploy Droplet**. 
Look for the **Private IP** (it will look something like `10.106.x.x`). Copy this.

## Step 4: Update Prometheus Targets on the Monitor Droplet
SSH into your Monitor Droplet. Open your `monitoring/prometheus/prometheus.yml` file. 
Replace the local targets with the **Private IP** of your Deploy Droplet. 

```yaml
scrape_configs:
  # This monitors the Prometheus server itself (leave as localhost)
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']

  - job_name: 'app'
    static_configs:
      - targets: ['10.106.x.x:8000'] # <-- Replace with actual Private IP

  - job_name: 'node'
    static_configs:
      - targets: ['10.106.x.x:9100'] # <-- Replace with actual Private IP

  - job_name: 'cadvisor'
    static_configs:
      - targets: ['10.106.x.x:8080'] # <-- Replace with actual Private IP
```

## Step 5: Restart Prometheus
On the Monitor Droplet, apply the changes by restarting Prometheus:
```bash
docker-compose -f docker-compose.monitoring.prod.yml restart prometheus
```

## Verification
Now, open your browser and go to `http://<Monitor-Droplet-Public-IP>:9090/targets`. 
You should see `app`, `node`, and `cadvisor` all listed as **UP**, proving that your Monitor Droplet is successfully reaching across the VPC to scrape the Deploy Droplet!
