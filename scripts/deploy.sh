#!/bin/sh

set -e

echo "Starting Kubernetes cluster..."
minikube start

echo "Building API docker image..."
minikube image build -t deployments-api:latest -f deployments/Dockerfile.api .

echo "Building Postgres docker image..."
minikube image build -t deployments-postgres:latest -f deployments/Dockerfile.postgres .

echo "Applying deployment secrets..."
kubectl apply -f deployments/k8s/secrets.yml

echo "Installing/updating Traefik..."
helm repo add traefik https://traefik.github.io/charts
helm repo update
helm upgrade --install traefik traefik/traefik \
  --namespace traefik \
  --create-namespace \
  -f deployments/k8s/traefik-values.yml

echo "Installing/updating VictoriaMetrics observability stack..."
helm repo add vm https://victoriametrics.github.io/helm-charts/
helm repo update
helm upgrade --install vmks vm/victoria-metrics-k8s-stack \
  --namespace monitoring \
  --create-namespace \
  -f deployments/k8s/victoria-values.yml

helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
helm repo update
helm upgrade --install otel-gateway open-telemetry/opentelemetry-collector \
  --namespace monitoring \
  -f deployments/k8s/otel-gateway-values.yml

echo "Deploying Postgres..."
kubectl apply -f deployments/k8s/postgres.yml

echo "Deploying database migration..."
kubectl apply -f deployments/k8s/migrate.yml

echo "Deploying Redis..."
kubectl apply -f deployments/k8s/redis.yml

echo "Deploying API..."
kubectl apply -f deployments/k8s/api.yml

echo "Deploying API ingress..."
kubectl apply -f deployments/k8s/api-ingress.yml

echo "Deployment Complete"
echo ""
echo "Kubernetes resources:"
kubectl get pods
echo ""
kubectl get svc
echo ""
kubectl get ingress
