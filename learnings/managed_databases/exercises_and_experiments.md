# Managed Databases Exercises and Experiments

This guide provides practical exercises and experiments to understand DigitalOcean Managed Databases (PostgreSQL and Redis). Since you have secured your databases by restricting access exclusively to your App Droplet via Trusted Sources, **all command-line operations must be performed from within your App Droplet**.

## Preparation: Accessing Your App Droplet

Before starting, SSH into your App Droplet. This is the only machine authorized to connect to your managed databases.

```bash
ssh deploy@<YOUR_APP_DROPLET_IP>
```

---

## Part 1: Exercises

### 1. Slow Query Log

DigitalOcean provides a Logs & Queries dashboard to identify slow-running queries that might be bottlenecking your application.

**Steps:**
1. From your App Droplet, connect to your managed PostgreSQL database using `psql`:
   ```bash
   psql "postgres://<user>:<password>@<host>:<port>/<dbname>?sslmode=require"
   ```
   *(Note: Use the connection string provided in your DigitalOcean dashboard).*
2. Execute a simulated slow query:
   ```sql
   SELECT pg_sleep(5);
   ```
3. Log into the DigitalOcean Control Panel.
4. Navigate to your Managed Database -> **Logs & Queries** tab.
5. Look for the `SELECT pg_sleep(5);` query to appear in the slow query logs. (It may take a few minutes to aggregate).

### 2. Connection Pooling

Managed PostgreSQL includes PgBouncer for connection pooling, running on port `25061`. The direct PostgreSQL connection runs on `25060`. We will use `pgbench` to see how the pooler handles connection spikes.

**Steps:**
1. Make sure `pgbench` is installed on your App Droplet:
   ```bash
   sudo apt-get update
   sudo apt-get install postgresql-contrib
   ```
2. Initialize the `pgbench` tables (using the direct port `25060`):
   ```bash
   pgbench -i -s 1 "postgres://<user>:<password>@<host>:25060/cloudpulse?sslmode=require"
   ```
3. **Test Direct Connection (Port 25060):** Run 200 concurrent clients.
   ```bash
   pgbench -c 200 -j 2 -t 10 "postgres://<user>:<password>@<host>:25060/cloudpulse?sslmode=require"
   ```
   *Observation: This will likely fail or throw connection limit errors because the default PostgreSQL max connections is much lower than 200.*
4. **Test Connection Pooler (Port 25061):** Run the same 200 clients through PgBouncer.
   ```bash
   pgbench -c 200 -j 2 -t 10 "postgres://<user>:<password>@<host>:25061/cloudpulse?sslmode=require"
   ```
   *Observation: This succeeds smoothly. PgBouncer queues the connections and feeds them to the backend PostgreSQL limits without crashing.*

### 3. Read Replica

Read replicas allow you to scale your read-heavy workloads geographically.

**Steps:**
1. In the DigitalOcean Control Panel, navigate to your Managed Database.
2. Click **Add a standby node or read-only node**.
3. Choose a region and create a **Read-Only Node**. Wait for it to provision.
4. Get the connection string specifically for the Read Replica from the DO dashboard.
5. On your App Droplet, connect to the replica:
   ```bash
   psql "postgres://<user>:<password>@<replica-host>:<port>/<dbname>?sslmode=require"
   ```
6. Verify reads work:
   ```sql
   SELECT * FROM tasks LIMIT 5;
   ```
7. Verify writes are blocked:
   ```sql
   INSERT INTO tasks (title, done) VALUES ('Replica Test', false);
   ```
   *Expected Error: `cannot execute INSERT in a read-only transaction`*
8. **Cleanup:** Delete the read replica in the DO dashboard to save money.

### 4. Backup Restoration

DigitalOcean takes daily backups of your managed database automatically.

**Steps:**
1. In the DO Control Panel, go to your Managed Database.
2. Click the **Actions** dropdown in the top right, and select **Restore from backup**.
3. Select a time from yesterday (or the earliest available).
4. Provide a name for the new restored database cluster (e.g., `cloudpulse-restored`).
5. Click **Restore Database**. DigitalOcean will spin up an entirely new database cluster containing the data as it existed at that exact point in time.

---

## Part 2: Experiments

### 1. Simulate Connection Exhaustion

Let's write a Go script to aggressively open 1000 connections and leave them sleeping to see how the connection pool behaves.

