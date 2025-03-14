Start mysql:

```
docker-compose up
```

Start dbproxy:

```
go run .
```

Connect to proxy:

```
mysql -h 127.0.0.1 -P 3307 -uroot -ppass --ssl-mode=disabled db
```
