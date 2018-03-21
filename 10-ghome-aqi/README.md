app.yaml
```
runtime: go
api_version: go1

handlers:
- url: /
  script: _go_app
```

```
gcloud app deploy
```