**Steps:**
1. On your App Droplet, create a file named `exhaust_conns.go`:
   ```bash
   nano exhaust_conns.go
   ```
2. Paste the following Go code:
   ```go
   package main

   import (
       "database/sql"
       "fmt"
       "os"
       "sync"

       _ "github.com/lib/pq"
   )

   func main() {
       dbURL := os.Getenv("DATABASE_URL")
       if dbURL == "" {
           fmt.Println("Please set DATABASE_URL")
           return
       }

       var wg sync.WaitGroup
       conns := 1000 // Attempt to open 1000 connections

       for i := 0; i < conns; i++ {
           wg.Add(1)
           go func(id int) {
               defer wg.Done()
               db, err := sql.Open("postgres", dbURL)
               if err != nil {
                   return
               }
               // Keep the connection open by executing a long sleep
               _, err = db.Exec("SELECT pg_sleep(30)")
               if err != nil {
                   fmt.Printf("Conn %d failed: %v\n", id, err)
               } else {
                   fmt.Printf("Conn %d succeeded\n", id)
               }
           }(i)
       }
       wg.Wait()
   }
   ```
3. Run the script against the **direct port (25060)**:
   ```bash
   export DATABASE_URL="postgres://<user>:<password>@<host>:25060/<dbname>?sslmode=require"
   go run exhaust_conns.go
   ```
   *Result:* Many connections will fail immediately with `too many clients already`.
4. Run the script against the **pooler port (25061)**:
   ```bash
   export DATABASE_URL="postgres://<user>:<password>@<host>:25061/<dbname>?sslmode=require"
   go run exhaust_conns.go
   ```
   *Result:* The script will block as PgBouncer cleanly queues the 1000 requests without overloading the database.

### 2. Kill App Connections

Test your Go application's resilience to sudden database disconnections.

**Steps:**
1. Ensure your Go backend (`cloudpulse` app) is running normally on the App Droplet and connected to the database.
2. In another terminal on the App Droplet, connect to PostgreSQL via `psql`.
3. Find your backend's Process IDs (PIDs):
   ```sql
   SELECT pid, application_name, state FROM pg_stat_activity WHERE datname = 'cloudpulse';
   ```
4. Terminate one or all of the backend's connections:
   ```sql
   SELECT pg_terminate_backend(<pid_from_above>);
   ```
5. Check your backend application logs (e.g., `docker logs cloudpulse-backend` or `journalctl -u cloudpulse`).
6. *Observation:* You should see a database disconnect error followed by an automatic reconnection (assuming your Go `database/sql` pool is properly handling reconnects, which it does natively).

### 3. Redis Memory Limit (Eviction Policy)

Your DO Managed Redis instance has a strict memory limit (e.g., 1GB). Let's see what happens when we exceed it.

**Steps:**
1. Ensure your Redis cluster eviction policy is set to `allkeys-lru` (the default for DO caching configurations). You can check this in the DO Settings tab for your Redis instance.
2. On your App Droplet, create `fill_redis.go`:
   ```bash
   nano fill_redis.go
   ```
3. Paste the following script, which continuously writes 1MB blocks to Redis:
   ```go
   package main

   import (
       "context"
       "crypto/rand"
       "fmt"
       "os"

       "github.com/redis/go-redis/v9"
   )

   func main() {
       redisURL := os.Getenv("REDIS_URL")
       opt, _ := redis.ParseURL(redisURL)
       client := redis.NewClient(opt)
       ctx := context.Background()

       // Generate 1MB of random data
       data := make([]byte, 1024*1024) 
       rand.Read(data)

       fmt.Println("Starting to fill Redis...")
       for i := 0; i < 2000; i++ { // Attempt to write 2GB total
           key := fmt.Sprintf("junk_data_%d", i)
           err := client.Set(ctx, key, data, 0).Err()
           if err != nil {
               fmt.Printf("Failed at %d MB: %v\n", i, err)
               break
           }
           if i%100 == 0 {
               fmt.Printf("Inserted %d MB\n", i)
           }
       }
   }
   ```
4. Run the script:
   ```bash
   export REDIS_URL="rediss://default:<password>@<host>:<port>"
   go run fill_redis.go
   ```
5. *Observation:* If the policy is `allkeys-lru`, Redis will silently start deleting older `junk_data_X` keys to make room for the new ones, and the script will finish successfully. If the policy is `noeviction`, the script will crash around ~1000MB with an `OOM command not allowed` error.
