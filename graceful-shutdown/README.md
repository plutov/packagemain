[Read the full article on packagemain.tech](https://packagemain.tech/p/graceful-shutdowns-k8s-go)

## Requirements

1. Kubernetes cluster. For example using colima: `colima start --kubernetes --network-address`
2. `kubectl` command installed
3. `docker` command installed

## Build Docker Images

We use 2 tags (v1 and v2) to initiate a Rolling Update later.

```
docker build -t server:v1 -t server:v2 .
```

## Deploy Redis to k8s

```
kubectl create namespace redis
helm install redis oci://registry-1.docker.io/bitnamicharts/redis --set auth.enabled=false --namespace redis
```

## Deploy Go server

```
kubectl apply -f server.yaml
```

## Send requests

Use [vegeta](https://github.com/tsenart/vegeta) to send 3000 requests (50 per second for 60s):

```
# reset redis counter
kubectl exec -it redis-master-0 -n redis -- redis-cli set counter 0

echo "GET http://localhost:30001/incr" | vegeta attack -duration=60s -rate=50 -http2=false | vegeta report
```

## Test with Kubernetes Rolling Update

While vegeta is running we can initiate a Rolling Update in Kubernetes by changing the image tag.

```
kubectl set image deployment server server=server:v2
```

## Verify counter in Redis

It should be 3000.

```
kubectl exec -it redis-master-0 -n redis -- redis-cli get counter
```

## Clean up

```
kubectl delete -f server.yaml
helm uninstall redis --namespace redis
kubectl delete namespace redis
```
