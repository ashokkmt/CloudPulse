# 🚀 11-Day Cloud Engineering Bootcamp

> **From zero cloud experience to production-ready infrastructure in 11 days.**
> Built around one application: **CloudPulse** — a real-time analytics and task management dashboard.

---

## Table of Contents

| Day | Topic | Key Services |
|-----|-------|-------------|
| — | [Overview & Architecture](#section-1-overview--application-architecture) | Application Design, Repo Structure, Cost Strategy |
| 1 | [Foundation: Local Dev + Docker Compose](#section-2-day-1---foundation-local-development--docker-compose) | Docker, Docker Compose |
| 2 | [First Cloud Deployment](#section-3-day-2---first-cloud-deployment-single-droplet) | Droplet, VPC, Firewall, DNS, SSL, Spaces, Registry |
| 3 | [Monitoring Stack](#section-4-day-3---monitoring-stack-prometheus-grafana-loki-tempo) | Prometheus, Grafana, Loki, Tempo, Alertmanager |
| 4 | [Managed Databases](#day-4---managed-databases-postgresql--redis) | Managed PostgreSQL, Managed Redis |
| 5 | [Object Storage](#day-5---object-storage-spaces--cdn--file-uploads) | Spaces, CDN |
| 6 | [CI/CD Pipeline](#day-6---cicd-pipeline-github-actions--container-registry) | GitHub Actions, Container Registry |
| 7 | [Horizontal Scaling](#day-7---horizontal-scaling-load-balancer--multi-droplet) | Load Balancer, Snapshots, Multi-Droplet |
| 8 | [Kubernetes](#day-8---kubernetes-digitalocean-kubernetes-doks) | DOKS, kubectl, Ingress, HPA |
| 9 | [Infrastructure as Code](#day-9---infrastructure-as-code-terraform) | Terraform, Remote State |
| 10 | [PaaS + Serverless](#day-10---platform-as-a-service-app-platform--serverless-functions) | App Platform, DO Functions |
| 11 | [Production Architecture](#day-11---production-architecture-security-disaster-recovery-and-final-review) | Security, DR, Final Review |
| — | [Final Architecture](#final-production-architecture) | Complete Production Diagram |
| — | [Service Reference](#complete-service-reference) | All DO Services Comparison |
| — | [Credit Summary](#credit-usage-summary) | Day-by-Day Cost Tracking |
| — | [What's Next](#whats-next) | Continued Learning |

---

# SECTION 1: OVERVIEW & APPLICATION ARCHITECTURE

## 1. Bootcamp Overview

Welcome to the Cloud Engineering Bootcamp! As a cloud architect, I've designed this 11-day program to take you from knowing the basics of Linux, Docker, Go, and Next.js, to mastering production-grade deployments on DigitalOcean. 

Our philosophy here is simple: **Learn by building, breaking, and fixing real things.** We aren't just going to spin up simple toy apps. We are going to deploy a multi-tier SaaS application, monitor it, break it, and scale it, applying production best practices at every step.

**What You'll Build:**
You'll be deploying **CloudPulse**, a real-time analytics and task management dashboard. It's a complete SaaS application that includes:
- A Next.js/TypeScript frontend for user interaction.
- A Go backend API handling business logic and worker tasks.
- PostgreSQL for persistent data.
- Redis for caching and session management.
- An Activity Simulator to generate continuous load, giving us real metrics to observe.

## 2. CloudPulse Application Architecture

Here is the high-level architecture of what we're building:

```text
                                     +-------------------+
                                     |                   |
                                     |   DigitalOcean    |
                                     |   Load Balancer   |
                                     |                   |
                                     +--------+----------+
                                              |
                                     +--------v----------+
                                     |                   |
                                     |   Nginx Reverse   |
                                     |   Proxy (Ingress) |
                                     |                   |
                                     +---+-----------+---+
                                         |           |
                           +-------------v-+       +-v-------------+
                           |               |       |               |
                           |  Next.js UI   |       |   Go API      |
                           |  (Frontend)   |       |  (Backend)    |
                           |               |       |               |
                           +---------------+       +---+---+-------+
                                                       |   |
                                 +---------------------+   +-------------------+
                                 |                                             |
                         +-------v-------+                             +-------v-------+
                         |               |                             |               |
                         |  Managed      |                             |  Managed      |
                         |  Redis        |                             |  PostgreSQL   |
                         |               |                             |               |
                         +---------------+                             +---------------+

                                      (Background Components)
                         +---------------+      +------------------+
                         |               |      |                  |
                         | Activity      |      | Background       |
                         | Simulator     +------> Worker (Go)      |
                         |               |      |                  |
                         +---------------+      +------------------+
```

## 3. Complete Repository Structure

Here is how your `cloudpulse` repository will be structured. Understanding where everything lives is the first step to owning the codebase.

```text
cloudpulse/
|-- frontend/              # Next.js app
|   |-- src/app/           # App router pages
|   |-- src/components/    # Reusable UI components
|   |-- src/lib/           # Utility functions, API clients
|   |-- public/            # Static assets (images, icons)
|   |-- Dockerfile         # Multi-stage build for frontend
|   |-- next.config.js     # Next.js configuration
|   +-- package.json       # Node dependencies
|-- backend/               # Go API
|   |-- cmd/
|   |   |-- api/           # Main API server entrypoint
|   |   |-- worker/        # Background worker entrypoint
|   |   +-- simulator/     # Activity simulator entrypoint
|   |-- internal/          # Private application code
|   |   |-- handler/       # HTTP handlers (controllers)
|   |   |-- middleware/    # Auth, logging, metrics middlewares
|   |   |-- model/         # Structs and DB schema definitions
|   |   |-- repository/    # DB and Redis interaction logic
|   |   |-- service/       # Core business logic
|   |   +-- metrics/       # Prometheus metrics definitions
|   |-- migrations/        # SQL migration files
|   |-- Dockerfile         # Multi-stage build for Go binaries
|   +-- go.mod             # Go dependencies
|-- monitoring/            # Observability stack config
|   |-- prometheus/        # Prometheus scrape configs
|   |-- grafana/           # Grafana configs
|   |   |-- provisioning/  # Auto-load datasources/dashboards
|   |   +-- dashboards/    # JSON dashboard definitions
|   |-- loki/              # Log aggregation config
|   |-- tempo/             # Distributed tracing config
|   +-- alertmanager/      # Alert routing and notification config
|-- infrastructure/        # Infrastructure as Code (Day 6+)
|   |-- terraform/         # DO resource definitions
|   |   |-- modules/       # Reusable TF modules
|   |   +-- environments/  # Dev/Prod state
|   |-- kubernetes/        # K8s manifests (Day 8+)
|   |   |-- base/          # Common K8s configs
|   |   +-- overlays/      # Env-specific configs (dev, prod)
|   +-- scripts/           # Helper bash scripts
|-- .github/               # CI/CD pipelines
|   +-- workflows/         # GitHub Actions definitions
|-- nginx/                 # Reverse proxy config
|   +-- nginx.conf         # Routing rules for backend/frontend
|-- docker-compose.yml             # Local development stack
|-- docker-compose.monitoring.yml  # Local observability stack
|-- docker-compose.prod.yml        # Production single-node stack
|-- docs/                  # Architecture docs, runbooks, ADRs
|   |-- architecture.md   # System architecture documentation
|   |-- runbook.md         # Operational procedures
|   +-- decisions/         # Architecture Decision Records
|-- Makefile               # Handy shortcuts for commands
+-- README.md              # Project documentation
```

*Senior Tip:* Keeping `cmd/` separate from `internal/` in Go is a standard best practice. It ensures your business logic (`internal/`) can be shared easily among the API server, background worker, and simulator without circular dependencies.

## 4. Resource Provisioning Strategy
> **Philosophy:** Use production-realistic sizes. The goal is **maximum hands-on time** with every service. Don't destroy resources between days unless we're done learning from them.

**Philosophy:** Use production-realistic sizes. Don't destroy resources between days unless we're done learning from them. The goal is **maximum hands-on time** with every service.

**Always-On Resources (Day 2 onward):**
- App Droplet: `s-4vcpu-8gb`  — room for app + DB containers early, then app-only later
- Monitoring Droplet: `s-2vcpu-4gb`  — dedicated Prometheus/Grafana/Loki/Tempo (Day 3+)
- Container Registry Professional:  — unlimited repos, 100 GiB storage
- Spaces Bucket: 
- VPC/Firewall/Project/Domain: Free

**Created and Kept Running Until Done:**
- Managed PostgreSQL `db-s-1vcpu-2gb`  — Day 4 onward, kept alive
- Managed Redis `db-s-1vcpu-2gb`  — Day 4 onward, kept alive
- Load Balancer  — Day 7 onward, kept alive
- Second App Droplet `s-4vcpu-8gb`  — Day 7 onward, kept alive
- DOKS Cluster 3×`s-2vcpu-4gb`  — Day 8 through Day 10
- App Platform Professional  — Day 10 through Day 11

- Try a HA PostgreSQL cluster  for a few hours to test failover
- Spin up a GPU Droplet briefly to see what's available
- Deploy to a second region to test cross-region latency

## 5. Daily Schedule Overview

- **Day 1: Foundation.** Local Development + Docker Compose  ✅ COMPLETED
- **Day 2: First Cloud Deployment.** App Droplet (8GB) + Nginx + SSL + VPC + Firewall 
- **Day 3: Observability.** Dedicated Monitoring Droplet + Full Stack 
- **Day 4: Managed Databases.** Managed PostgreSQL + Redis (2GB nodes) + HA experiment 
- **Day 5: Object Storage.** Spaces + CDN + File Uploads + Backups 
- **Day 6: CI/CD Pipeline.** GitHub Actions + Professional Registry + Auto-Deploy 
- **Day 7: Horizontal Scaling.** Load Balancer + Second Droplet (8GB) + Failover 
- **Day 8: Kubernetes.** DOKS Cluster (3 nodes) + Keep Droplets running for comparison 
- **Day 9: Infrastructure as Code.** Terraform all resources, keep K8s alive 
- **Day 10: PaaS + Serverless.** App Platform Pro + Functions + Destroy K8s 
- **Day 11: Production Architecture.** Security + DR + HA Database Experiment + Final Review 

---

# SECTION 2: DAY 1 - Foundation: Local Development + Docker Compose

## Objective
Build the CloudPulse application locally and containerize everything with Docker Compose. This day focuses strictly on local development to ensure the architecture is sound before touching the cloud.

## Resources to Create
- **DigitalOcean Project:** Create a new project called `CloudPulse` in the DO dashboard. This is just an organizational folder .
- **

## Local Architecture Diagram

```text
                      +-------------------------+
                      |      localhost:8080     |
                      +-------------------------+
                                  |
                +-----------------v-----------------+
                |                                   |
                |          docker-compose           |
                |                                   |
                |  +-------------+   +-----------+  |
                |  | Frontend    |   | Backend   |  |
                |  | (Next.js)   |   | (Go)      |  |
                |  | :3000       |   | :8000     |  |
                |  +-------------+   +-----------+  |
                |                          |        |
                |  +-------------+   +-----v-----+  |
                |  | Activity    |   | Worker    |  |
                |  | Simulator   |   | (Go)      |  |
                |  +------|------+   +-----|-----+  |
                |         |                |        |
                |  +------v------+   +-----v-----+  |
                |  | PostgreSQL  |   | Redis     |  |
                |  | :5432       |   | :6379     |  |
                |  +-------------+   +-----------+  |
                +-----------------------------------+
```

## Detailed Implementation Steps

### 1. Create the GitHub Repository
Initialize the repository locally and map out the folder structure defined in Section 1. 

*Why?* Getting the structure right from Day 1 prevents messy refactors later.

### 2. Set up the Go Backend
Write the core API in Go using standard `net/http` or a router like `chi`.

**API Specification:**
- `GET /healthz` - Liveness probe (returns 200 OK)
- `GET /metrics` - Prometheus metrics exposure
- `POST /api/auth/register` - Register user (hashes password with bcrypt)
- `POST /api/auth/login` - Authenticate and return JWT
- `GET /api/tasks` - List tasks (requires JWT)
- `POST /api/tasks` - Create task
- `PUT /api/tasks/{id}` - Update task
- `DELETE /api/tasks/{id}` - Delete task
- `GET /api/analytics` - Return dashboard metrics (cached in Redis)

*Senior Tip:* Don't build your own JWT implementation from scratch. Use `golang-jwt/jwt`. Always hash passwords using `golang.org/x/crypto/bcrypt`.

### 3. Set up Next.js Frontend
Bootstrap a Next.js application (App Router) to interact with the API.
- **Pages:** Login, Register, Dashboard (Task List + Analytics charts).
- Use `fetch` to talk to the backend, attaching the JWT to the `Authorization: Bearer <token>` header.

### 4. Create Dockerfiles
Create highly optimized, multi-stage `Dockerfile`s.

*Backend Dockerfile snippet:*
```dockerfile
# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o api ./cmd/api

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api .
EXPOSE 8000
CMD ["./api"]
```
*Why multi-stage?* It keeps the final image tiny (~15MB instead of 800MB), which speeds up deployments and reduces the attack surface area.

### 5. Create `docker-compose.yml`
Define your entire stack.
Include `postgres:15-alpine` and `redis:7-alpine`. Mount volumes for DB data so it persists across container restarts!

### 6. Describe the Activity Simulator
Write a simple Go script in `cmd/simulator` that authenticates as a dummy user and loops infinitely, making random GET, POST, and PUT requests to the API. 
*Why?* A dashboard is boring without data. The simulator will give us actual traffic to monitor in Day 3.

### 7. Test Everything Locally
Run `docker-compose up --build`. Verify that the frontend loads at `localhost:3000`, the API responds at `localhost:8000`, and the database migrations run successfully.

## Exercises
1. Modify the Go backend to return a new metric: total users registered.
2. Intentionally introduce a syntax error in the Next.js code and see how Docker Compose reacts (it shouldn't crash the container, just show an error in the UI).
3. Connect to the PostgreSQL container using `psql` or DBeaver locally to inspect the tables.
4. Use `curl` to hit the `/api/auth/login` endpoint directly and decode the resulting JWT using jwt.io.
5. Scale the worker container to 3 replicas in `docker-compose.yml` and observe the logs.

## Experiments
- **Break things:** Stop the Redis container (`docker stop <redis_container_id>`). What happens to the frontend? Does the backend crash or handle it gracefully? (It *should* degrade gracefully, but likely won't on your first try!).
- **Observe behavior:** Run the Activity Simulator and watch the PostgreSQL container's CPU usage spike using `docker stats`.

## Checklist
- [ ] DigitalOcean Project created.
- [ ] Repo structure matches the spec.
- [ ] Go API responds to `/healthz`.
- [ ] Next.js UI loads locally.
- [ ] Multi-stage Dockerfiles written.
- [ ] `docker-compose up` brings up all 5 services successfully.
- [ ] Simulator generates traffic.

## What I Learned
- Multi-stage Docker builds reduce image size significantly.
- Local networking in Docker Compose automatically resolves service names (e.g., the backend can connect to `postgres:5432` without knowing its IP).
- Designing API routes upfront saves time.

## Artifacts Produced
- `docker-compose.yml`
- Backend and Frontend codebases with Dockerfiles.

---

# SECTION 3: DAY 2 - First Cloud Deployment: Single Droplet

## Objective
Deploy CloudPulse to a single DigitalOcean Droplet using Docker Compose and Nginx as a reverse proxy. This is the "old school" but incredibly reliable way to host an application.

## Resources to Create
- **VPC:** Custom Virtual Private Cloud network .
- **Firewall:** Cloud-level firewall rules .
- **SSH Key:** For secure access.
- **Droplet:** `s-4vcpu-8gb` running Ubuntu 24.04 .
- **Container Registry:** Professional tier .
- **Domain + DNS:** Add your domain to DO, set an A record.

**

## Architecture Diagram

```text
                                  +-------------------+
                                  |    Internet       |
                                  +---------+---------+
                                            |
                                (Firewall: 80, 443, 22)
                                            |
+-------------------------------------------v-----------------------------------------+
| DO VPC (10.106.0.0/20)                                                              |
|                                                                                     |
|   +-----------------------------------------------------------------------------+   |
|   | Droplet (Ubuntu 24.04)                                                      |   |
|   |                                                                             |   |
|   |   +-----------------+                                                       |   |
|   |   |      Nginx      | (Reverse Proxy, SSL Termination)                      |   |
|   |   |   (Port 80/443) |                                                       |   |
|   |   +--------+--------+                                                       |   |
|   |            |                                                                |   |
|   |    +-------v-------+--------------------+--------------------+              |   |
|   |    |               |                    |                    |              |   |
|   | +--v-----+      +--v-----+         +----v-----+         +----v-----+        |   |
|   | | Frontend|     | Backend |         | Postgres |         |  Redis   |        |   |
|   | | (:3000) |     | (:8000) |         | (:5432)  |         | (:6379)  |        |   |
|   | +---------+     +---------+         +----------+         +----------+        |   |
|   |                                                                             |   |
|   +-----------------------------------------------------------------------------+   |
+-------------------------------------------------------------------------------------+
```

## Detailed Implementation Steps

### 1. Create VPC
Create a VPC in your chosen region (e.g., `nyc3`). Name it `cloudpulse-vpc`. CIDR: `10.106.0.0/20`. 
*Why?* Placing resources in a custom VPC isolates them at the network layer. Later, when we add managed databases, they will communicate over this secure private network rather than the public internet.

### 2. Create Firewall
Create a DO Firewall named `cloudpulse-fw`.
- **Inbound:** 
  - SSH (TCP 22) from your IP only (Security best practice).
  - HTTP (TCP 80) from All IPv4/IPv6.
  - HTTPS (TCP 443) from All IPv4/IPv6.
- **Outbound:** All traffic permitted.

> **Cloud Firewalls vs `ufw`:** DigitalOcean Cloud Firewalls drop traffic *before* it reaches your Droplet's OS — they run on the hypervisor. The host-based firewall `ufw` (Uncomplicated Firewall) runs *inside* your VM, so malicious packets already consumed bandwidth to reach your kernel. In production, use both layers: DO Cloud Firewall as the primary defense, and `ufw` as a secondary defense-in-depth layer. Practice both:
> ```bash
> # On the Droplet (after SSH), set up ufw as a secondary layer
> sudo ufw default deny incoming
> sudo ufw default allow outgoing
> sudo ufw allow 22/tcp    # SSH
> sudo ufw allow 80/tcp    # HTTP
> sudo ufw allow 443/tcp   # HTTPS
> sudo ufw enable
> sudo ufw status verbose
> ```
> *When to rely on `ufw` alone:* Only on bare-metal or non-cloud providers that don't offer external firewalls.

### 3. Generate SSH Key
Run `ssh-keygen -t ed25519 -C "your_email@example.com"`. Upload the public key to DO.
*Why ed25519?* It's faster and more secure than RSA. 

### 4. Create Droplet
- **Region:** Same as VPC (e.g., `nyc3`).
- **Image:** Ubuntu 24.04 LTS.
- **Size:** Basic Premium Intel `s-4vcpu-8gb` . *Why 8GB?* We have credits to burn — use a production-realistic size. Running Go, Node.js, Postgres, Redis, the full monitoring stack, and an OS comfortably needs 8GB. On Day 3 we'll offload monitoring to its own Droplet, but until then this gives breathing room.
- **VPC:** Select `cloudpulse-vpc`.
- **Authentication:** SSH Key.
- **Metrics:** Enable DO Monitoring.
- **Tags:** `env:prod`, `app:cloudpulse`.

### 5. SSH into Droplet and Create a Deploy User
```bash
# SSH in as root (first time only)
ssh root@<droplet_ip>

# Create a non-root deploy user (NEVER run production apps as root)
adduser deploy
usermod -aG sudo deploy

# Copy your SSH key to the new user
mkdir -p /home/deploy/.ssh
cp ~/.ssh/authorized_keys /home/deploy/.ssh/
chown -R deploy:deploy /home/deploy/.ssh
chmod 700 /home/deploy/.ssh
chmod 600 /home/deploy/.ssh/authorized_keys

# Test: open a new terminal and SSH as deploy
# ssh deploy@<droplet_ip>
```
*Why a deploy user?* Running as `root` means any compromised container or misconfigured service has full system access. A dedicated `deploy` user with `sudo` gives you admin capabilities when needed while limiting the blast radius of mistakes. This is Linux Administration 101.

> **Linux User Management Quick Reference:**
> - `adduser <name>` — Create user with home directory
> - `usermod -aG sudo <name>` — Grant sudo privileges
> - `su - <name>` — Switch to user
> - `id <name>` — Show user's groups
> - `passwd -l root` — Lock root password (after setting up deploy user)

### 6. Install Docker + Docker Compose
SSH in as your `deploy` user. Follow the official Docker `apt` repository instructions:
```bash
# Update package index and install prerequisites
sudo apt update && sudo apt install -y ca-certificates curl gnupg

# Add Docker's official GPG key and apt repo
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo $VERSION_CODENAME) stable" | sudo tee /etc/apt/sources.list.d/docker.list

# Install Docker
sudo apt update && sudo apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# Add deploy user to docker group (no sudo needed for docker commands)
sudo usermod -aG docker deploy
newgrp docker
```
*Why `apt` instead of `snap`?* The `apt` package gives you more control over Docker versions and update timing — critical in production.

### 7. Push to Container Registry
Instead of building images on the Droplet (which uses CPU/Memory), build them locally or in CI and push to DO Container Registry (DOCR).
`docker tag cloudpulse-backend registry.digitalocean.com/cloudpulse/backend:v1`
`docker push registry.digitalocean.com/cloudpulse/backend:v1`

### 8. Create `docker-compose.prod.yml`
Create a prod-specific compose file on the Droplet.
- *Differences from dev:* Use DOCR image URLs instead of `build: .`. Set `restart: always` on all services so they survive reboots. Pass production secrets via `.env`.

### 9. Set up Nginx Reverse Proxy
Install Nginx directly on the host (`apt install nginx`).
Configure it to route traffic:
```nginx
server {
    listen 80;
    server_name _;

    location /api/ {
        proxy_pass http://127.0.0.1:8000;

        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location / {
        proxy_pass http://127.0.0.1:3000;

        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```
*Why Nginx?* It acts as a gatekeeper, handles SSL efficiently, and maps your domain to the correct internal Docker ports.

### 10. Configure Domain and DNS
Point your domain's A record to the Droplet's public IP in the DO Networking tab.

### 11. Set up SSL with Let's Encrypt
Run `certbot --nginx -d yourdomain.com`.
*Why?* Never deploy a production app without HTTPS. Certbot automates SSL certificate provisioning and renewal.

### 12. Deploy and Verify
Run `docker-compose -f docker-compose.prod.yml up -d`. Check your domain in the browser!

## Exercises
1. Check Nginx access logs (`tail -f /var/log/nginx/access.log`) while clicking around the app.
2. View Docker logs for the backend container (`docker logs -f <backend_container>`).
3. Scale the simulator locally, push a new image, and run a rolling update via docker-compose on the droplet.
4. Verify the Let's Encrypt cron job is installed (`systemctl list-timers | grep certbot`).
5. Use the Cloud Console to look at the Droplet's CPU utilization graphs.

## Experiments
- **Kill a container:** Run `docker kill <frontend_container>`. Watch how `restart: always` brings it back up.
- **Test firewall rules:** Try to access Postgres directly from your local machine (e.g., `psql -h <droplet_ip> -p 5432 -U user`). It should hang and timeout because port 5432 is blocked by the DO Firewall.
- **Reboot Droplet:** Run `reboot` in the SSH session. Wait 1 minute. Verify the website comes back online automatically without manual intervention.

## Checklist
- [ ] Custom VPC and Firewall created.
- [ ] Droplet provisioned and secured with SSH key.
- [ ] Docker installed on Droplet.
- [ ] Images pushed to DO Container Registry.
- [ ] `docker-compose.prod.yml` running on Droplet.
- [ ] Nginx configured as reverse proxy.
- [ ] Domain mapped via DNS.
- [ ] HTTPS enabled via Let's Encrypt.
- [ ] Droplet survives a reboot.

## What I Learned
- Cloud firewalls are superior to host-based firewalls for edge security.
- Nginx is incredibly powerful for routing subpaths (`/api`) to different backend services.
- Building images externally and pulling them to the server saves the server from CPU spikes.

## Artifacts Produced
- `docker-compose.prod.yml`
- `/etc/nginx/sites-available/cloudpulse`
- Production SSL Certificates.

---

# SECTION 4: DAY 3 - Monitoring Stack: Prometheus, Grafana, Loki, Tempo

## Objective
Deploy a complete observability stack alongside CloudPulse. Running code is easy; knowing *what* your code is doing is hard. We will implement metrics, logs, and traces (The Three Pillars of Observability).

## Resources to Create
- **New Droplet:** `s-2vcpu-4gb`  dedicated to the monitoring stack.
- **

> **Why a separate Monitoring Droplet?** In production, you *never* run your monitoring on the same server as your application. If your app crashes and takes down the server, you lose visibility at the exact moment you need it most. A dedicated monitoring node ensures Grafana/Prometheus stay alive even when your app servers are on fire. With credits expiring, this is exactly the kind of production-realistic setup we should practice.

> **Create the droplet and a monitor user because production containers should not run as root user. Take help from Day 2 section to create new user named `monitor` and setup docker too**

## Observability Architecture Diagram

```text
+---------------------------------------------------------------------------------+
| Droplet                                                                         |
|                                                                                 |
|  +-------------+       +---------------+       +-------------+                  |
|  |             | Logs  |               |       |             |                  |
|  | Docker Apps +------->  Loki         <-------+             |                  |
|  | (CloudPulse)|       |               |       |             |                  |
|  +------+------+       +---------------+       |             |                  |
|         |                                      |   Grafana   |-----> Browser UI |
|         | Traces       +---------------+       |             |                  |
|         +-------------->  Tempo        <-------+             |                  |
|         |              |               |       |             |                  |
|         | Metrics      +---------------+       +------^------+                  |
|  +------v------+                                      |                         |
|  | Prometheus  |       +---------------+              |                         |
|  | Scraper     <-------+  Prometheus   +--------------+                         |
|  +------+------+       |  DB           |                                        |
|         |              +-------+-------+                                        |
|         |                      |                                                |
|  +------v------+       +-------v-------+                                        |
|  | Node Exporter       | Alertmanager  |-----> Email/Slack                      |
|  | cAdvisor    |       |               |                                        |
|  +-------------+       +---------------+                                        |
|                                                                                 |
+---------------------------------------------------------------------------------+
```

## Implementation Steps

### 1. Create `docker-compose.monitoring.yml`
Create a separate compose file for the monitoring stack. This keeps things modular. 
Include: `prometheus`, `grafana`, `loki`, `promtail` (or configure Docker logging driver directly), `tempo`, `alertmanager`, `node-exporter`, and `cadvisor`.

### 2. Configure Prometheus
Create `prometheus.yml`. 
*Why?* Prometheus uses a pull model. You must tell it where to scrape metrics.
Configure jobs for:
- `app` (Your Go API `/metrics` endpoint)
- `node` (Node Exporter for host metrics)
- `cadvisor` (Docker container metrics)

### 3. Configure Grafana
Mount a `provisioning/` folder to automatically load datasources (Prometheus, Loki, Tempo) and your custom dashboards on startup. This prevents you from having to click through the UI every time you deploy.

### 4. Configure Loki
Create `loki-config.yml`. Loki is "Prometheus, but for logs." It indexes metadata (labels) rather than the full log text, making it highly efficient.

### 5. Configure Tempo
Create `tempo-config.yml`. Tempo handles distributed tracing via the OpenTelemetry (OTLP) protocol.

### 6. Configure Alertmanager
Create `alertmanager.yml`. Define routing to send critical alerts to a webhook (like Discord or Slack). Add a `rules.yml` to Prometheus to trigger these alerts.

### 7. Deploy Node Exporter
Runs on port 9100. It exposes OS-level metrics (CPU, disk I/O, memory, network).

### 8. Deploy cAdvisor
Runs on port 8080. It exposes resource usage and performance characteristics of running Docker containers.

### 9. Configure Docker Logging Driver
Modify your app's `docker-compose.yml` to send logs directly to Loki:
```yaml
logging:
  driver: loki
  options:
    loki-url: "http://localhost:3100/loki/api/v1/push"
```
*Why?* This bypasses the need for Promtail and streams `stdout/stderr` natively.

### 10. Add OpenTelemetry Tracing to Go Backend
Instrument your Go code using the OpenTelemetry SDK. Wrap HTTP handlers and database calls with spans. Point the exporter to the Tempo container.

### 11. Create Grafana Dashboards
Build and export JSON dashboards:
- **Application Dashboard:** Request rates (RPS), P95 Latency, Error rates (HTTP 5xx).
- **Infrastructure Dashboard:** Node CPU, memory utilization, disk space, network tx/rx.
- **Database Dashboard:** PostgreSQL active connections, slow query counts.
- **Redis Dashboard:** Cache hit rate, memory usage, command execution rate.
- **Business Dashboard:** Active users, tasks created per minute (extracted from logs/metrics).

### 12. Set Up Alerts
Define alerting rules in Prometheus:
- **InstanceDown:** If Node Exporter is unreachable for 1m.
- **HighCPUUsage:** If CPU > 85% for 5m.
- **HighErrorRate:** If HTTP 5xx errors > 5% of total requests for 2m.

### 13. Verify Everything Works
Deploy the monitoring stack (`docker-compose -f docker-compose.monitoring.yml up -d`). Open Grafana, check the dashboards, and ensure data is flowing.

## Grafana Dashboard Specifications
- **App Rate/Errors:** `rate(http_requests_total[5m])` and `rate(http_requests_total{status=~"5.."}[5m])`.
- **Latency:** `histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))`.
- **Node Memory:** `100 - ((node_memory_MemAvailable_bytes * 100) / node_memory_MemTotal_bytes)`.

## Alert Rules
- **Thresholds:** Always use duration conditions (e.g., `for: 5m`) in alerts to prevent flapping (alerts firing and resolving rapidly due to micro-spikes).

## Exercises
1. Log into Grafana (default admin/admin).
2. Write a PromQL query to find the container using the most memory.
3. Write a LogQL query in the Explore tab to find all Go errors containing the word "timeout".
4. Add a new span to a specific function in your Go code and view it in Tempo.
5. Create a new panel on the Business Dashboard tracking the exact number of file uploads.

## Experiments
- **Generate Load:** Run your Activity Simulator at a high rate. Watch the Application Dashboard. Do you see the RPS climb?
- **Trigger an Alert:** SSH into the droplet and run `stress --cpu 4`. Wait 5 minutes. Did Alertmanager fire the HighCPU alert?
- **Trace a Request:** Click on a slow request in the Grafana logs. Use the attached Trace ID to jump directly to Tempo and see exactly which database query slowed it down.
- **Kill Redis:** Stop the Redis container. Watch the Error Rate spike in Grafana.
- **Search Logs:** Filter Loki logs for `{container="cloudpulse-backend"} |= "panic"`.

## Checklist
- [ ] `docker-compose.monitoring.yml` running.
- [ ] Prometheus scraping app, node, and cadvisor successfully.
- [ ] Docker configured to send logs to Loki.
- [ ] Go backend instrumented with OpenTelemetry.
- [ ] 5 Dashboards provisioned in Grafana.
- [ ] Alertmanager rules configured and tested.

## What I Learned
- Collecting data is useless without good dashboards to visualize it.
- Correlation is key: having logs, metrics, and traces in one tool (Grafana) allows you to find a spike in metrics, check the logs for errors, and trace the specific request that failed.
- Setting up observability early pays massive dividends when debugging production issues later.

## Artifacts Produced
- `docker-compose.monitoring.yml`
- Prometheus/Loki/Tempo/Alertmanager configurations.
- Grafana Dashboard JSON files.
# DAY 4 - Managed Databases: PostgreSQL + Redis

## Objective
Migrate from containerized databases to DigitalOcean Managed PostgreSQL and Managed Redis. Learn about managed database operations, connection pooling, backups, and failover.

## WHY Managed Databases?
Running databases in containers (as we did in Days 1-3) is great for development and learning, but it's risky for production. Here is why we are migrating to Managed Databases:

1. **Automated Backups & Point-in-Time Recovery (PITR)**: DO handles automated daily backups and write-ahead log (WAL) archiving, allowing you to restore data to any point within the last 7 days. Building this yourself is complex and error-prone.
2. **High Availability & Failover**: If a node goes down, DO automatically promotes a standby node. You don't have to wake up at 3 AM to fix a broken replication setup.
3. **Automated Patching**: Security patches and minor version updates are applied automatically during your specified maintenance window.
4. **Built-in Connection Pooling (PgBouncer)**: Postgres creates a new OS process per connection, which is memory-heavy. DO provides a built-in connection pooler to handle thousands of connections efficiently.
5. **Monitoring & Alerting**: Native integration with DO's monitoring dashboard for slow queries, IOPS, and memory usage without configuring custom exporters.
6. **Focus**: Your job is to build the application, not to be a DBA.

## Resources to Create & Destroy

**Create:**
- Managed PostgreSQL Cluster (db-s-1vcpu-2gb, same region as Droplet, assigned to VPC)
- Managed Redis Cluster (db-s-1vcpu-2gb, same region as Droplet, assigned to VPC)

**Destroy:**
- PostgreSQL Docker container (remove from `docker-compose.prod.yml`)
- Redis Docker container (remove from `docker-compose.prod.yml`)

## Estimated Cost Change
- Droplet (s-4vcpu-8gb): 
- Registry (Basic): 
- Managed PostgreSQL (db-s-1vcpu-2gb): +
- Managed Redis (db-s-1vcpu-2gb): +
- **New Running Total**: ~

## Architecture Diagram

```ascii
                      DigitalOcean Cloud
+-------------------------------------------------------------+
|                            VPC                              |
|                                                             |
|  +-------------------+        +--------------------------+  |
|  |     Droplet       |        |    Managed Databases     |  |
|  |                   |        |                          |  |
|  |  [ Nginx ]        |        |  +--------------------+  |  |
|  |      |            |        |  | Managed PostgreSQL |  |  |
|  |      v            |=======>|  | (db-s-1vcpu-2gb)   |  |  |
|  |  [ Go Backend ]   |        |  +--------------------+  |  |
|  |                   |        |                          |  |
|  |  [ Next.js ]      |        |  +--------------------+  |  |
|  |                   |=======>|  | Managed Redis      |  |  |
|  |  [ Monitoring ]   |        |  | (db-s-1vcpu-2gb)   |  |  |
|  +-------------------+        |  +--------------------+  |  |
|                               +--------------------------+  |
+-------------------------------------------------------------+
```

## Implementation Steps

### 1. Create Managed PostgreSQL Cluster
Navigate to the Cloud Provider Control Panel -> Databases -> Create Database Cluster.
- **Database Engine**: PostgreSQL (choose the latest stable major version, e.g., 15 or 16).
- **Node Plan**: Basic `db-s-1vcpu-2gb` .
- **Choose a Datacenter**: Must be the *exact same region* as your Droplet to minimize latency and ensure it can join the VPC.
- **VPC Network**: Select your existing VPC.
- **Cluster Name**: `cloudpulse-pg-prod`

> [!IMPORTANT]
> **Connection Pooling**: PostgreSQL handles connections poorly. Each connection consumes ~10MB RAM. With a 1GB RAM database, 100 connections = OOM kill. DO includes PgBouncer. You will set up a pool in the DO dashboard.
> - **Transaction mode (Default/Recommended)**: Pool reclaims the connection after each transaction. Best for typical web apps.
> - **Session mode**: Pool reclaims connection after the client disconnects. Use if your app relies on session-level state (like prepared statements or advisory locks).
> - **Statement mode**: Reclaims after each statement. Rarely used.

**Trusted Sources**: Once created, go to the DB settings and restrict access under "Trusted Sources". Add your Droplet or VPC. This acts as a firewall, preventing external internet access to your DB.

**Maintenance Window**: Set this to a low-traffic time (e.g., Sunday 3:00 AM UTC) for automated updates.

### 2. Create Managed Redis Cluster
Navigate to Databases -> Create Database Cluster.
- **Database Engine**: Redis
- **Node Plan**: Basic `db-s-1vcpu-2gb` .
- **Datacenter**: Same region as Droplet.
- **VPC Network**: Existing VPC.
- **Cluster Name**: `cloudpulse-redis-prod`

> [!TIP]
> **Eviction Policies**: By default, DO Redis uses `noeviction` (returns errors if memory is full) or `allkeys-lru`. If you use Redis *strictly* as a cache, set the eviction policy to `allkeys-lru` in the settings so it automatically drops old keys when full. If you use it for job queues (like Sidekiq/Celery), you might want `volatile-lru` or `noeviction` to prevent data loss.

### 3. Migration Strategy (Zero Data Loss)
We need to move data from the Docker containers to the Managed DBs.

1. **Export data from Docker PG**:
   SSH into your droplet.
   ```bash
   docker exec -t cloudpulse-db pg_dumpall -c -U myuser > dump.sql
   ```
2. **Import to Managed PG**:
   Get the connection string from the Cloud Console. Ensure you use the *VPC private network* string (usually starts with `private-db-postgresql...`).
   ```bash
   # Install postgresql-client on the droplet if needed
   sudo apt-get install postgresql-client
   psql "postgres://doadmin:PASSWORD@private-db-postgresql-xyz.db.ondigitalocean.com:25060/defaultdb?sslmode=require" < dump.sql
   ```
3. **Verify Data Integrity**:
   Connect via `psql` and run `SELECT count(*) FROM users;` to ensure data migrated.
4. **Redis**: Since our Redis is just cache/session data, we will let it start fresh. If it had critical data, we would use a tool like `redis-dump` or `RIOT`.

### 4. Update docker-compose.prod.yml
Remove the local databases.

```diff
-  db:
-    image: postgres:15-alpine
-    ...
-  redis:
-    image: redis:7-alpine
-    ...
```

Update your backend service environment variables:
```yaml
  backend:
    environment:
      # Use the connection pooler port (usually 25061), NOT the direct port 
      - DATABASE_URL=postgres://doadmin:PASSWORD@private-db-postgresql-xyz.db.ondigitalocean.com:25061/defaultdb?sslmode=require
      - REDIS_URL=rediss://doadmin:PASSWORD@private-db-redis-xyz.db.ondigitalocean.com:25061
```
*Note the `rediss://` (with two s's) which indicates TLS/SSL connection for Redis, required by DO.*

### 5. Test Application
Run `docker-compose -f docker-compose.prod.yml up -d` to restart the app. Verify it connects, loads data, and sessions work.

### 6. Update Monitoring
Your Prometheus `docker-compose` previously scraped the Docker PG and Redis exporters. You now need to update this. DigitalOcean provides metrics via an API, but for Prometheus integration, you often rely on DO's built-in dashboard for DB metrics, OR you can run PG/Redis exporters on your Droplet pointing at the Managed DBs.

> [!WARNING]
> **Common Mistake**: Connecting to the public connection string from your Droplet. This routes traffic out to the internet and back, incurring latency and egress bandwidth costs. Always use the VPC `private-...` connection string.

## Security Best Practices
- **Never expose DBs to the internet**: Rely completely on Trusted Sources.
- **SSL**: DO requires SSL (`sslmode=require`) for connections.
- **Rotation**: Rotate database passwords periodically via the DO dashboard.

## Exercises
- [ ] **Slow Query Log**: In DO dashboard -> Logs & Queries. Run `SELECT pg_sleep;` in psql, then find it in the dashboard.
- [ ] **Connection Pooling**: Run `pgbench` against the direct port  with 200 clients. Watch it fail. Then run it against the pool port  and watch it succeed.
- [ ] **Read Replica**: Create a read replica in the DO dashboard. Connect a `psql` session to it and verify you can `SELECT` but not `INSERT`. (Delete it after to save money).
- [ ] **Backup Restoration**: Create a new database from yesterday's backup using the DO UI.

## Experiments
- **Simulate connection exhaustion**: Write a simple Go script that opens 1000 connections in a loop and sleeps. Watch how PgBouncer handles the queue versus direct connections.
- **Kill app connections**: Use `SELECT pg_terminate_backend(pid) FROM pg_stat_activity;` to kill your backend's connections. Observe your Go app automatically reconnecting (assuming you configured your connection pool correctly in Go).
- **Redis memory limit**: Write a script to insert 2GB of random strings into your 1GB Redis instance. Observe the eviction policy in action.

## Checklist
- [ ] Created Managed Postgres in same VPC
- [ ] Created Managed Redis in same VPC
- [ ] Migrated data using pg_dump/psql
- [ ] Configured Trusted Sources
- [ ] Updated application to use connection pooler URL
- [ ] Removed containerized databases from docker-compose
- [ ] Verified application functionality

## What I Learned
- **Self-hosted vs Managed Tradeoffs**: Self-hosting saves money but costs time and risk. Managed DBs offer peace of mind, automated backups, and built-in connection pooling, but come at a premium.
- Connection pooling is absolutely critical for PostgreSQL in memory-constrained environments.
- VPC networking is essential for security and performance when connecting Droplets to Managed Services.

---

# DAY 5 - Object Storage: Spaces + CDN + File Uploads

## Objective
Integrate DigitalOcean Spaces (S3-compatible object storage) for file uploads, static assets, and backups. Configure CDN for performance.

## WHY Object Storage?
Why not just save user uploads to the Droplet's local disk?
1. **Durability**: Droplet disks are ephemeral relative to block storage. Spaces are highly durable (replicated across multiple servers).
2. **Statelessness**: If files are on the Droplet, you can't easily scale to multiple Droplets (Load Balancing) because user A's avatar is on Droplet 1, but they might be routed to Droplet 2. Object storage makes your app stateless.
3. **CDN Integration**: Spaces has a built-in CDN. Serving images directly from the CDN edge nodes is vastly faster and reduces load on your Droplet.
4. **Cost**: Spaces is  for 250GB, which is much cheaper than resizing your Droplet just for disk space.

## Resources to Create
- Spaces Bucket (Created Day 2)
- Spaces API Keys
- Spaces CDN Endpoint (Custom Subdomain)

## Estimated Cost Change
- Spaces:  (already factored in from Day 2 creation, but now utilized)
- **New Running Total**: ~

## Architecture Diagram

```ascii
                      User / Browser
                             |
         +-------------------+-------------------+
         | (1. Get pre-signed URL)               | (3. Serve via CDN)
         v                                       v
+-------------------+                    +---------------+
|     Droplet       |                    |  Spaces CDN   |
| [ Go Backend ]    |                    |  Edge Nodes   |
+-------------------+                    +---------------+
         | (2. Upload directly)                  ^
         |                                       |
         +---------------------------------------+
                             | (Sync)
                     +---------------+
                     |  DO Spaces    |
                     |  (Origin)     |
                     +---------------+
```

## Implementation Steps

### 1. Configure Spaces Bucket
- Go to Spaces in DO dashboard.
- Ensure the bucket is set to **Restrict File Listing** (private by default).
- Go to Settings -> **CORS Configurations**. Add your domain (e.g., `https://cloudpulse.com`) to allow `GET`, `PUT`, `POST` with `*` allowed headers. This is required for frontend direct uploads.

### 2. Generate Spaces Access Keys
- Go to API -> Spaces Keys -> Generate New Key.
- Note the **Access Key** and **Secret Key**. (You will never see the secret key again, save it securely).

### 3. Enable CDN on the Spaces Bucket
- In Bucket Settings -> CDN -> Enable CDN.
- **Custom Subdomain**: Add a subdomain like `cdn.cloudpulse.com`. You will need to add a CNAME record in your DO Networking DNS settings pointing to the Spaces edge endpoint (e.g., `cloudpulse.nyc3.cdn.digitaloceanspaces.com`).
- Add an SSL certificate (DO manages this for free using Let's Encrypt).

### 4. Implement File Upload in Go Backend
Use the official AWS SDK for Go (`github.com/aws/aws-sdk-go-v2`), as Spaces is S3-compatible.

```go
// Example configuration (Best Practice)
cfg, err := config.LoadDefaultConfig(context.TODO(),
    config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(spacesKey, spacesSecret, "")),
    config.WithRegion("nyc3"),
    // CRITICAL: Point the endpoint resolver to DigitalOcean
    config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
        func(service, region string, options ...interface{}) (aws.Endpoint, error) {
            return aws.Endpoint{URL: "https://nyc3.digitaloceanspaces.com"}, nil
        })),
)
```

**Pre-signed URLs**: Never proxy uploads through your backend (Client -> Backend -> Spaces). This ties up your Droplet's memory, CPU, and bandwidth. Instead:
1. Client requests upload URL.
2. Backend generates a pre-signed S3 URL (valid for 5 mins) and returns it.
3. Client uploads directly to Spaces using the URL.
4. Client tells Backend the upload finished.

### 5. Update Frontend for File Uploads
Create a React component that fetches the pre-signed URL from your Go backend, then uses `fetch` or `axios` with a `PUT` request to upload the file directly to Spaces.
Ensure you display files using the CDN URL (`https://cdn.cloudpulse.com/filename.jpg`), NOT the origin URL.

### 6. Automated Database Backups to Spaces
Write a simple bash script to dump the DB and push to Spaces.
```bash
#!/bin/bash
DATE=$(date +%Y-%m-%d)
# Dump from managed DB
pg_dump $DATABASE_URL > backup-$DATE.sql
# Use s3cmd or aws-cli to push to Spaces
s3cmd put backup-$DATE.sql s3://cloudpulse-bucket/backups/
# Clean up local
rm backup-$DATE.sql
```
Add this to `crontab` on your Droplet to run daily.

> [!TIP]
> **Spaces vs Block Storage (Volumes)**
> - **Spaces (Object Storage)**: S3 compatible, HTTP access, infinite scaling, slow per-operation, cheap. Use for: user uploads, static assets, backups.
> - **Volumes (Block Storage)**: Appears as a physical hard drive to the OS, fast IOPS, mounted to one Droplet at a time. Use for: databases (if self-hosted), large cache files.

## Exercises
- [ ] **CORS Testing**: Try to upload a file from a different domain (like `localhost` without adding it to CORS) and watch the browser block it.
- [ ] **Latency Comparison**: Upload an image. Load it in your browser using the origin URL (e.g., `https://cloudpulse.nyc3.digitaloceanspaces.com/img.jpg`) vs the CDN URL. Check the Chrome DevTools Network tab for TTFB (Time to First Byte).
- [ ] **Restore Backup**: Run your backup script, download the file from Spaces, and verify you can restore it locally.

## Experiments
- **Large File Upload**: Use the pre-signed URL to upload a 500MB file from the frontend. Notice how your Go backend uses 0% CPU during the transfer.
- **Cache Invalidation**: Update a static asset (like a logo) in Spaces with the exact same filename. Notice the CDN still serves the old one. Use the DO dashboard to "Purge Cache" for that file and observe the change.
- **Wrong Credentials**: Deliberately break your Spaces Secret Key in the Go environment variables. Observe the exact AWS SDK error returned.

## Checklist
- [ ] Spaces bucket configured with CORS
- [ ] CDN endpoint created and mapped to custom domain
- [ ] Go backend configured with AWS SDK for DO Spaces
- [ ] Pre-signed URL generation implemented
- [ ] Frontend direct-to-Spaces upload implemented
- [ ] DB Backup cron job configured
- [ ] Verified assets load via CDN

## What I Learned
- Pre-signed URLs are the industry standard for handling uploads at scale, keeping heavy I/O off application servers.
- DigitalOcean Spaces is highly S3-compatible, meaning the vast ecosystem of AWS tools (SDKs, CLI, s3cmd) works seamlessly.
- Edge caching via CDN drastically reduces latency for static assets compared to serving them from a central server.

---

# DAY 6 - CI/CD Pipeline: GitHub Actions + Container Registry

## Objective
Build a complete CI/CD pipeline with GitHub Actions that tests, builds, pushes images to DigitalOcean Container Registry (DOCR), and deploys to the Droplet automatically.

## WHY CI/CD?
Currently, deploying requires SSHing into the server, pulling code, running docker build, and restarting containers. This is manual toil.
1. **Consistency**: The build environment is identical every time. No "it works on my machine" issues.
2. **Speed & Reliability**: Automation removes human error during deployment.
3. **Auditability**: You can look at GitHub Actions and see exactly what code is deployed and who deployed it.
4. **Rollback**: If a deployment breaks production, you can instantly deploy the previous Docker image.

## Resources to Create
None. We use existing GitHub (free tier) and the DOCR created in Day 2.

## Estimated Cost Change
- $0 (Using existing resources)
- **New Running Total**: ~

## Architecture Diagram

```ascii
[ Developer ]
      | (git push)
      v
[ GitHub Repository ]
      | (triggers)
      v
+-------------------------------------------------------+
|                 GitHub Actions                        |
|                                                       |
|  [ CI Job ]                [ CD Job ]                 |
|  - Linting                 - Build Docker Images      |
|  - Unit Tests              - Tag with Git SHA         |
|  - Security Scan           - Push to DOCR             |
+-------------------------------------------------------+
                                  |
                                  v
                       +-------------------+
                       |    DO Container   |
                       |    Registry       |
                       +-------------------+
                                  | (SSH execution)
+---------------------------------|---------------------+
|                      DigitalOcean Droplet             |
|                                 v                     |
|                      [ deploy.sh script ]             |
|                      - docker pull (new images)       |
|                      - docker-compose up -d           |
|                      - health check & rollback        |
+-------------------------------------------------------+
```

## Implementation Steps

### 1. Understand DO Container Registry (DOCR)
DOCR is private. To allow GitHub Actions to push, and your Droplet to pull, both need authentication.
- **Droplet**: Authenticate via `doctl registry login` (done on Day 2).
- **GitHub**: Needs a Personal Access Token (PAT) from DO.

**Tagging Strategy**: Never deploy just `latest`. Tag images with the Git commit SHA (e.g., `registry.digitalocean.com/cloudpulse/backend:a1b2c3d`). This allows exact version tracking and instant rollbacks.

### 2. Set up GitHub Repository Secrets
Go to GitHub Repo -> Settings -> Secrets and Variables -> Actions. Add:
- `DIGITALOCEAN_ACCESS_TOKEN`: DO API token.
- `REGISTRY_NAME`: Your DOCR name.
- `DROPLET_SSH_KEY`: The *private* key that corresponds to the public key on your Droplet.
- `DROPLET_IP`: Your Droplet's public IP address.
- `DB_PASSWORD`: Managed PG password (for tests if needed).

### 3. Create CI Workflow (`.github/workflows/ci.yml`)
This runs on *every* push and PR.

```yaml
name: CI

on:
  push:
    branches:
      - main
  pull_request:
    branches:
      - main

jobs:
  test-backend:
    runs-on: ubuntu-latest

    defaults:
      run:
        working-directory: backend

    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version-file: backend/go.mod

      - name: Download Go modules
        run: go mod download

      - name: Vet
        run: go vet ./...

      - name: Test
        run: go test -v ./...

  test-frontend:
    runs-on: ubuntu-latest

    defaults:
      run:
        working-directory: frontend

    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-node@v4
        with:
          node-version-file: frontend/.nvmrc
          cache: npm
          cache-dependency-path: frontend/package-lock.json

      - name: Install dependencies
        run: npm ci

      - name: Build frontend
        run: npm run build
```

### 4. Create CD Workflow (`.github/workflows/deploy.yml`)
This runs *only* on pushes to `main`.

```yaml
name: Deploy Production

on:
  workflow_run:
    workflows:
      - CI
    types:
      - completed

jobs:
  deploy:
    if: >
      github.event.workflow_run.conclusion == 'success' &&
      github.event.workflow_run.head_branch == 'main'

    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Install doctl
        uses: digitalocean/action-doctl@v2
        with:
          token: ${{ secrets.DIGITALOCEAN_ACCESS_TOKEN }}

      - name: Login to Container Registry
        run: doctl registry login --expiry-seconds 600

      ####################################################
      # Backend
      ####################################################

      - name: Build Backend
        run: |
          docker build \
            --target api \
            -t ${{ secrets.REGISTRY_NAME }}/backend:${{ github.sha }} \
            ./backend

      - name: Push Backend
        run: |
          docker push \
            ${{ secrets.REGISTRY_NAME }}/backend:${{ github.sha }}

      ####################################################
      # Worker
      ####################################################

      - name: Build Worker
        run: |
          docker build \
            --target worker \
            -t ${{ secrets.REGISTRY_NAME }}/worker:${{ github.sha }} \
            ./backend

      - name: Push Worker
        run: |
          docker push \
            ${{ secrets.REGISTRY_NAME }}/worker:${{ github.sha }}

      ####################################################
      # Frontend
      ####################################################

      - name: Build Frontend
        run: |
          docker build \
            -t ${{ secrets.REGISTRY_NAME }}/frontend:${{ github.sha }} \
            ./frontend

      - name: Push Frontend
        run: |
          docker push \
            ${{ secrets.REGISTRY_NAME }}/frontend:${{ github.sha }}
            
      ####################################################
      # Deploy
      ####################################################

      - name: Deploy
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
```

### 5. Implement Deployment Script (`deploy.sh` on Droplet)
Create this on your Droplet at `/opt/cloudpulse/deploy.sh`:

```bash
#!/usr/bin/env bash

set -Eeuo pipefail

COMPOSE_FILE="docker-compose.prod.yml"

echo "===================================="
echo "Deploying ${APP_VERSION}"
echo "===================================="

# Login to registry
doctl registry login --expiry-seconds 600

# Pull latest tagged images
docker compose -f "$COMPOSE_FILE" pull

# Start new containers
docker compose -f "$COMPOSE_FILE" up -d --remove-orphans

echo "Waiting for backend..."

for i in {1..30}; do
    if curl -fsS http://localhost:8000/healthz >/dev/null; then
        echo "Deployment successful."

        # Optional cleanup
        docker image prune -af

        exit 0
    fi

    sleep 2
done

echo "Deployment failed."

docker compose -f "$COMPOSE_FILE" logs backend

exit 1
```
Update your `docker-compose.prod.yml` to use the tag:
`image: registry.digitalocean.com/cloudpulse/backend:${APP_VERSION:-latest}`

> [!WARNING]
> **Common Mistake**: Hardcoding secrets in GitHub Actions workflows. ALWAYS use GitHub Secrets. Furthermore, restrict the SSH key. Consider creating a specific `deploy` user on the Droplet with restricted `sudo` access rather than deploying as `root`.

## Security Best Practices
- **Garbage Collection**: DOCR Basic tier is 500MB. If you push a new image on every commit, you will run out of space quickly. Set up a Garbage Collection policy in the Cloud Console -> Container Registry -> Settings -> Garbage Collection. Keep only the last 5 tags.

## Exercises
- [ ] **Push and Watch**: Make a small visible change in the frontend, `git push origin main`, and watch the Actions tab in GitHub. Verify the change is live on your Droplet.
- [ ] **Catch Failure**: Intentionally break a unit test and push to a PR. Verify the Action fails and prevents merging (if branch protection is on).
- [ ] **Badge**: Add a GitHub Actions status badge to your `README.md`.

## Experiments
- **Measure Deployment Time**: How long does the `deploy.yml` workflow take? Look at the logs. Where is the bottleneck? (Usually building the Docker images).
- **Broken Image Rollback**: Push code that passes tests but fails at runtime (e.g., misconfigured ENV var causing immediate crash). Watch the `deploy.sh` health check fail and simulate a manual rollback by running `export APP_VERSION=<previous-sha> && docker-compose up -d`.
- **Concurrent Deployments**: Push two commits to `main` within 10 seconds of each other. How does GitHub Actions handle it? (Hint: Check out Actions `concurrency` groups to prevent overlapping deployments).

## Checklist
- [ ] Configured GitHub Repository Secrets
- [ ] Created CI workflow (lint & test)
- [ ] Created CD workflow (build, push, ssh)
- [ ] Implemented deploy.sh on Droplet with health check
- [ ] Updated docker-compose.prod.yml to use version tags
- [ ] Configured DOCR Garbage Collection
- [ ] Successfully triggered automated deployment

## What I Learned
- CI/CD shifts the burden of deployment from humans to machines, significantly reducing risk.
- Tagging Docker images with Git SHAs creates an unbreakable link between your repository state and your production environment.
- Automated health checks and rollbacks are critical for zero-downtime (or minimal downtime) deployments on single-node setups.
# DAY 7 - Horizontal Scaling: Load Balancer + Multi-Droplet

## Objective
Scale CloudPulse horizontally by adding a second Droplet behind a DigitalOcean Load Balancer. You will learn about stateless application design, session management, and traffic distribution.

## WHY Horizontal Scaling?
Scaling vertically (adding more CPU/RAM to a single Droplet) has a hard ceiling and requires downtime. Horizontal scaling (adding more Droplets) provides:
- **Availability & Fault Tolerance:** If one Droplet dies, the other handles traffic. No single point of failure at the compute layer.
- **Capacity:** Distributes CPU/memory load across multiple servers.
- **Zero-Downtime Deployments:** You can update one server while the other handles user traffic.

## Prerequisites: The Stateless Application
Before adding a second server, the application MUST be stateless. What does this mean? Any request from any user must be able to hit *any* server and succeed.
- **Session State:** You can't store user sessions in memory. If a user logs in on Droplet 1, and their next request hits Droplet 2, they will be logged out unless the session is stored in a shared location. For CloudPulse, sessions must be managed in our Managed Redis.
- **File Storage:** You can't save user uploads to the local disk. They must go to our Spaces object storage bucket.
- **Background Jobs:** Tasks cannot just run randomly in memory. They must use a shared queue (Redis) to ensure jobs aren't duplicated or dropped.

## Resources to Create & Update
- **Create:** Second Droplet (`s-2vcpu-4gb`, same VPC, same Firewall)
- **Create:** Load Balancer (Regional, )
- **Update:** Firewall (add new Droplet to the target tags)
- **Update:** DNS (point the domain A record to the Load Balancer IP instead of Droplet 1)

## Estimated Cost Change
- Droplet 1: 
- **Droplet 2: +**
- **Load Balancer: +**
- Managed PG: 
- Managed Redis: 
- Spaces: 
- Registry: 
- **Total:** ~ (+)

## Architecture Diagram

```ascii
                                +-------------------+
                                |                   |
                                |   Internet / DNS  |
                                |                   |
                                +---------+---------+
                                          |
                                          v (HTTPS :443)
+-------------------------------------------------------------------------+
| VPC                                                                     |
|                               +-------------------+                     |
|                               | DigitalOcean      | SSL Termination     |
|                               | Load Balancer     |                     |
|                               +---------+---------+                     |
|                                         | (HTTP :80)                    |
|                         +---------------+---------------+               |
|                         | Round Robin / Least Conn      |               |
|                         v                               v               |
|               +-------------------+           +-------------------+     |
|               |                   |           |                   |     |
|               |     Droplet 1     |           |     Droplet 2     |     |
|               |  (CloudPulse App) |           |  (CloudPulse App) |     |
|               |                   |           |                   |     |
|               +---------+---------+           +---------+---------+     |
|                         |                               |               |
|                         +---------------+---------------+               |
|                                         |                               |
|          +----------------------+-------+-------+--------------------+  |
|          |                      |               |                    |  |
|          v                      v               v                    v  |
| +-----------------+   +-----------------+   +--------+      +-----------+
| | Managed         |   | Managed         |   | Spaces |      | Container |
| | PostgreSQL      |   | Redis           |   | (CDN)  |      | Registry  |
| +-----------------+   +-----------------+   +--------+      +-----------+
+-------------------------------------------------------------------------+
```

## Implementation Steps

### 1. Verify Application is Stateless
Audit your application configuration before proceeding:
- Ensure `SESSION_DRIVER=redis` (or equivalent) in your `.env`.
- Ensure `STORAGE_DRIVER=s3` pointing to DO Spaces.
- Double-check that no SQLite databases or local JSON files are being used for application logic.

### 2. Create a Droplet Snapshot
Instead of provisioning a new server from scratch and re-running all our Ansible/bash scripts, we'll use a snapshot.
- **Why:** Fast, identical provisioning. It ensures Droplet 2 has the exact same OS dependencies, Docker versions, and monitoring agents as Droplet 1.
- **Cost:** Snapshots cost $0.06/GiB/mo. A 20GB disk snapshot costs about .
- **Action:** Go to Droplet 1 -> Snapshots -> Take Snapshot. Name it `cloudpulse-base-snapshot`.

### 3. Create Second Droplet
- Go to Create -> Droplet.
- Choose "Custom Images" and select your `cloudpulse-base-snapshot`.
- **Size:** `s-2vcpu-4gb` .
- **VPC:** Select the existing CloudPulse VPC.
- **Tags:** Add the same tags (e.g., `cloudpulse-app`) so the firewall rules automatically apply.
- **Name:** `cloudpulse-app-02`.
- SSH in and verify Docker is running and the `.env` file is present. Pull the latest image and start the container to ensure it connects to the DB and Redis successfully.

### 4. Create Load Balancer
- Go to Networking -> Load Balancers -> Create.
- **Type:** Regional.
- **Region:** Same region as your Droplets (e.g., NYC3).
- **Forwarding Rules:**
  - `HTTP` on port `80` -> `HTTP` on port `80`
  - `HTTPS` on port `443` -> `HTTP` on port `80` (This is called SSL Termination. The LB decrypts the HTTPS traffic and forwards it as HTTP within the secure VPC).
- **SSL Certificate:** Choose "Let's Encrypt" and map it to your domain. DigitalOcean will manage the auto-renewal.
- **Health Checks:**
  - Protocol: HTTP
  - Path: `/healthz` (ensure your app has a lightweight health endpoint that doesn't do heavy DB queries).
  - Port: 80
  - Interval: 10 seconds.
  - Healthy Threshold: 3 (needs 3 consecutive successes to send traffic).
  - Unhealthy Threshold: 5 (needs 5 consecutive failures to stop sending traffic).
- **Algorithm:** Round Robin (distributes evenly) vs Least Connections (sends to the server with fewest active connections). Use Round Robin for now.
- **Sticky Sessions:** Leave disabled. We are stateless. (Sticky sessions use cookies to ensure a user always hits the same Droplet, which is an anti-pattern for modern stateless apps).
- **Targets:** Add both `cloudpulse-app-01` and `cloudpulse-app-02`.

### 5. Update DNS
- Go to Networking -> Domains.
- Find your `A` record for `cloudpulse.yourdomain.com`.
- Change the target IP from Droplet 1's public IP to the **Load Balancer's public IP**.
- *Note on TTL (Time to Live):* DNS changes can take time to propagate. Lower the TTL to 60 seconds before making this change if you anticipate issues.

### 6. Configure SSL at Load Balancer Level
- SSH into your Droplets and disable/remove the Nginx SSL configuration or Certbot containers if you had them.
- Your app container/Nginx should now only listen on port 80.
- All backend communication between the LB and Droplets happens securely over the private VPC.

### 7. Update CI/CD Pipeline
Your GitHub Actions deploy script currently SSHes into one Droplet. Update it to deploy to both:
```yaml
deploy:
  runs-on: ubuntu-latest
  steps:
    - name: Deploy to Droplet 1
      uses: appleboy/ssh-action@master
      with:
        host: ${{ secrets.DROPLET_1_IP }}
        script: docker compose pull && docker compose up -d

    - name: Deploy to Droplet 2
      uses: appleboy/ssh-action@master
      with:
        host: ${{ secrets.DROPLET_2_IP }}
        script: docker compose pull && docker compose up -d
```
*Pro-tip:* For zero-downtime, take Droplet 1 out of the LB, deploy, wait for health check, put it back, then repeat for Droplet 2. We'll implement this later.

### 8. Update Monitoring
- Ensure Prometheus on your monitoring node is scraping both `cloudpulse-app-01` and `cloudpulse-app-02` by adding the new private IP to the `prometheus.yml` scrape targets.
- DigitalOcean automatically provides Load Balancer metrics (bandwidth, latency, error rates) in the DO control panel.

## Deep Dive: Health Checks and Session Affinity
**Health Checks** are the heartbeat of horizontal scaling. If a server gets stuck in a CPU loop or runs out of memory, the `/healthz` endpoint will time out. The LB detects this and stops routing traffic to it, saving your users from errors. 
**Session Affinity (Sticky Sessions)** is a crutch for legacy apps that store state in memory. The LB injects a cookie (`DO-LB`) to route the user back to the exact same server. Avoid this. If that server dies, the user loses their session. Using Redis for sessions is the modern, robust approach.

## Exercises
- [ ] Verify traffic is distributed: Tail the access logs on both Droplets simultaneously and refresh the website. You should see logs appearing on both.
- [ ] Test zero-downtime deploy: Run your CI/CD pipeline while running a load test. You shouldn't see any 502/503 errors.
- [ ] Monitor Load Balancer metrics in the DO panel (look at connection counts).
- [ ] Add a third Droplet temporarily just by cloning the snapshot and adding it to the LB target group.

## Experiments
- **Kill Droplet 1:** Stop the Docker container on Droplet 1. Watch the LB health checks mark it as down. Verify the website remains accessible (served entirely by Droplet 2).
- **Simulate Latency:** Add an artificial `sleep` to your app's health check endpoint. Watch the LB mark it unhealthy due to timeouts.
- **Break the Health Endpoint:** Return a `500 Internal Server Error` from `/healthz`. Observe the LB remove the node from the pool.

## Checklist
- [ ] App confirmed stateless (Redis for sessions, Spaces for files).
- [ ] Snapshot created from Droplet 1.
- [ ] Droplet 2 provisioned from snapshot.
- [ ] Load Balancer created with SSL termination.
- [ ] DNS updated to point to Load Balancer.
- [ ] CI/CD pipeline updated to deploy to both nodes.
- [ ] Failover tested successfully.

## What I Learned
- The fundamental difference between vertical and horizontal scaling.
- Why statelessness is a hard requirement for modern web applications.
- How load balancers use health checks to route around failure.
- How SSL termination offloads CPU work from the application servers.

---

## Bonus: Block Storage Volumes (Day 7 Lab)

DigitalOcean Block Storage Volumes are network-attached SSD drives that you can mount to a Droplet like a physical hard drive. While CloudPulse uses Spaces (object storage) for file uploads, Block Storage is the right choice for database data directories, large cache files, or any workload that needs fast random IOPS.

> **Spaces vs Block Storage:**
> - **Spaces** = S3-like object storage. Great for: file uploads, backups, static assets, CDN.
> - **Block Storage** = SSD disk attached to one Droplet. Great for: databases, application logs, high-IOPS workloads.

### Create, Format, Mount, Resize, Snapshot, Restore

**1. Create a Block Storage Volume**
- Name: `cloudpulse-data`
- Size: 10 GiB (minimum, can resize later)
- Region: Same as your Droplet (e.g., `nyc3`)
- Filesystem: ext4
- Attach to: Droplet 1
- Cost:  (10 GiB × $0.10/GiB)

**2. Format and Mount the Volume**
```bash
# SSH into the Droplet
ssh deploy@<droplet_ip>

# List block devices — you should see /dev/sda (the volume)
lsblk

# Format the volume with ext4
sudo mkfs.ext4 /dev/sda

# Create a mount point
sudo mkdir -p /mnt/cloudpulse-data

# Mount the volume
sudo mount -o defaults,noatime /dev/sda /mnt/cloudpulse-data

# Verify it's mounted
df -h /mnt/cloudpulse-data

# Make it persistent across reboots (add to /etc/fstab)
echo '/dev/sda /mnt/cloudpulse-data ext4 defaults,noatime 0 2' | sudo tee -a /etc/fstab
```
*Why `noatime`?* It skips updating the "last accessed" timestamp on every read, reducing unnecessary disk writes.

**3. Resize the Volume**
```bash
# In DO Console: Resize the volume from 10 GiB to 20 GiB
# Then on the Droplet, extend the filesystem:
sudo resize2fs /dev/sda
df -h /mnt/cloudpulse-data  # Should now show 20 GiB
```
*Note:* DO Block Storage volumes can be resized up (never down) without unmounting — the filesystem can be extended live.

**4. Create a Volume Snapshot**
- In DO Console: Volumes → cloudpulse-data → Create Snapshot
- Name: `cloudpulse-data-snap-day7`
- Cost: $0.06/GiB/month

*Why snapshot?* Snapshots are point-in-time copies. If you corrupt data on the volume, you can restore from the snapshot.

**5. Restore from Snapshot**
- In DO Console: Snapshots → Volumes → Select snapshot → Create Volume
- This creates a brand-new volume from the snapshot
- Detach the old volume, attach the new one, mount it at the same path

### Block Storage Exercises
- [ ] Write a file to the volume: `echo "hello" > /mnt/cloudpulse-data/test.txt`
- [ ] Reboot the Droplet and verify the volume remounts automatically
- [ ] Resize the volume from 10 to 20 GiB without unmounting
- [ ] Create a snapshot, delete the test file, restore from snapshot, verify the file is back
- [ ] Detach the volume and attach it to a different Droplet

### Block Storage Experiments
- **Performance test:** Run `dd if=/dev/zero of=/mnt/cloudpulse-data/benchfile bs=1M count=500` to measure write speed
- **Detach while mounted:** Try detaching a volume that's mounted — observe the error
- **Destroy and restore:** Delete the volume entirely, then recreate from snapshot

> **When to use Block Storage in production:** Mount a volume for PostgreSQL data (`/var/lib/postgresql/data`) if you're self-hosting the database. This separates your data from the Droplet's boot disk — if the Droplet fails, the data volume survives.

**Cleanup:** Destroy the Block Storage Volume after completing the exercises (it's a temporary learning resource). Keep the snapshot for the DR lab on Day 11.

---

# DAY 8 - Kubernetes: DigitalOcean Kubernetes (DOKS)

## Objective
Deploy CloudPulse to a DigitalOcean Kubernetes cluster (DOKS). Learn Kubernetes fundamentals through hands-on deployment of a real application.

> [!NOTE]
> **Credits are expiring — keep K8s alive through Day 10.** We'll run the DOKS cluster alongside our Droplet setup for 3 days (Days 8-10). This lets you compare both approaches side-by-side and gives you ample time to experiment. We'll destroy it on Day 10 when we deploy to App Platform.

## WHY Kubernetes?
If our Load Balancer + Droplet setup works, why use Kubernetes (K8s)?
- **Orchestration & Self-Healing:** If a Docker container crashes on a Droplet, it stays dead unless a process manager restarts it. K8s actively monitors desired state. If a pod dies, K8s creates a new one instantly.
- **Auto-Scaling:** K8s can automatically add more Pods (HPA) or more physical Nodes (Cluster Autoscaler) based on CPU/Memory usage.
- **Declarative Configuration:** You define *what* you want (e.g., "I want 3 instances of the backend"), and K8s figures out *how* to make it happen.
- **Service Discovery & Routing:** Built-in internal DNS.

**Honest Comparison:** 
- *App Platform (PaaS):* Best for small teams prioritizing speed over control. High cost at scale.
- *Droplets (IaaS):* Best for predictable workloads and tight budgets. High operational overhead.
- *Kubernetes (CaaS):* Best for complex microservices, massive scale, and high availability. Steep learning curve, higher base cost.

## Resources to Create
- **Create:** DOKS Cluster (3 worker nodes, `s-2vcpu-4gb` each, Standard control plane).
- **Keep Running:** Both Droplets, Load Balancer, Monitoring Droplet — run everything in parallel for comparison.

> **Why keep Droplets alive alongside K8s?** With credits expiring, we can afford to run both environments simultaneously. This is incredibly valuable — you can deploy the same app to Droplets AND Kubernetes and compare performance, latency, deployment speed, and operational complexity side-by-side. In the real world you rarely get this luxury.

## Estimated Cost (All Resources Running)
- App Droplet 1: 
- App Droplet 2: 
- Monitoring Droplet: 
- Load Balancer: 
- DOKS Nodes (3 × $24): 
- Managed PG + Redis: 
- Spaces + Registry: 
- **Total during Days 8-9:** ~ (~)

## Architecture Diagram

```ascii
                                +-------------------+
                                |   Internet / DNS  |
                                +---------+---------+
                                          |
+-----------------------------------------|---------------------------------+
| DigitalOcean Kubernetes (DOKS)          v                                 |
|                               +-------------------+                       |
|                               | DO Load Balancer  | (Created by Ingress)  |
|                               +---------+---------+                       |
|                                         |                                 |
|                               +---------v---------+                       |
|                               | Nginx Ingress     |                       |
|                               | Controller Pod    |                       |
|                               +----+---------+----+                       |
|                                    |         |                            |
|             /api/* routing         |         |  / routing                 |
|            +-----------------------+         +-----------------------+    |
|            v                                                         v    |
|  +-------------------+                                     +-------------------+
|  | Backend Service   |                                     | Frontend Service  |
|  +---------+---------+                                     +---------+---------+
|            |                                                         |    |
|   +--------+--------+                                       +--------+    |
|   v                 v                                       v             v
| +-------+         +-------+                               +-------+     +-------+
| | Pod 1 |         | Pod 2 |                               | Pod 1 |     | Pod 2 |
| | (Go)  |         | (Go)  |                               | (Web) |     | (Web) |
| +-------+         +-------+                               +-------+     +-------+
+-------------------------------------------------------------------------------+
           |                 |                                      |
           v                 v                                      v
    +--------------+  +--------------+                       +--------------+
    | Managed PG   |  | Managed Redis|                       | Spaces (CDN) |
    +--------------+  +--------------+                       +--------------+
```

## Kubernetes Concepts Primer
- **Cluster & Node:** A Cluster is the entire system. Nodes are the physical VMs (Droplets) running the workloads.
- **Pod:** The smallest deployable unit in K8s. Usually contains one Docker container.
- **Deployment:** Manages Pods. Defines how many replicas you want and how to update them.
- **Service:** A stable internal IP address that routes traffic to a set of ephemeral Pods.
- **Ingress:** Manages external access to Services (usually HTTP/HTTPS routing rules).
- **ConfigMap / Secret:** Injects environment variables and configs into Pods.

## Implementation Steps

### 1. Create DOKS Cluster
- Go to Kubernetes -> Create Cluster.
- **Region:** Same as your DB/VPC.
- **Version:** Latest stable.
- **VPC:** Select your existing VPC.
- **Node Pool:** 3 nodes, size `s-2vcpu-4gb` .
- **Name:** `cloudpulse-k8s`.

### 2. Install Tools & Authenticate
Install `kubectl` and `doctl` on your local machine.
```bash
# Authenticate doctl
doctl auth init

# Pull kubeconfig to connect kubectl to your cluster
doctl kubernetes cluster kubeconfig save cloudpulse-k8s

# Verify
kubectl get nodes
```

### 3. Connect Registry to DOKS
K8s needs permission to pull your private Docker images.
```bash
doctl kubernetes cluster registry add cloudpulse-k8s
```

### 4. Create Kubernetes Manifests
Create a directory `k8s/` and add these files. *Note: As a senior engineer, I enforce setting resource requests/limits. Without them, one bad pod can consume an entire node's CPU and crash the node.*

`k8s/namespace.yaml`
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: cloudpulse
```

`k8s/secret.yaml` (Base64 encode your secrets before applying, or use external secret operators in production)
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: cloudpulse-secrets
  namespace: cloudpulse
type: Opaque
data:
  DATABASE_URL: <base64-encoded-url>
  REDIS_URL: <base64-encoded-url>
```

`k8s/backend-deployment.yaml`
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend
  namespace: cloudpulse
spec:
  replicas: 3
  selector:
    matchLabels:
      app: backend
  template:
    metadata:
      labels:
        app: backend
    spec:
      containers:
      - name: backend
        image: registry.digitalocean.com/your-registry/backend:latest
        ports:
        - containerPort: 8080
        envFrom:
        - secretRef:
            name: cloudpulse-secrets
        resources:
          requests:
            cpu: "100m"
            memory: "128Mi"
          limits:
            cpu: "500m"
            memory: "256Mi"
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
```

`k8s/backend-service.yaml`
```yaml
apiVersion: v1
kind: Service
metadata:
  name: backend-svc
  namespace: cloudpulse
spec:
  selector:
    app: backend
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP
```

*(Create similar Deployments and Services for Frontend and Worker)*

### 5. Deploy Ingress Controller
We need an Ingress Controller to route outside traffic into the cluster. DO makes it easy to deploy Nginx Ingress, which automatically provisions a DO Load Balancer.
```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.2/deploy/static/provider/do/deploy.yaml
```

`k8s/ingress.yaml`
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: cloudpulse-ingress
  namespace: cloudpulse
  annotations:
    kubernetes.io/ingress.class: nginx
spec:
  rules:
  - host: cloudpulse.yourdomain.com
    http:
      paths:
      - path: /api
        pathType: Prefix
        backend:
          service:
            name: backend-svc
            port: 
              number: 80
      - path: /
        pathType: Prefix
        backend:
          service:
            name: frontend-svc
            port: 
              number: 80
```

### 6. Apply the Manifests
```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/secret.yaml
kubectl apply -f k8s/backend-deployment.yaml
kubectl apply -f k8s/backend-service.yaml
kubectl apply -f k8s/ingress.yaml
```

### 7. Configure DNS
Run `kubectl get svc -n ingress-nginx` to find the External-IP of the Load Balancer K8s created. Update your domain's A record to this IP.

### Bonus: PersistentVolume and PersistentVolumeClaim

While CloudPulse is stateless (databases are managed by DO), you should understand how Kubernetes handles persistent storage — it's critical for self-hosted databases, log aggregation, or any stateful workload.

```yaml
# k8s/pvc.yaml — PersistentVolumeClaim using DO Block Storage
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: cloudpulse-data-pvc
  namespace: cloudpulse
spec:
  accessModes:
    - ReadWriteOnce        # Block storage mounts to one node at a time
  storageClassName: do-block-storage  # DO's built-in CSI driver
  resources:
    requests:
      storage: 5Gi         # Provisions a 5 GiB DO Block Storage Volume
```

Apply it: `kubectl apply -f k8s/pvc.yaml`

To use it in a Deployment (example: mounting to the simulator for writing test data):
```yaml
# Add to the simulator deployment spec:
      volumes:
        - name: data-vol
          persistentVolumeClaim:
            claimName: cloudpulse-data-pvc
      containers:
        - name: simulator
          volumeMounts:
            - name: data-vol
              mountPath: /data
```

> **Key concepts:**
> - **PersistentVolume (PV):** The actual disk (DO Block Storage Volume, auto-provisioned by the CSI driver).
> - **PersistentVolumeClaim (PVC):** A request for storage. K8s binds a PVC to a PV.
> - **StorageClass:** `do-block-storage` tells K8s to use DigitalOcean's CSI driver to dynamically provision volumes.
> - **AccessModes:** `ReadWriteOnce` = one node at a time (block storage). For shared storage, use `ReadWriteMany` (requires NFS or similar).

**Exercise:** Create the PVC, exec into the simulator pod, write a file to `/data`, delete the pod, wait for K8s to recreate it, and verify the file persists.

## Common Kubernetes Mistakes
- **Not setting resource requests/limits:** Pods will cannibalize each other, leading to Out Of Memory (OOM) kills and node instability.
- **Not configuring health probes:** If a Go app deadlocks but the container process doesn't exit, K8s won't know it's broken unless a `livenessProbe` fails.
- **Using `latest` image tags:** Always tag images with Git SHAs. If you use `latest`, K8s doesn't know the image changed, so it won't pull the new one on pod restarts.

## Exercises
- [ ] **Scale manually:** `kubectl scale deployment backend --replicas=5 -n cloudpulse`
- [ ] **Rollout status:** `kubectl rollout status deployment/backend -n cloudpulse`
- [ ] **Exec into a Pod:** `kubectl exec -it <pod-name> -n cloudpulse -- /bin/sh`
- [ ] **View logs:** `kubectl logs -f <pod-name> -n cloudpulse`

## Experiments
- **Delete a Pod:** Run `kubectl delete pod <pod-name>`. Watch it terminate and instantly be replaced by a new one to maintain the desired replica count.
- **Exceed Resource Limits:** Deploy a pod with a memory limit of `10Mi` and run a script that allocates memory. Watch it get `OOMKilled`.
- **Break Readiness Probe:** Change the `readinessProbe` to a bad path. Watch the pod start, but notice it never gets added to the Service endpoints (it receives no traffic).

## Checklist
- [ ] `kubectl` and `doctl` configured.
- [ ] Registry credentials added to cluster.
- [ ] Deployments, Services, and Ingress applied.
- [ ] App accessible via browser through Ingress.
- [ ] Manual scaling verified.

## What I Learned
- K8s separates the desired state (YAML) from the actual state (running infrastructure).
- How internal networking (Services) abstracts away the ephemeral nature of Pods.
- The critical importance of resource limits and health probes in a container orchestration environment.

---

# DAY 9 - Infrastructure as Code: Terraform

## Objective
Define ALL DigitalOcean infrastructure as code using Terraform. Import existing resources, then tear down and recreate everything from code.

## WHY Infrastructure as Code (IaC)?
Clicking through a web UI to create servers is fine for tutorials, but unacceptable in production.
- **Reproducibility:** You can spin up an identical staging environment in 5 minutes.
- **Version Control:** Infrastructure changes are PR'd, reviewed, and audited in Git just like application code.
- **Drift Detection:** If someone manually changes a firewall rule, Terraform detects it and reverts it to the defined state.
- **Disaster Recovery:** If your account is compromised and resources deleted, you can rebuild the core infra immediately.

## Resources to Destroy
- **Destroy the DOKS cluster entirely.** We learned K8s, but keeping it running is too expensive for this bootcamp. Use the DO Console to delete the cluster (this will also delete the associated Load Balancer).
- We are now back to our core resources: Droplet 1 (powered back on), Managed PG, Managed Redis, Spaces, Registry, VPC, Firewall, Domain.

## Architecture Diagram

```ascii
+-----------------------------------------------------------------------+
| GitHub Repository                                                     |
|                                                                       |
|  +----------------+    git push     +------------------------------+  |
|  |                | --------------> | GitHub Actions (CI/CD)       |  |
|  | main.tf        |                 |                              |  |
|  | droplet.tf     |                 |  terraform plan              |  |
|  | database.tf    |                 |  terraform apply             |  |
|  | firewall.tf    |                 +---------------+--------------+  |
|  +----------------+                                 |                 |
+-----------------------------------------------------|-----------------+
                                                      |
                                                      v (DO API)
+-----------------------------------------------------------------------+
| DigitalOcean Cloud Environment                                        |
|                                                                       |
|  +--------------+  +-------------+  +-------------+  +-------------+  |
|  |              |  |             |  |             |  |             |  |
|  |  Droplet     |  | Managed PG  |  | Managed     |  |  Spaces     |  |
|  |              |  |             |  | Redis       |  |             |  |
|  +--------------+  +-------------+  +-------------+  +-------------+  |
|                                                                       |
|  +--------------+  +-------------+  +-------------+                   |
|  |              |  |             |  |             |                   |
|  |  VPC         |  | Firewall    |  | Registry    |                   |
|  |              |  |             |  |             |                   |
|  +--------------+  +-------------+  +-------------+                   |
+-----------------------------------------------------------------------+
```

## Terraform Concepts Primer
- **Provider:** Plugins that translate Terraform code into API calls (e.g., the DigitalOcean provider).
- **Resource:** A piece of infrastructure (e.g., `digitalocean_droplet`).
- **State File (`terraform.tfstate`):** A JSON file mapping your `.tf` code to the actual resource IDs in the cloud. *Critical:* Never lose this, never commit it to Git with sensitive data in plain text.
- **Plan:** A dry-run showing what will be created/modified/destroyed.
- **Apply:** Executes the plan.

## Implementation Steps

### 1. Install Terraform
Download and install the Terraform CLI for your OS.

### 2. Set Up Project Structure
Create an `infrastructure/` directory in your repo:
```text
infrastructure/terraform/
â”œâ”€â”€ main.tf           # Provider config
â”œâ”€â”€ variables.tf      # Input variables
â”œâ”€â”€ vpc.tf            # VPC resource
â”œâ”€â”€ firewall.tf       # Firewall rules
â”œâ”€â”€ droplet.tf        # Droplet resources
â”œâ”€â”€ database.tf       # Managed databases
â”œâ”€â”€ spaces.tf         # Object storage
â””â”€â”€ registry.tf       # Container registry
```

### 3. Configure Provider (`main.tf`)
```hcl
terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.30.0"
    }
  }
  
  # We will configure remote backend later
  # backend "s3" {} 
}

variable "do_token" {}

provider "digitalocean" {
  token = var.do_token
}
```

### 4. Define Resources (Examples)

`droplet.tf`
```hcl
resource "digitalocean_droplet" "web" {
  image    = "ubuntu-22-04-x64"
  name     = "cloudpulse-app-01"
  region   = "nyc3"
  size     = "s-2vcpu-4gb"
  vpc_uuid = digitalocean_vpc.main.id
  tags     = ["cloudpulse-app"]
}
```

`database.tf`
```hcl
resource "digitalocean_database_cluster" "postgres" {
  name       = "cloudpulse-db"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-2gb"
  region     = "nyc3"
  node_count = 1
  private_network_uuid = digitalocean_vpc.main.id
}
```

`loadbalancer.tf` (ready to apply when scaling horizontally)
```hcl
resource "digitalocean_loadbalancer" "web" {
  name   = "cloudpulse-lb"
  region = "nyc3"

  forwarding_rule {
    entry_port     = 443
    entry_protocol = "https"
    target_port     = 80
    target_protocol = "http"
    certificate_name = digitalocean_certificate.cert.name
  }

  forwarding_rule {
    entry_port     = 80
    entry_protocol = "http"
    target_port     = 80
    target_protocol = "http"
  }

  healthcheck {
    port     = 80
    protocol = "http"
    path     = "/healthz"
  }

  droplet_tag = "cloudpulse-app"
  vpc_uuid    = digitalocean_vpc.main.id
}
```
*Note:* This resource can be kept commented out or in a separate workspace. Apply it when you need horizontal scaling — Terraform makes it trivial to add/remove a Load Balancer.

### 5. Import Existing Resources
Since we already created these resources in the UI, if we just run `apply`, Terraform will try to create *new* ones. We must import the existing IDs into our state.

1. Get the Droplet ID from the DO Console URL or DO API.
2. Run the import command:
   ```bash
   export TF_VAR_do_token="your_do_token"
   terraform init
   terraform import digitalocean_droplet.web 123456789
   terraform import digitalocean_database_cluster.postgres <db-id>
   ```
3. Run `terraform plan`. Terraform will compare your `.tf` files with the imported state. 
4. **Iterate:** If the plan says it wants to change something (e.g., tags), update your `.tf` file to match reality until the plan says `No changes. Infrastructure is up-to-date.`

### 6. Remote State in Spaces
Storing state locally is dangerous for teams. We will store it in our DO Spaces bucket (which is S3 compatible).
Update `main.tf`:
```hcl
terraform {
  backend "s3" {
    endpoint                    = "nyc3.digitaloceanspaces.com"
    region                      = "us-east-1" # S3 standard requirement, leave as is
    bucket                      = "cloudpulse-tf-state" # Create this bucket manually first
    key                         = "terraform.tfstate"
    skip_credentials_validation = true
    skip_metadata_api_check     = true
  }
}
```
Run `terraform init` to migrate your local state to Spaces.

### 7. Create a Reusable Module
Modules are like functions in IaC. Let's create a Droplet module.
Create `modules/droplet/main.tf`:
```hcl
variable "name" {}
variable "size" {}
variable "vpc_id" {}

resource "digitalocean_droplet" "this" {
  image    = "ubuntu-22-04-x64"
  name     = var.name
  region   = "nyc3"
  size     = var.size
  vpc_uuid = var.vpc_id
}

output "ip_address" {
  value = digitalocean_droplet.this.ipv4_address
}
```
Now in your root `droplet.tf`, you can call it:
```hcl
module "app_droplet" {
  source = "./modules/droplet"
  name   = "cloudpulse-app-01"
  size   = "s-2vcpu-4gb"
  vpc_id = digitalocean_vpc.main.id
}
```

## Common Terraform Mistakes
- **Not using remote state:** If your laptop dies, your state is gone. If a coworker runs terraform, they overwrite your work. Always use a remote backend.
- **Hardcoding values:** Don't hardcode region `nyc3` in 15 files. Use a variable `var.region`.
- **Not pinning versions:** Provider APIs change. If you don't pin `version = "~> 2.30.0"`, a future `terraform init` might pull a breaking update and destroy your infra.
- **Forgetting `lifecycle` blocks:** Use `prevent_destroy = true` on your production database resource to prevent accidental deletion.

## Exercises
- [ ] Define the VPC and Firewall in Terraform.
- [ ] Import all existing resources successfully so `terraform plan` is clean.
- [ ] Migrate state to DO Spaces.
- [ ] Refactor your Droplet code to use a Module.

## Experiments
- **Observe Drift:** Go into the DO Web Console and manually change the name of your Droplet. Run `terraform plan`. Terraform will detect the drift and propose changing it back. Run `terraform apply` to fix it.
- **Destroy-Recreate:** Change the `image` string in your Droplet resource. Notice that `terraform plan` says `-/+ destroy and then create replacement`. Some properties cannot be updated in-place.
- **Add to CI:** Create a GitHub Action that runs `terraform plan` on every Pull Request so reviewers can see infrastructure changes.

## Checklist
- [ ] DOKS cluster kept alive for Day 10 comparison.
- [ ] Terraform installed and provider configured.
- [ ] All `.tf` files written for existing infra.
- [ ] State successfully imported.
- [ ] Remote backend configured using Spaces.

## What I Learned
- The fundamental power of declarative infrastructure.
- How Terraform maps code to cloud resources using the state file.
- The workflow of importing existing, unmanaged infrastructure into code.
- How to structure Terraform projects for reusability.
# DAY 10 - Platform as a Service: App Platform + Serverless Functions

**Objective:** Deploy CloudPulse to DigitalOcean App Platform (PaaS) and create serverless functions. Compare deployment models: Droplets vs Kubernetes vs App Platform.

> **With credits expiring:** Keep App Platform running through Day 11 alongside your other environments. This gives you three simultaneous deployment models to compare: Droplet + LB, Kubernetes, and App Platform. Destroy the DOKS cluster on Day 10 (you've had 3 days with it), but keep App Platform alive until final cleanup.

## WHY PaaS / App Platform?
As a cloud engineer, you'll constantly face the "build vs. buy" tradeoff. Managing Droplets gives you ultimate control, but you have to patch the OS, configure networking, set up systemd services, and build CI/CD pipelines. Kubernetes (DOKS) abstracts the OS but requires managing cluster state, Helm charts, and complex manifests.

Platform as a Service (PaaS) like DigitalOcean App Platform abstracts away the underlying infrastructure entirely. You simply point it to your GitHub repository or Container Registry, define a configuration file (App Spec), and the platform handles:
- Zero infrastructure management (no OS updates, no SSH)
- Built-in CI/CD (push to `main` -> automatic build and deploy)
- Managed SSL and custom domains out of the box
- Auto-scaling (on Professional tiers)
- Distributed CDN routing automatically configured

## WHEN to use App Platform vs Droplets vs Kubernetes

| Feature / Platform | Droplets (IaaS) | Kubernetes (DOKS) | App Platform (PaaS) |
| :--- | :--- | :--- | :--- |
| **Control** | Full root access | High (cluster admin) | Minimal (app only) |
| **Maintenance** | High (patching, scaling) | Medium (upgrades, nodes) | Low (fully managed) |
| **Deployment** | Manual or custom CI/CD | Manifests, Helm | Auto from Git/Registry |
| **Scaling** | Manual / Load Balancers | Auto (HPA, Cluster Auto) | Auto (Slider in UI/Spec) |
| **Cost** | Lowest ( min) | High ( min) | Medium ( min) |
| **Best For** | Custom stacks, legacy apps, tight budgets | Microservices, complex orchestrations | Rapid iterations, startups, standard web apps |

*Senior Engineer Tip:* Always start with the simplest solution that meets your business requirements. For many modern web apps, PaaS is the right choice until you hit specific technical limitations or unit economics break down at extreme scale.

## Resources to Create & Destroy
- **Create**: App Platform application (Professional tier for auto-scaling tests), DO Functions namespace + functions.
- **Destroy**: DOKS cluster (you've had 3 full days with it — Days 8-10).
- **Keep Running**: App Platform (through Day 11 for comparison), Droplets, LB, Managed DBs.
- **Estimated Cost**: ~ (App Platform Pro replaces DOKS cost).

## Architecture Diagram (App Platform)

```text
                                     DIGITALOCEAN APP PLATFORM
                                     
      +----------------+        +-------------------------------------------------+
      |                |        |                                                 |
      |  GitHub Repo   |---+    |  +-------------------+   +-------------------+  |
      |                |   |    |  |                   |   |                   |  |
      +----------------+   |    |  | Static Site Comp  |   | Service Component |  |
                           |    |  | (Next.js FE)      |   | (Go Backend API)  |  |
      +----------------+   +--->|  |                   |   |                   |  |
      |                |        |  +--------+----------+   +---------+---------+  |
      | Container Reg. |---+    |           |                        |            |
      |                |   |    |           | API Calls              |            |
      +----------------+   |    |           v                        v            |
                           |    |  +-------------------+   +-------------------+  |
      +----------------+   +--->|  |                   |   |                   |  |
      |                |        |  | Worker Component  |   | Serverless Functs |  |
      | End User (Web) |---+    |  | (Go Background)   |   | (doctl deployed)  |  |
      |                |   |    |  |                   |   |                   |  |
      +----------------+   |    |  +-------------------+   +-------------------+  |
                           |    |                                                 |
                           |    +-------------------------------------------------+
                           |           |            |              |
                           +-----------+            |              |
                                                    v              v
                                           +----------------+ +----------------+
                                           | Managed PG     | | Managed Redis  |
                                           +----------------+ +----------------+
```

## Implementation Steps

### 1. Understand App Platform Concepts
App Platform applications consist of **Components**. A single app can have multiple components:
- **Service**: A web server accessible from the internet (e.g., our Go backend).
- **Worker**: A background process not directly accessible via HTTP (e.g., our task processor).
- **Job**: A short-lived task (e.g., database migrations) that runs before/after deployments.
- **Static Site**: HTML/CSS/JS served globally via CDN (e.g., Next.js export).
- **Function**: Serverless HTTP endpoint.

You define these in an `App Spec` (YAML). App Platform can build directly from source code using Cloud Native Buildpacks, or deploy pre-built images from a Container Registry. We will use our Container Registry images since we already built CI/CD in Day 6.

### 2. Create the App Spec (YAML)
Create a file `cloudpulse-app.yaml`. We will connect this App to our existing Managed PostgreSQL and Redis.

```yaml
name: cloudpulse-app
region: nyc3
services:
  - name: backend-api
    image:
      registry: my-registry
      registry_type: DOCR
      repository: cloudpulse-backend
      tag: latest
    envs:
      - key: DATABASE_URL
        scope: RUN_TIME
        value: ${db.DATABASE_URL}
      - key: REDIS_URL
        scope: RUN_TIME
        value: ${redis.REDIS_URL}
    http_port: 8080
    instance_count: 1
    instance_size_slug: basic-xxs
    routes:
      - path: /api

workers:
  - name: task-processor
    image:
      registry: my-registry
      registry_type: DOCR
      repository: cloudpulse-worker
      tag: latest
    envs:
      - key: DATABASE_URL
        scope: RUN_TIME
        value: ${db.DATABASE_URL}
      - key: REDIS_URL
        scope: RUN_TIME
        value: ${redis.REDIS_URL}
    instance_count: 1
    instance_size_slug: basic-xxs

static_sites:
  - name: frontend-ui
    source_dir: /
    github:
      branch: main
      deploy_on_push: true
      repo: mygithubuser/cloudpulse-frontend
    build_command: npm run build
    envs:
      - key: NEXT_PUBLIC_API_URL
        scope: BUILD_TIME
        value: "https://cloudpulse-app-xxxxx.ondigitalocean.app/api"
    routes:
      - path: /

databases:
  - name: db
    engine: PG
    production: true
    cluster_name: cloudpulse-postgres-nyc3
  - name: redis
    engine: REDIS
    production: true
    cluster_name: cloudpulse-redis-nyc3
```

*Explanation:* 
- We expose the backend on `/api` and the frontend on `/`. App Platform automatically configures the routing!
- We use `${db.DATABASE_URL}` cross-references. App Platform securely fetches credentials from our existing Managed DBs without us hardcoding secrets.
- `basic-xxs` is the  container size.

### 3. Deploy via doctl
Deploy the app using the DO CLI:
```bash
doctl apps create --spec cloudpulse-app.yaml
```
Once deployed, get the app ID and watch the build logs:
```bash
doctl apps list
doctl apps logs <APP_ID> --type build
```
Common Mistake: Forgetting that Next.js static exports need the API URL at *build time*. Ensure your `envs` scope is set to `BUILD_TIME` for the frontend.

### 4. Configure Custom Domain
If you wanted to take this to production, you'd add a domain:
```bash
doctl apps create-domain <APP_ID> --domain cloudpulse.mydomain.com
```
App Platform automatically provisions and renews Let's Encrypt TLS certificates.

### 5. Compare with Droplet Deployment
- **Performance**: Droplets have dedicated resources; App Platform Basic runs on shared compute. App Platform Professional provides dedicated vCPUs. Latency is similar, as both sit in the same region.
- **Complexity**: App Platform is vastly simpler. No reverse proxy (Caddy/Nginx) config required.
- **Build Time**: App Platform builds can be slower than custom GitHub Actions because they rely on generalized Buildpacks if not using pre-built images.

### 6. Create Serverless Functions
DigitalOcean Functions are serverless compute. Great for event-driven glue code.

1. Install the serverless plugin:
```bash
doctl serverless install
```
2. Connect to a namespace (creates one if it doesn't exist):
```bash
doctl serverless connect
```
3. Initialize a function project:
```bash
doctl serverless init cloudpulse-functions --language go
```
4. Create a webhook handler (`cloudpulse-functions/packages/default/webhook/main.go`):
```go
package main

import "fmt"

func Main(args map[string]interface{}) map[string]interface{} {
    name, ok := args["name"].(string)
    if !ok {
        name = "stranger"
    }
    return map[string]interface{}{
        "body": fmt.Sprintf("Webhook received for %s!", name),
    }
}
```
5. Deploy:
```bash
doctl serverless deploy cloudpulse-functions
```
6. Get the URL and invoke:
```bash
doctl serverless functions get default/webhook --url
curl -X POST <URL> -H "Content-Type: application/json" -d '{"name":"CloudPulse"}'
```

### 7. Compare Functions vs Worker Containers
- **Functions**: Pay-per-invocation. Infinite scale to zero. Harder to test locally. Cold starts (delay on first request). Maximum execution time limits (e.g., 15 mins).
- **Workers**: Always running (pay monthly). No cold starts. Better for continuous processing (like consuming a busy Redis stream).

## Exercises
- [ ] Deploy the entire app using the App Spec.
- [ ] Scale the worker component to 2 instances via CLI (`doctl apps update --spec ...`).
- [ ] Trigger an auto-deploy by pushing a dummy commit to your frontend Git repo.
- [ ] View build and runtime logs via `doctl` and the UI.
- [ ] Test the serverless webhook function with `curl`.

## Experiments
- **Push a broken build:** Push code with a syntax error. Observe how App Platform cancels the deployment and keeps the previous healthy version running (Blue/Green deployment built-in).
- **Cold start times:** Invoke the serverless function after 30 minutes of inactivity. Measure the latency difference between the first invoke and the second.


## Cleanup Instructions
Destroy the DOKS cluster (you've had 3 days with it). Keep App Platform running through Day 11 for the final comparison.
```bash
# Get your App ID
doctl apps list
# Destroy the app
doctl apps delete <APP_ID> -f

# Clean up functions
doctl serverless undeploy cloudpulse-functions
```

## What I Learned
- App Platform heavily reduces operational burden at the cost of granular control.
- App Specs allow for Infrastructure as Code (IaC) without needing Terraform.
- Serverless functions are excellent for bursty, event-driven workloads (cron jobs, webhooks, file uploads) but suffer from cold starts compared to always-on Worker containers.

---

# DAY 11 - Production Architecture: Security, Disaster Recovery, and Final Review

**Objective:** Harden the CloudPulse deployment for production, implement disaster recovery, perform a final architecture review, and consolidate all learnings.

> **THIS IS THE FINAL DAY!** We will take everything we built, secure it, prove we can recover from disaster, and review the final architecture.

## Part 1: Security Hardening

Security is not a checkbox; it's a posture. As a Senior Engineer, you must assume your infrastructure *will* be attacked. 

### Implementation Steps

1. **SSH Hardening (Droplets)**
   - Never use passwords. Enforce SSH keys.
   - Disable root login in `/etc/ssh/sshd_config`: `PermitRootLogin no`.
   - Disable password authentication: `PasswordAuthentication no`.
   - Install `fail2ban` to protect against brute-force attacks:
     ```bash
     sudo apt update && sudo apt install fail2ban -y
     sudo systemctl enable fail2ban --now
     ```

2. **Automatic Security Updates**
   - Configure `unattended-upgrades` so the OS automatically installs critical security patches:
     ```bash
     sudo apt install -y unattended-upgrades
     sudo dpkg-reconfigure -plow unattended-upgrades  # Select "Yes"
     ```
   - Verify it's active: `systemctl status unattended-upgrades`
   - *Why?* Unpatched kernel or OpenSSL vulnerabilities are the #1 cause of cloud breaches. Automatic updates for security patches are a must.
   - *Caveat:* Only enable for security updates, not feature updates, to avoid breaking changes.

3. **Firewall & Network Security**
   - **Principle of Least Privilege**: Only open what is strictly necessary.
   - Ensure CloudPulse Droplets only accept HTTP/HTTPS from the open internet, and SSH from your specific IP (or a jump host).
   - Databases should **NEVER** have public IPs. Configure your Managed PostgreSQL and Redis to only accept connections from your VPC (trusted sources).

4. **Application Security**
   - **HTTPS Everywhere**: Handled by Nginx + Let's Encrypt (Day 2) or the Load Balancer (Day 7).
   - Ensure the Next.js frontend sends secure headers (set in `next.config.js`):
     - `Strict-Transport-Security` (HSTS)
     - `X-Content-Type-Options: nosniff`
     - `X-Frame-Options: DENY`
     - `Content-Security-Policy`
   - Use parameterized queries in Go (already handled by most ORMs like GORM or standard `database/sql`) to prevent SQL Injection.
   - **Rate Limiting**: Add rate limiting to Nginx to prevent API abuse and DDoS:
     ```nginx
     # In nginx.conf, add a rate-limiting zone
     limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;

     server {
         location /api/ {
             limit_req zone=api burst=20 nodelay;
             limit_req_status 429;
             proxy_pass http://localhost:8000/;
         }
     }
     ```
   - *Why rate limiting?* Without it, a single client can overwhelm your API with thousands of requests per second, starving legitimate users. The `burst=20` allows short spikes while `rate=10r/s` enforces the steady-state limit. Return `429 Too Many Requests` so clients can back off gracefully.

5. **Secret Management**
   - Never commit `.env` files.
   - In production, inject secrets via CI/CD (GitHub Secrets -> Terraform variables -> Droplet environment variables).

6. **Container Security**
   - Modify Dockerfiles to run as non-root users:
     ```dockerfile
     RUN addgroup -S appgroup && adduser -S appuser -G appgroup
     USER appuser
     ```
   - Use minimal base images like `alpine` or `scratch` to reduce the attack surface.

## Part 2: Disaster Recovery

Disaster Recovery (DR) is measured in RTO (Recovery Time Objective - how long to get back up) and RPO (Recovery Point Objective - how much data you can afford to lose).

### Implementation Steps

1. **Backup Strategy**
   - **Database**: DO Managed DBs automatically backup daily (7-day retention). For extra safety, configure a cron job to run `pg_dump` and push to Spaces.
   - **Application**: Code is safely in GitHub. Built images are in DO Container Registry.
   - **Infrastructure**: Terraform state should be backed up (e.g., using Spaces as a remote backend).

2. **Recovery Procedures (Runbooks)**
   - *Scenario A: Droplet is accidentally destroyed.*
     - RTO: ~5 minutes.
     - Procedure: Run `terraform apply`. Terraform spins up a new Droplet. GitHub Actions redeploys the latest container images via SSH/Docker Compose.
   - *Scenario B: Database corruption.*
     - RTO: ~15 minutes. RPO: Maximum 24 hours (based on daily backup).
     - Procedure: Use DO Control Panel to "Restore from Backup" to a new cluster. Update Terraform with the new cluster ID. Apply.

3. **Snapshot Strategy**
   - For stateful Droplets, you can configure automatic weekly snapshots via the DO UI or Terraform.

### Bonus: HA PostgreSQL Experiment (Use Those Expiring Credits!)

Since credits are expiring, this is the perfect time to test **High Availability PostgreSQL**. In production, a HA database cluster has a standby node that automatically promotes to primary if the primary fails.

**Steps:**
1. In DO Console, resize your existing Managed PostgreSQL cluster to add a **standby node** (HA configuration).
   - Or create a brand-new HA cluster: `db-s-1vcpu-2gb`, 2 nodes (primary + standby), same VPC.
   - Cost: ~ ( for HA minus your existing $25 single node, prorated to pennies for a few hours).

2. **Observe the HA topology** in the DO Console — you'll see Primary and Standby nodes.

3. **Test Failover:**
   - Click "Promote Standby" in the DO Console (or use `doctl databases migrate <db-id>`).
   - Watch your application's connection momentarily drop and reconnect to the new primary.
   - Measure the failover time — typically 15-30 seconds.
   - Check: Did your Go app's connection pool reconnect automatically? (It should, if using `database/sql` with proper retry logic.)

4. **Observe replication lag** in the DO metrics dashboard during active writes.

5. **Create a Read Replica** (separate from HA) and configure your Go app to send analytics/read-heavy queries to the replica while writes go to the primary.

> **Why this matters:** Every production SaaS database should be HA. Testing failover *before* an actual incident is the difference between a 30-second blip and an hours-long outage. You'll never get a cheaper opportunity to practice this.

**Cleanup:** You can destroy the HA/replica resources after experimenting, or keep them alive — credits are expiring anyway.

## Part 3: Performance Optimization

1. **Database Connection Pooling**: DO Managed DBs come with PgBouncer. Ensure your Go app connects to the PgBouncer port (usually 25061) rather than direct Postgres  to prevent connection exhaustion.
2. **Redis Caching**: Ensure expensive API calls (e.g., analytics aggregations) check Redis before hitting Postgres.
3. **CDN**: Ensure Spaces is configured with the DO CDN edge cache to serve static assets (images, user uploads) to reduce latency for global users.

## Part 4: Final Architecture Review

### The 10x Scale Consideration
If traffic grows 10x:
- Move from Docker Compose on a single Droplet to DOKS (Kubernetes) or App Platform Professional.
- Add a Read-Replica to the Managed PostgreSQL cluster to offload analytics queries.

### The 100x Scale Consideration
If traffic grows 100x:
- Multi-region deployment (Active/Active).
- Shard the database or move to a distributed SQL engine.
- Implement aggressive Redis cluster sharding.

## Part 5: Knowledge Consolidation

### Decision Framework: When to use what?

| Workload | Recommended DO Service | Reason |
| :--- | :--- | :--- |
| **Static Website** | App Platform (Static) | Free, CDN included, Gitops built-in. |
| **Standard Web App (API+UI)**| App Platform (Basic/Pro) | Best balance of cost, CI/CD, and zero-maintenance. |
| **Custom Network/Legacy OS** | Droplets | Full control over the kernel and network stack. |
| **Microservices Team** | DOKS (Kubernetes) | Standardized orchestration, helm ecosystem, autoscaling. |
| **Background Jobs** | DO Functions / App Platform Worker | Event-driven scaling, isolated from main web traffic. |

## Part 6: Final Cleanup and Credit Report

This is the end of the bootcamp. If this is just a learning environment, **destroy everything** to prevent ongoing charges.

```bash
# In your terraform directory
terraform destroy -auto-approve
```

Verify in the DigitalOcean console that:
1. All Droplets are gone.
2. Managed Databases are deleted.
3. Load Balancers are deleted.
4. Spaces buckets are emptied and deleted.
5. Container Registry is deleted.

## Exercises & Experiments

### Disaster Recovery Drills (Mandatory)
- [ ] **Droplet Failure Recovery**: Run `terraform destroy -target=digitalocean_droplet.web`, then run `terraform apply`. Measure RTO — how long until the site is back online?
- [ ] **Database Restore**: In the DO Console, restore Managed PostgreSQL from its latest automatic backup to a new cluster. Update your app's connection string. Verify data integrity.
- [ ] **Snapshot Restore**: Restore your Droplet snapshot from Day 7 to a brand-new Droplet. Verify the app boots and serves traffic.
- [ ] **Volume Restore**: If you created a Block Storage snapshot on Day 7, create a new volume from it. Mount it on a Droplet and verify the data is intact.
- [ ] **Terraform State Recovery**: Simulate losing your Terraform state: download the state file from Spaces, delete it from the remote backend, then re-upload it. Verify `terraform plan` shows no changes. *Alternative:* Delete the local `.terraform` directory and run `terraform init` to re-download state from Spaces.

### Security Drills
- [ ] **Security Audit**: Run an SSH audit tool (`ssh-audit`) against your Droplet IP. Fix any findings.
- [ ] **Rate Limit Test**: Use `ab` (Apache Benchmark) or `hey` to send 100 requests/second to your API. Verify Nginx returns `429 Too Many Requests` for requests exceeding the limit.
- [ ] **Firewall Audit**: Attempt to connect to Postgres port  and Redis port  from your local machine. Verify they're blocked by the Cloud Firewall.
- [ ] **Automatic Updates Verify**: Check `unattended-upgrades` logs: `cat /var/log/unattended-upgrades/unattended-upgrades.log`

### Final Documentation
- [ ] **Runbook Creation**: Write a markdown runbook in `docs/runbook.md` explaining how to deploy the app from scratch if a new engineer joins the team.
- [ ] **Architecture Document**: Create `docs/architecture.md` with the final production architecture diagram and service inventory.

## What I Learned
- Security must be baked in at multiple layers: Network (VPC/Firewalls), OS (SSH keys/fail2ban), and Application (Headers, Non-root containers).
- Infrastructure as Code (Terraform) is not just for provisioning; it is the ultimate Disaster Recovery tool.
- A true production architecture involves balancing cost, complexity, maintainability, and reliability.

---

# FINAL PRODUCTION ARCHITECTURE

This represents the ideal, highly available, secure production setup for CloudPulse on DigitalOcean.

```text
                                           INTERNET
                                              |
                                              v
                              +-------------------------------+
                              |    Cloudflare / Route53       |
                              |    (DNS & DDoS Protection)    |
                              +-------------------------------+
                                              |
      +-------------------------------------------------------------------------------+
      | DIGITALOCEAN REGION (NYC3)                                                    |
      |                                                                               |
      |                               +---------------+                               |
      |                               |  DO Managed   |                               |
      |                               | Load Balancer |                               |
      |                               +---------------+                               |
      |                                /             \                                |
      |                              /                 \                              |
      |       +---------------------------------------------------------------+       |
      |       | VPC (Virtual Private Cloud) - 10.116.0.0/20                   |       |
      |       |                                                               |       |
      |       |    +-------------------+           +-------------------+      |       |
      |       |    |   App Droplet 1   |           |   App Droplet 2   |      |       |
      |       |    |-------------------|           |-------------------|      |       |
      |       |    | - Nginx / Caddy   |           | - Nginx / Caddy   |      |       |
      |       |    | - Next.js (FE)    |           | - Next.js (FE)    |      |       |
      |       |    | - Go API (BE)     |           | - Go API (BE)     |      |       |
      |       |    | - Go Worker       |           | - Go Worker       |      |       |
      |       |    +-------------------+           +-------------------+      |       |
      |       |             |                                |                |       |
      |       |             |          Internal VPC          |                |       |
      |       |             +--------------------------------+                |       |
      |       |             |                                |                |       |
      |       |    +-------------------+           +-------------------+      |       |
      |       |    | DO Managed        |           | DO Managed        |      |       |
      |       |    | PostgreSQL        |           | Redis             |      |       |
      |       |    | (Primary)         |           | (Cache & Pub/Sub) |      |       |
      |       |    +-------------------+           +-------------------+      |       |
      |       |             |                                                 |       |
      |       |    +-------------------+                                      |       |
      |       |    | PostgreSQL        |                                      |       |
      |       |    | (Standby Node)    |                                      |       |
      |       |    +-------------------+                                      |       |
      |       +---------------------------------------------------------------+       |
      |                                                                               |
      |   +-------------------+  +-------------------+  +-------------------+         |
      |   | DO Spaces (S3)    |  | DO Container      |  | DO Serverless     |         |
      |   | w/ Built-in CDN   |  | Registry (Basic)  |  | Functions         |         |
      |   | (Static Assets)   |  | (Docker Images)   |  | (Cron / Webhooks) |         |
      |   +-------------------+  +-------------------+  +-------------------+         |
      +-------------------------------------------------------------------------------+

        Out-of-band:
        - GitHub Actions (CI/CD -> pushes to Registry, triggers deployments)
        - Prometheus / Grafana / Loki (Monitoring Stack - ideally on separate admin droplet)
```

---

# COMPLETE SERVICE REFERENCE

| Service | What it does | When to use it | When NOT to use it | Pricing (Starting) | Learned |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Droplets** | Virtual Machines (IaaS) | Custom OS needs, legacy apps, maximizing performance/cost ratio. | You don't want to manage OS patches and systemd. |  | Day 1 |
| **VPC** | Private Networking | Securing traffic between DO resources, preventing public internet access. | N/A (Always use it). | Free | Day 3 |
| **Cloud Firewalls** | Network Security Group | Blocking ports, restricting IP access to Droplets. | You want application-layer (WAF) blocking. | Free | Day 3 |
| **Managed Databases** | Hosted PostgreSQL/MySQL/Redis | Production data storage, automated backups, high availability. | Scratch data, extremely tight budgets where you run DB on the app Droplet. |  | Day 5 |
| **Spaces (Object Storage)** | S3-compatible file storage | User uploads, backups, static web assets (CDN). | Fast random read/write disk access (use Block Storage). |  | Day 5 |
| **Container Registry** | Private Docker Image Host | Storing built application images securely for deployments. | Using public Docker Hub images exclusively. |  | Day 6 |
| **Load Balancers** | Distributes web traffic | Scaling across multiple Droplets, high availability, TLS termination. | Single-droplet setups. |  | Day 8 |
| **DOKS (Kubernetes)** | Managed K8s Cluster | Microservices, advanced rolling deployments, auto-scaling. | Small teams, monolithic apps, strict budgets. |  | Day 9 |
| **App Platform** | Platform as a Service (PaaS) | Zero-maintenance app hosting, auto CI/CD from Git. | Applications requiring custom kernel tuning or specific background daemons. |  | Day 10 |
| **DO Functions** | Serverless Compute | Event-driven code (Spaces uploads, cron jobs, webhooks). | Long-running processes, predictable heavy loads. | Free (25k) | Day 10 |
---

# WHAT'S NEXT?

You have built a production-grade infrastructure on DigitalOcean! Where do you go from here?

1. **Advanced Topics to Explore:**
   - **Infrastructure Testing**: Write `terratest` to validate your Terraform code automatically.
   - **Advanced Observability**: Implement OpenTelemetry distributed tracing across your Go backend and Next.js frontend.
   - **GitOps**: If you choose the Kubernetes route, learn ArgoCD or Flux for declarative cluster management.

2. **Alternative Cloud Platforms:**
   - Now that you understand the primitives (Compute, VPC, Managed DBs, Object Storage, IAM), applying this to AWS, GCP, or Azure is mostly translating vocabulary (e.g., Droplet -> EC2, Spaces -> S3, App Platform -> AWS AppRunner).
   - *Challenge*: Try recreating the Terraform configuration for AWS.

3. **Community & Certifications:**
   - Read the DigitalOcean community tutorials; they are among the best technical documentation on the internet.
   - Consider pursuing the *Cloud Native Computing Foundation (CNCF)* certifications like CKA (Certified Kubernetes Administrator) if you want to specialize in Kubernetes.

Congratulations on completing the Cloud Engineering Bootcamp! You now have the skills to architect, provision, secure, and monitor real-world web applications.
