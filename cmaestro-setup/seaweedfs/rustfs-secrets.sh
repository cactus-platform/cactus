ACCESS_KEY=$(openssl rand -hex 16)
SECRET_KEY=$(openssl rand -base64 48)

echo "ACCESS_KEY = [$ACCESS_KEY]"
echo "SECRET_KEY = [$SECRET_KEY]"

kubectl -n cmaestro-db create secret generic seaweedfs-s3-credentials \
  --from-literal=accessKey="$ACCESS_KEY" \
  --from-literal=secretKey="$SECRET_KEY"