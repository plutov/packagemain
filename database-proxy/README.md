# mysql proxy

This proxy solves a very simple use case: intercept SQL query and rewrite table name if it matches a pattern.

```sql
-- Application-generated query
SELECT * FROM orders;

-- Rewritten query
SELECT * FROM orders_v1;
```

## run locally with Docker Compose

```
docker-compose up --build
```

## connect to proxy

```
mysql
```