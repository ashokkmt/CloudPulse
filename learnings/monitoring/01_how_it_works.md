# Complete Guide to CloudPulse Observability

This document provides a deep dive into the "Three Pillars of Observability" (Metrics, Logs, and Traces), explaining exactly how the industry-standard open-source tools work in general, and precisely how they are integrated into our CloudPulse application.

---

## 1. Prometheus (The Metrics Engine)

### How it works in general:
Prometheus is a **Time-Series Database (TSDB)**. Unlike a traditional SQL database that stores rows of data, Prometheus stores numbers over time. 
*   **Pull Model:** Instead of your application pushing data to a database, Prometheus uses a "pull" model. It reaches out (scrapes) a specific HTTP endpoint (usually `/metrics`) on your servers every few seconds to gather data.
*   **Metric Types:**
    *   *Counter:* A number that only goes up (e.g., total website visits).
    *   *Gauge:* A number that can go up and down (e.g., current active users, CPU usage).
    *   *Histogram:* A way to measure the distribution of events into buckets (e.g., how many requests took less than 0.1s, 0.5s, 1s).

### How it works in CloudPulse:
1.  **Generation:** Inside the Go backend (`internal/middleware/metrics.go`), we imported the Prometheus Go client. Every time an HTTP request hits the backend, the `MetricsMiddleware` increments a Counter (`http_requests_total`) and updates a Histogram (`http_request_duration_seconds`). It then hosts this raw data at `http://localhost:8000/metrics`.
2.  **Scraping:** Our `monitoring/prometheus/prometheus.yml` is configured with four "jobs":
    *   `app`: Scrapes the Go backend.
    *   `node`: Scrapes `node-exporter` (a tool running on the server that measures the host machine's RAM, CPU, and Disk).
    *   `cadvisor`: Scrapes `cAdvisor` (a tool that monitors how much CPU/RAM each individual Docker container is using).
    *   `prometheus`: Scrapes itself to monitor its own health.
3.  **Storage:** Every 15 seconds, Prometheus grabs the numbers from all four sources and saves them in its internal database.

---

## 2. Grafana (The Visualization Layer)

### How it works in general:
Raw time-series data looks like a wall of text. Grafana is a powerful visualization web application. It doesn't store data; instead, it connects to databases (Datasources) like Prometheus, Loki, or even PostgreSQL, runs queries against them, and draws beautiful, real-time charts. 
It uses **PromQL** (Prometheus Query Language) to do math on the metrics (like calculating the rate of change over 5 minutes).

### How it works in CloudPulse:
*   **Provisioning:** Usually, you have to click through the Grafana UI to add datasources and build dashboards. To make our setup production-ready and reproducible, we use "provisioning". The files inside `monitoring/grafana/provisioning/` tell Grafana to automatically load Prometheus, Loki, and Tempo on startup.
*   **Dashboards:** The JSON files inside `monitoring/grafana/dashboards/` contain the exact queries for our charts. For example, the Application Dashboard runs the query `rate(http_requests_total[5m])`. Grafana asks Prometheus: *"How fast did the 'http_requests_total' counter grow over the last 5 minutes?"* and draws the Requests-Per-Second (RPS) line graph.

---

## 3. Loki (The Logging Engine)

### How it works in general:
Logs are the raw text output of your application (`console.log()` or `fmt.Println()`). Traditional systems like Elasticsearch index every single word of a log message so you can search it quickly, but this requires massive amounts of RAM and CPU.
Loki takes a different approach: it is "Prometheus, but for logs." Instead of indexing the full text, it only indexes **labels** (metadata like `container="backend"`, `level="error"`, or `region="us-east"`). The log text itself is compressed and stored cheaply.

### How it works in CloudPulse:
*   **Local Development:** We run a small agent called **Promtail** (defined in `docker-compose.monitoring.yml`). Promtail mounts to the Docker socket, reads the console output of all your running containers, attaches labels to them (like the container name), and pushes them to Loki.
*   **Production Deployment:** On the Droplet, we configure the Docker daemon itself to use the `loki` logging driver. Docker automatically streams `stdout` and `stderr` natively to the Loki container, completely removing the need for the Promtail middleman.
*   **Viewing:** You open the "Explore" tab in Grafana, select Loki, and run a **LogQL** query like `{container="backend"} |= "timeout"`. Grafana fetches those specific logs from Loki.

---

## 4. Tempo (Distributed Tracing)

### How it works in general:
If a user clicks a button on the frontend, and it takes 10 seconds to load, metrics will tell you the request was slow, but they won't tell you *where* the time went. Logs might show a database error, but they don't easily link back to the specific user's click.
**Tracing** solves this. When a request enters your system, it is given a unique `Trace ID`. As the request moves from the Frontend -> API Gateway -> Backend -> Database, it creates "Spans" that record exactly how long each step took. Tempo is the database that stores these spans.

### How it works in CloudPulse:
Our `tempo-config.yml` starts a Tempo server listening for OTLP (OpenTelemetry Protocol) traffic on port 4317. 
As the application grows, the Go backend will be instrumented with the OpenTelemetry SDK. When a slow request occurs, you will see a log line in Loki with a `Trace ID`. In Grafana, you simply click that ID, and Tempo instantly renders a "waterfall chart" showing the exact millisecond breakdown of the request's journey.

---

## 5. Alertmanager (Automated Notifications)

### How it works in general:
Prometheus is smart enough to know when metrics look bad, but it isn't designed to send emails or Slack messages. It also isn't smart enough to know that if a server completely crashes, 50 different microservices on that server will fail simultaneously, generating 50 alerts.
**Alertmanager** takes raw alert signals from Prometheus, groups related alerts together (so you get one Slack message instead of 50), and routes them to the correct team (e.g., Database errors go to the DBA team's Slack, Frontend errors go to PagerDuty).

### How it works in CloudPulse:
1.  **The Rules:** In `monitoring/prometheus/rules.yml`, we wrote three rules using PromQL:
    *   `InstanceDown`: Checks if `up == 0` (meaning a server or container is dead).
    *   `HighCPUUsage`: Checks if the Droplet CPU is over 85%.
    *   `HighErrorRate`: Checks if the ratio of 500-level errors to total requests is over 5%.
2.  **The Trigger:** If Prometheus sees the CPU hit 90%, it waits 5 minutes (to ensure it's not a temporary micro-spike). After 5 minutes, it sends a signal to Alertmanager.
3.  **The Notification:** Alertmanager receives the signal, formats it into a readable message, and makes an HTTP POST request to the Slack Webhook URL we defined in `alertmanager.yml` (as detailed in `slack-setup.md`). You get pinged on your phone instantly.
