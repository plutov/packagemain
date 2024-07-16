[Read the full article on packagemain.tech](https://packagemain.tech/p/graceful-shutdowns-k8s-go)

## Build Docker Images

We use 2 tags (v1 and v2) to initiate a Rolling Update later.

```
docker build -t hard-shutdown:v1 -t hard-shutdown:v2 ./hard-shutdown
docker build -t graceful-shutdown:v1 -t graceful-shutdown:v2 ./graceful-shutdown
```

## Deploy Redis to k8s

```
kubectl create namespace redis
helm install redis oci://registry-1.docker.io/bitnamicharts/redis --set auth.enabled=false --namespace redis
```

## Deploy APIs to k8s

```
kubectl apply -f k8s.yaml
```

## Send requests

Use [vegeta](https://github.com/tsenart/vegeta) to send 1000 requests:

hard-shutdown:
```
echo "GET http://localhost:30001/incr" | vegeta attack -duration=40s -rate=25 -http2=false | vegeta report
```

graceful-shutdown:
```
echo "GET http://localhost:30002/incr" | vegeta attack -duration=40s -rate=25 -http2=false | vegeta report
```

## Test with Kubernetes Rolling Update

While vegeta is running we can initiate a Rolling Update in Kubernetes by changing the image tag.

```
kubectl set image deployment hard-shutdown hard-shutdown=hard-shutdown:v2
kubectl set image deployment graceful-shutdown graceful-shutdown=graceful-shutdown:v2
```

## Verify counter in Redis

It should "1000".

```
kubectl exec -it redis-master-0 -n redis -- redis-cli get counter
```

## Clean up

```
kubectl delete -f k8s.yaml
helm uninstall redis --namespace redis
kubectl delete namespace redis
```