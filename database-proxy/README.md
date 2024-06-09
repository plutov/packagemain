# mysql proxy

This proxy solves a very simple use case: intercept SQL query and rewrite table name if it matches a pattern.

```sql
-- Application-generated query
SELECT * FROM table1;

-- Rewritten query
SELECT * FROM table2;
```

## run locally with Docker Compose

```
docker-compose up --build
```

## connect to proxy

```
mysql -uroot -P 3307 -h 127.0.0.1 -proot --ssl-mode=disable
```