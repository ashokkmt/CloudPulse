# Kubernetes Manual Deployment Guide

A quick reference cheatsheet of commands for deploying and managing CloudPulse on DOKS.

### Cluster Setup
```bash
doctl kubernetes cluster kubeconfig save cloudpulse-k8s
doctl kubernetes cluster registry add cloudpulse-k8s
```

### Deploy Application
```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/secret.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/backend.yaml
kubectl apply -f k8s/worker.yaml
kubectl apply -f k8s/frontend.yaml
```

### Install Ingress
```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.11.3/deploy/static/provider/do/deploy.yaml
kubectl apply -f k8s/ingress.yaml
kubectl get svc -n ingress-nginx  # get public IP
```

### Deploy Monitoring
```bash
helm install cloudpulse-monitoring prometheus-community/kube-prometheus-stack --namespace monitoring
helm install loki grafana/loki-stack --namespace monitoring --set promtail.enabled=true
helm install tempo grafana/tempo --namespace monitoring
```

### Debugging Cheatsheet
```bash
kubectl get pods -n cloudpulse
kubectl logs <pod> -n cloudpulse
kubectl logs <pod> -n cloudpulse --previous
kubectl describe pod <pod> -n cloudpulse
kubectl get events -n cloudpulse --sort-by='.lastTimestamp'
kubectl exec -it <pod> -n cloudpulse -- /bin/sh
```

### Scaling & Management
```bash
kubectl scale deployment cloudpulse-backend --replicas=5 -n cloudpulse
kubectl rollout status deployment/cloudpulse-backend -n cloudpulse
kubectl rollout undo deployment/cloudpulse-backend -n cloudpulse
```

### Update Image (New Deployment)
```bash
kubectl set image deployment/cloudpulse-backend backend=registry.digitalocean.com/cloudpulse-registry/backend:<NEW_SHA> -n cloudpulse
```

### Cleanup
```bash
kubectl delete -f k8s/
helm uninstall cloudpulse-monitoring -n monitoring
helm uninstall loki -n monitoring
helm uninstall tempo -n monitoring
```
