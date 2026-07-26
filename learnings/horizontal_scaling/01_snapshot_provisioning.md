# Bootstrapping a Droplet from a Snapshot

When you create `cloudpulse-app-02` from a snapshot of `cloudpulse-app-01`, you are making an exact byte-for-byte clone of the disk at that specific moment in time.

## Why is it broken immediately after booting?
When the new Droplet boots up, the Docker daemon starts automatically. Docker sees that containers were running when the snapshot was taken, so it immediately tries to start them again.

However, you can't reach `localhost:8000/healthz` because:
1. **Missing Environment Variables**: When CI/CD deployed to Droplet 1, it executed `APP_VERSION=<sha> docker compose up -d`. Droplet 2 doesn't know this environment variable, so if the containers crashed or restarted, Docker Compose might be trying to pull/run an invalid tag or falling back to a missing `latest` tag.
2. **Stale State**: The containers might be in a crash-loop due to IP changes or networking state captured during the snapshot.

## How to Fix It (Step-by-Step)

To get Droplet 2 into a healthy state, you need to manually clean up the "ghost" containers and start them fresh with the correct version tag.

### 1. SSH into the new Droplet
```bash
ssh deploy@<DROPLET_02_IP>
```

### 2. Navigate to the App Directory
```bash
cd ~/cloudpulse
```

### 3. Wipe the Stale Containers
Force Docker Compose to stop and remove the cloned containers that started automatically:
```bash
docker compose down
```

### 4. Re-Authenticate with DigitalOcean Container Registry (DOCR)
Your `doctl` authentication token might still be valid from the snapshot, but it's always best to ensure you can pull private images:
```bash
doctl registry login
```

### 5. Start the Containers with the Correct Version
Since you aren't running this via GitHub Actions right now, you need to provide the `APP_VERSION` environment variable manually. 

You have two options:

**Option A (Using a specific Git commit SHA - Recommended):**
Look at your DigitalOcean Container Registry to find the latest valid SHA tag, then run:
```bash
export APP_VERSION=a1b2c3d  # Replace with your actual commit SHA from DOCR
docker compose pull
docker compose up -d
```

**Option B (Fallback to Latest):**
If you have a `latest` tag in your registry, you can use that (though exact SHAs are better for production):
```bash
export APP_VERSION=latest
docker compose pull
docker compose up -d
```

### 6. Verify Health
Once the containers are recreated, verify the backend is running correctly:
```bash
docker ps
curl http://localhost:8000/healthz
```

You should receive a `{"status": "ok"}` response.

---

## Future Automation
In a fully automated setup (like using Terraform or DigitalOcean User Data scripts), you would automate this cleanup process by passing a startup script to the Droplet that automatically runs `docker compose down` and pulls the latest image upon its first boot. For now, doing it manually helps you understand exactly what the snapshot cloned!
