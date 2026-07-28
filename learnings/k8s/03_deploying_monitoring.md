# Deploying the Monitoring Stack

This guide explains how to deploy Prometheus, Grafana, Loki, and Tempo to monitor your Kubernetes cluster using Helm.

### Prerequisites
- Helm installed locally (`brew install helm` on Mac / `choco install kubernetes-helm` on Windows).
- The `monitoring` namespace must already be created (applied via `namespace.yaml`).

### Step 1: Prometheus + Grafana + Alertmanager
The `kube-prometheus-stack` is an incredible all-in-one chart. It includes Prometheus, Grafana, Alertmanager, pre-built Kubernetes dashboards, ServiceMonitors for dynamic discovery, and node-exporter DaemonSets for node-level metrics.
```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install cloudpulse-monitoring prometheus-community/kube-prometheus-stack \
  --namespace monitoring
```

### Step 2: Access Grafana
You can access Grafana securely using port-forwarding:
```bash
kubectl port-forward svc/cloudpulse-monitoring-grafana 3000:80 -n monitoring
```
Then visit `http://localhost:3000`.
**Default login**: `admin` / `prom-operator`

> **Note**: In K8s, Prometheus auto-discovers pods via `ServiceMonitors`. No manual IP configuration is needed unlike our old Droplet setup!

### Step 3: Loki (Log Aggregation)
Loki aggregates logs. Promtail is its agent, running as a DaemonSet (one on each node), automatically collecting logs from ALL containers running in the cluster.
```bash
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
helm install loki grafana/loki-stack \
  --namespace monitoring \
  --set promtail.enabled=true
```

### Step 4: Tempo (Distributed Tracing)
Tempo handles application tracing data.
```bash
helm install tempo grafana/tempo \
  --namespace monitoring
```

### Step 5: Connect Data Sources in Grafana
Once you've port-forwarded into Grafana, go to **Connections** → **Data Sources** → **Add data source**:
- **Loki**: `http://loki.monitoring.svc.cluster.local:3100`
- **Tempo**: `http://tempo.monitoring.svc.cluster.local:3100`

> **Internal DNS**: Notice the URL format! Kubernetes internal DNS uses `<service-name>.<namespace>.svc.cluster.local`. This allows Grafana to securely talk to Loki and Tempo without going over the public internet.

### Step 6: Useful Monitoring Commands
Check on your monitoring tools:
```bash
kubectl get pods -n monitoring
kubectl top pods -n cloudpulse  # requires metrics-server to be installed
kubectl top nodes
```

### Step 7: Exposing Grafana to the Internet (Optional)
If you prefer not to port-forward, you can add an Ingress rule in your `ingress.yaml` to route to the Grafana service (e.g., at a `/grafana` path or a dedicated subdomain). However, for security, using port-forwarding is often preferred unless you configure secure authentication (like OAuth) for Grafana.
