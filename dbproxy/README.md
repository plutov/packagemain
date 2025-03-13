Start postgres:

```
docker-compose up
```

Start dbproxy:

```
go run .
```

Connect to proxy:

```
psql -h 127.0.0.1 -p 55432 -U user -d
```
