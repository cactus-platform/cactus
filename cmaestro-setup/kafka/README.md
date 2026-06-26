# Cactus Maestro - Kafka Configuration

This setup comes from https://github.com/adobe/koperator/  

Adobe documentation : https://opensource.adobe.com/koperator/docs/e2e-tutorial/koperator-install/


## Configuration Steps

> helm install zookeeper-operator --repo https://charts.pravega.io zookeeper-operator --namespace=cmaestro-kafkasys --create-namespace

```bash
kubectl create -f - <<EOF
apiVersion: zookeeper.pravega.io/v1beta1
kind: ZookeeperCluster
metadata:
    name: zookeeper-server
    namespace: cmaestro-kafkasys
spec:
    replicas: 1
    image:
        repository: ghcr.io/adobe/zookeeper-operator/zookeeper
        tag: 3.8.4-0.2.15-adobe-20250923
    persistence:
        reclaimPolicy: Delete
EOF
```

> kubectl apply -f https://raw.githubusercontent.com/adobe/koperator/refs/heads/master/config/base/crds/kafka.banzaicloud.io_cruisecontroloperations.yaml  
> kubectl apply -f https://raw.githubusercontent.com/adobe/koperator/refs/heads/master/config/base/crds/kafka.banzaicloud.io_kafkaclusters.yaml  
> kubectl apply -f https://raw.githubusercontent.com/adobe/koperator/refs/heads/master/config/base/crds/kafka.banzaicloud.io_kafkatopics.yaml  
> kubectl apply -f https://raw.githubusercontent.com/adobe/koperator/refs/heads/master/config/base/crds/kafka.banzaicloud.io_kafkausers.yaml  

Check the latest version at : https://github.com/orgs/adobe/packages  

> helm install kafka-operator oci://ghcr.io/adobe/helm-charts/kafka-operator --version 0.28.0-adobe-20250923 --skip-crds --namespace=cmaestro-kafkasys --create-namespace

> helm pull oci://ghcr.io/adobe/koperator/kafka-operator --version 0.28.0-adobe-20250923