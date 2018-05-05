```
dep init
```

.gitignore
```
vendor
```

Let's run tests:
```
go test ./...
```

Strange licenses fodler created, fix:
```
os.Mkdir(fullPath+"licenses", 0777)
```