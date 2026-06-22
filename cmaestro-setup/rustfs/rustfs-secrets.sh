kubectl -n cmaestro-db create secret generic seaweedfs-s3-credentials \
  --from-literal=accessKey="$(openssl rand -hex 16)" \
  --from-literal=secretKey="$(openssl rand -base64 48)"