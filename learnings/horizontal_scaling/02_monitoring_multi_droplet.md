# Monitoring Multi-Droplet Architecture

When scaling horizontally by adding a second Droplet (`cloudpulse-app-02`), you need to update your centralized **Monitoring Droplet** (running Prometheus, Grafana, Loki, and Tempo) to scrape metrics and collect logs/traces from both app Droplets.

---

## 1. Network Requirements (VPC & Firewalls)

Ensure all Droplets are inside the same **CloudPulse VPC** so they can communicate over private IP addresses (e.g. `10.106.0.x`).

Make sure your **Monitoring Droplet Firewall** allows inbound connections on:
- **Port `3100`** (Loki) from `cloudpulse-app-02` (for log pushing).
- **Port `4318` / `4317`** (Tempo OTLP) from `cloudpulse-app-02` (for trace pushing).

Make sure your **App Droplet Firewall** allows inbound connections from the **Monitoring Droplet** on:
- **Port `8000`** (Backend `/metrics`).
- **Port `9100`** (Node Exporter - if installed on App Droplets).
- **Port `8080`** (cAdvisor - if installed on App Droplets).

---

## 2. Configuring Prometheus to Scrape Both Droplets

On your **Monitoring Droplet**, edit `/etc/prometheus/prometheus.yml` (or `~/monitoring/prometheus/prometheus.yml`):

```yaml
scrape_configs:
  - job_name: 'cloudpulse-backend'
    metrics_path: '/metrics'
    static_configs:
      - targets:
          - '<DROPLET_01_PRIVATE_IP>:8000'
          - '<DROPLET_02_PRIVATE_IP>:8000'
        labels:
          environment: 'production'

  - job_name: 'node-exporter'
    static_configs:
      - targets:
          - '<DROPLET_01_PRIVATE_IP>:9100'
          - '<DROPLET_02_PRIVATE_IP>:9100'

  - job_name: 'cadvisor'
    static_configs:
      - targets:
          - '<DROPLET_01_PRIVATE_IP>:8080'
          - '<DROPLET_02_PRIVATE_IP>:8080'
```

After updating `prometheus.yml`, reload Prometheus on your Monitoring Droplet:
```bash
docker exec -it prometheus kill -HUP 1
# or: docker compose -f docker-compose.monitoring.prod.yml restart prometheus
```

---

## 3. Configuring Promtail on Droplet 2 (Log Collection)

On `cloudpulse-app-02`, Promtail (or a lightweight log shipper) must send container logs to the central Loki instance on the Monitoring Droplet.

In `promtail-config.yml` on Droplet 2:
```yaml
clients:
  - url: http://<MONITORING_DROPLET_PRIVATE_IP>:3100/loki/api/v1/push
```

In `docker-compose.prod.yml` on Droplet 2, ensure Promtail runs and mounts `/var/run/docker.sock`.

---

## 4. Configuring OpenTelemetry Traces (Tempo)

In the `.env` file on `cloudpulse-app-02`, point the OpenTelemetry exporter to your Monitoring Droplet's private IP:

```env
OTEL_EXPORTER_OTLP_ENDPOINT=http://<MONITORING_DROPLET_PRIVATE_IP>:4318
```

Restart the backend container on Droplet 2:
```bash
docker compose restart backend
```

---

## 5. Verification Checklist

1. **Prometheus UI (`http://<MONITORING_DROPLET_IP>:9090/targets`)**:
   - Verify that both `<DROPLET_01_PRIVATE_IP>:8000` and `<DROPLET_02_PRIVATE_IP>:8000` show state **UP**.

2. **Grafana Dashboards (`http://<MONITORING_DROPLET_IP>:3000`)**:
   - Check your HTTP requests dashboard. You should see incoming metrics labeled from both target IPs.
   - Filter logs in Loki (`Explore` tab) by `container="cloudpulse-backend-1"` to verify logs are arriving from both Droplets.
