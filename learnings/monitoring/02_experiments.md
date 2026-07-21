# Monitoring Experiments

The best way to understand the monitoring stack is to break things and watch the dashboards react. 

Here are 3 experiments you can run right now to see the monitoring stack in action.

## Experiment 1: Generate a Traffic Spike (RPS)
Your backend is already generating metrics, but if there are no users, the graphs will be flat. Let's create an artificial spike in traffic.

**Steps:**
1. Open a new terminal.
2. Run the included Node.js script to simulate a flood of traffic:
   ```bash
   node learnings/monitoring/generate_traffic.js
   ```
3. Open Grafana (`http://localhost:3001`) and navigate to the **Application Dashboard**.
4. **Observe:** Within 15-30 seconds, you should see the **RPS (Requests Per Second)** graph spike upwards dramatically. 

*What happened?* The script sent hundreds of requests to the backend. The backend updated its internal counters. Prometheus scraped those counters 15 seconds later. Grafana queried Prometheus and drew the spike.

## Experiment 2: Trigger a High Error Rate Alert
Let's simulate a broken endpoint to see if we can trigger an alert.

**Steps:**
1. Open `backend/internal/handler/router.go` and temporarily add a broken route that always returns a 500 Internal Server Error:
   ```go
   mux.HandleFunc("/api/broken", func(w http.ResponseWriter, r *http.Request) {
       http.Error(w, "Simulated Server Error", http.StatusInternalServerError)
   })
   ```
2. Restart your backend container.
3. Modify the `generate_traffic.js` script to hit `http://localhost:8000/api/broken` instead of the health check.
4. Run the traffic script for about 2 minutes.
5. Open Grafana and watch the **Error Rate** graph spike above 5%.
6. Open Alertmanager (`http://localhost:9093`).
7. **Observe:** You should see a red `HighErrorRate` alert active in the Alertmanager UI.

*What happened?* The script generated 500-level errors. The Prometheus rule evaluated `rate(http_requests_total{status=~"5.."}[2m])` and saw it exceeded the 5% threshold, firing the alert.

## Experiment 3: Find the Most CPU-Intensive Container
Let's see how infrastructure metrics from `cAdvisor` help us identify resource hogs.

**Steps:**
1. Open Grafana (`http://localhost:3001`).
2. Go to **Explore** (the compass icon on the left menu).
3. Ensure the **Prometheus** datasource is selected.
4. Enter the following PromQL query to see the CPU usage of all containers:
   ```promql
   rate(container_cpu_usage_seconds_total{name!=""}[1m])
   ```
5. Click **Run query**.
6. **Observe:** You will see a list of lines representing every Docker container currently running on your machine.
7. Now, start the `generate_traffic.js` script again. Re-run the query. You should see the line corresponding to the `backend` container jump significantly higher as it struggles to process the flood of requests.
