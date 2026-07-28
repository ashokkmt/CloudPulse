# Deploying CloudPulse on DigitalOcean Kubernetes

This is a complete step-by-step guide to getting CloudPulse up and running on DOKS.

### Phase 1: Create the DOKS Cluster
- Go to DO Dashboard → Kubernetes → Create Cluster
- **Region**: Same as your VPC/databases.
- **VPC**: The existing `cloudpulse` VPC.
- **Node Pool**: 3 nodes, size `s-2vcpu-4gb`.
- **Name**: `cloudpulse-k8s`.

### Phase 2: Connect Your Local Machine
Authenticate and pull your new cluster's kubeconfig so you can control it via `kubectl`:
```bash
doctl auth init
doctl kubernetes cluster kubeconfig save cloudpulse-k8s
kubectl get nodes
```

### Phase 3: Grant Registry Access
Your cluster needs permission to pull images from your private DigitalOcean registry (`cloudpulse-registry`):
```bash
doctl kubernetes cluster registry add cloudpulse-k8s
```

### Phase 4: Understanding the Manifest Files
The deployment manifests live in `d:/projects/Digital Ocean Labs/k8s/`:

- `namespace.yaml` — Creates the `cloudpulse` and `monitoring` namespaces for logical separation.
- `secret.yaml` — Base64-encoded sensitive variables (`DATABASE_URL`, `REDIS_URL`, `SPACES_SECRET_ACCESS_KEY`). Encode locally with: `echo -n 'value' | base64`
- `configmap.yaml` — Plain text non-sensitive env vars (`PORT`, `APP_ENV`, `SPACES_ENDPOINT`, etc.).
- `backend.yaml` — Deployment with 3 replicas and a ClusterIP Service. Handles backend APIs (port 8000). Also includes resource limits, liveness probes, and uses `envFrom` to pull in ConfigMap/Secret.
- `worker.yaml` — Deployment with 1 replica for background processing. No Service is needed since it just reads from Redis/DB and doesn't take incoming requests.
- `frontend.yaml` — Deployment with 3 replicas and a ClusterIP Service. Serves the Next.js frontend (port 3000).
- `ingress.yaml` — The routing rules. Routes `/api` to the `backend-service` and `/` to the `frontend-service`. We leave the host field empty so it works with raw IPs (no domain needed).

### Phase 5: Install Nginx Ingress Controller
The Ingress Controller is what actually reads our `ingress.yaml` rules. Installing it creates a LoadBalancer Service that provisions a real DigitalOcean Load Balancer with a public IP. Then the Nginx pod routes external traffic into our cluster.
```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.11.3/deploy/static/provider/do/deploy.yaml
```

### Phase 6: Apply the Manifests
Order matters! You must create namespaces and configuration before deploying the apps.
```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/secrets.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/backend.yaml
kubectl apply -f k8s/worker.yaml
kubectl apply -f k8s/frontend.yaml
kubectl apply -f k8s/ingress.yaml
```

### Phase 7: Get Your Public IP
Once the Load Balancer is provisioned, you'll see a public IP:
```bash
kubectl get svc -n ingress-nginx
```
Wait for `EXTERNAL-IP`. Visit `http://<EXTERNAL-IP>` for the frontend and `http://<EXTERNAL-IP>/api/healthz` for the backend.

### Phase 8: Debugging

Critical commands for troubleshooting:
- `kubectl get pods -n cloudpulse` — Check pod status.
- `kubectl describe pod <name> -n cloudpulse` — See events (e.g., image pull errors, scheduling reasons, crash reasons).
- `kubectl logs <pod-name> -n cloudpulse` — View application logs.
- `kubectl logs <pod-name> -n cloudpulse --previous` — See logs from the LAST crashed container (critical for diagnosing CrashLoopBackOff).
- `kubectl get events -n cloudpulse --sort-by='.lastTimestamp'` — Timeline of cluster events.

### Common Issues & Fixes
1. **CrashLoopBackOff**: Pod keeps crashing and restarting. Check logs with `--previous`. Usually means missing env vars or a DB connection failure.
2. **ImagePullBackOff**: Wrong image name/tag, or registry not connected. Run `doctl kubernetes cluster registry add cloudpulse-k8s`. Check that your Git SHA tags match.
3. **Pending Pods**: Not enough resources on nodes. Check `kubectl describe pod` for CPU/Memory scheduling errors.
4. **Ingress not working**: Make sure `ingressClassName: nginx` is set in your `ingress.yaml`. Check `kubectl get ingress -n cloudpulse`.
5. **DB connection refused from K8s**: The DOKS cluster MUST be in the same VPC as your DO Managed Databases. You must also add the K8s cluster to the database's Trusted Sources in the DO dashboard.
