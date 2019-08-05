## How to Rate Limit HTTP Requests in Go by IP


Build:

```
go build .
./17-ratelimit
```

Test:

```
vegeta attack -duration=10s -rate=100 -targets=vegeta.conf | vegeta report
```