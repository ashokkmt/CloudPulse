# CloudPulse Development Environment

CloudPulse is a robust application designed to help master DigitalOcean cloud services. Below is the complete listing of all services available in the local development environment using Docker Compose.

## Available Services

### Application Services
* **Frontend**
  * **URL**: http://localhost:3000
  * **Purpose**: Next.js user interface.
* **Backend API**
  * **URL**: http://localhost:8000
  * **Purpose**: Go REST API.
* **MinIO (DigitalOcean Spaces Simulator)**
  * **API URL**: http://localhost:9000
  * **Console URL**: http://localhost:9001
  * **Credentials**: `minioadmin` / `minioadmin`
  * **Purpose**: S3-compatible object storage simulating DO Spaces.

### Database & Cache Management
* **PostgreSQL Admin UI (pgAdmin 4)**
  * **URL**: http://localhost:5050
  * **Credentials**: `admin@cloudpulse.com` / `admin`
  * **Database Connection Details**: 
    * Host: `db`
    * Port: `5432`
    * Username: `postgres`
    * Password: `postgrespassword`
  * **Purpose**: Web-based PostgreSQL management tool for inspecting the `tasks` and `users` tables.
* **Redis Admin UI (RedisInsight)**
  * **URL**: http://localhost:8001
  * **Connection Details**: Click "Add Redis Database" and use Host: `redis` and Port: `6379`.
  * **Purpose**: Real-time Redis browser for inspecting cache keys, TTLs, and observing cache invalidation.

### Observability & Monitoring Stack
* **Grafana**
  * **URL**: http://localhost:3001
  * **Credentials**: `admin` / `admin` (or configured via `.env.local`)
  * **Purpose**: Centralized dashboarding for metrics, logs, and traces.
* **Prometheus**
  * **URL**: http://localhost:9090
  * **Purpose**: Time-series metrics collection and querying.
* **Tempo**
  * **URL**: http://localhost:3200
  * **Purpose**: Distributed tracing backend (OpenTelemetry).
* **Loki**
  * **URL**: http://localhost:3100
  * **Purpose**: Log aggregation system.
* **Alertmanager**
  * **URL**: http://localhost:9093
  * **Purpose**: Handles alerts sent by Prometheus.

## Development Workflow

To start the full development environment, run:

```bash
# 1. Start the main application stack
docker-compose --env-file ./backend/.env.local up -d --build

# 2. Start the observability stack
docker-compose -f docker-compose.monitoring.yml up -d
```

You can now open the respective browser tabs to verify the application's behavior.
