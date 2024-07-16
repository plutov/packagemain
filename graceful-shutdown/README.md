[Read the full article on packagemain.tech]()

## Build applications

```
docker build -t hard-shutdown ./hard-shutdown
docker build -t graceful-shutdown ./graceful-shutdown
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

## Test with Kubernetes Termination

While vegeta is running init the Kubernetes Termination process.

```
kubectl delete pods --all
```

## Verify counter in Redis

```
kubectl exec -it redis-master-0 -n redis -- redis-cli
127.0.0.1:6379> get counter
"1000"
```

## Clean up

```
kubectl delete -f k8s.yaml
helm uninstall redis --namespace redis
kubectl delete namespace redis
```