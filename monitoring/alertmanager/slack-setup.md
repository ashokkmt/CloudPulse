# Slack Alertmanager Setup

To configure Alertmanager to send notifications to your Slack workspace, follow these steps:

1. **Create a Slack App:**
   - Go to [https://api.slack.com/apps](https://api.slack.com/apps) and click "Create New App".
   - Choose "From scratch", give it a name (e.g., "CloudPulse Alerts"), and select your workspace.

2. **Enable Incoming Webhooks:**
   - In the app settings menu, navigate to "Incoming Webhooks".
   - Toggle the switch to "On".

3. **Add New Webhook to Workspace:**
   - Click "Add New Webhook to Workspace" at the bottom of the page.
   - Select the channel where you want the alerts to be posted (e.g., `#alerts`) and click "Allow".

4. **Copy the Webhook URL:**
   - You will see a new Webhook URL. Treat it as a secret and do not commit it to git.
   - Copy this URL.

5. **Update Alertmanager Configuration:**
   - Open `alertmanager.yml` located in this directory.
   - Replace the `api_url` placeholder with your actual webhook URL.

```yaml
receivers:
  - name: 'slack-notifications'
    slack_configs:
      - api_url: '<YOUR_WEBHOOK_URL>'
        channel: '#alerts'
        send_resolved: true
```

6. **Restart Alertmanager:**
   - Apply the changes by restarting the Alertmanager container from the project root:
     `docker-compose -f docker-compose.monitoring.yml restart alertmanager`
