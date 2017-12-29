### Using go-bindata with with html/template

#### What is go-bindata and why do we need it?

[go-bindata](https://github.com/jteeuwen/go-bindata) converts any text or binary file into Go source code, which is useful for embedding data into Go programs. So you can build your whole project into 1 binary file for easier delivery.

```
go get -u github.com/jteeuwen/go-bindata/...
```

#### html/template

[html/template](https://golang.org/pkg/html/template/)'s functions `Parse`, `ParseFiles` works only with files on the filesystem, so we need to implement a port to work with both approaches: files or go-bindata. Files:

```
go build && ./go-bindata-tpl

<!DOCTYPE html>
<html lang="en">
<body>
Hello
</body>
</html>
```

#### Generate templates with go-bindata

We need to install go-bindata CLI and generate a .go file from our templates:

```
go get -u github.com/jteeuwen/go-bindata/...
go-bindata -o tpl.go tpl
```

I prefer to add last command to `go:generate`:

```
//go:generate go-bindata -o tpl.go tpl
```

And use:
```
go generate
```

#### Use go-bindata templates

I made it by providing a flag `-go-bindata`:

```
./go-bindata-tpl -go-bindata
<!DOCTYPE html>
<html lang="en">
<body>
Hello
</body>
</html>
```

#### Conclusion

 - With `go-bindata` you can simplify your deployment with only one binary file.
 - `go-bindata` can give you a little faster templates reading.
 - Note that if you use `ParseFiles` you have to change it to work with `Assert` function.