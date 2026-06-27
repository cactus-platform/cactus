PASSWORD=$(openssl rand -base64 48)

echo "PASSWORD = [$PASSWORD]"

kubectl create secret generic postgres-credentials \
  --namespace cmaestro-db \
  --from-literal=POSTGRES_PASSWORD="$PASSWORD"