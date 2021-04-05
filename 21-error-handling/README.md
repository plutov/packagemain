Today we're going to talk about error handling in Go and some of its best practices.

## packagemain #21: Error Handling in Go

Proper handling errors is an essential element of a solid code. Error handling in Go has been the topic for hot discussions a lot due to its unconventional approach, when an error is just a value that a function can return if something unexpected happened, compared to `try...catch` block.

I personally like this approach, it reminds me that I should check every single place where an error may occur.

In Go there is a built-in type `error` with zero value as `nil`, so in order to handle an error we can check if the value is not `nil`. The code could be therefore a bit verbose, and `errors` package doesn't provide much functionality... So let's review some of the best practices and the `errors` package(s).

First of al let's review the built-in `errors` package.

There are 2 functions to create an error variable with a specified text message: `errors.New()` and `fmt.Errorf()`.

It's also possible to check what kind of error was returned, by checking the value of the error or by checking its type:

1. The first one is called `Sentinel Errors`, there are multiple examples in Go standard library, for example `os.ErrNotExist`  or `io.EOF`. These error variables have to be exported with your package, so the caller can check them.

2. There is another way to check the error, is to check its type. Callers can use type assertion to check the type.

```go
switch err.(type) {
case *os.PathError:
    fmt.Println("file doesn't exist")
default:
    fmt.Println("unknown error")
}
```

The `errors` package has its own limitations, let's review the case when we need to ad some context before returning an error.

```go
f, err := os.Open(path)
if err != nil {
    return fmt.Errorf("unable to open a file: %s", err.Error())
}

// ...

switch err.(type) {
// this won't work anymore, since we dropped the type
case *os.PathError:
    fmt.Println("file doesn't exist")
```

The stack trace was dropped as well. So what can we do to keep the stack and the type?

There are several packages to resolve the problem. They can replace `errors` package since they are built on top of it, so you won't need to do any changes in case you used the standard `errors` package.

### github.com/pkg/errors or golang.org/x/xerrors

These packages add few key methods. `Wrap()` - used to wrap the underlying error, add contextual text information, and attach the call stack. `WithMessage()` is used to add contextual text information to underlying error without attaching call stack. And `Cause()` method is for determining the underlying error.

Using these packages we can now keep the stack trace and original type:

```go
f, err := os.Open(path)
if err != nil {
    return errors.Wrap(err, "unable to open a file")
}

// ...

switch errors.Cause(err).(type) {
case *os.PathError:
    fmt.Println("file doesn't exist")
default:
    fmt.Printf("unknown error: %+v", err)
}
```

### Need something more sophisticated?

Working on distributed systems, there may be additional needs for your error handling mechanism. I enjoyed using `cocroachdb/errors` package which is also compatible with `pkg/errors` or `errors` and adds more features, for example: stripping PII details, supporting Sentry.io format, separating system message vs hint that can be shown to the user.

### Conclusion

Overall I enjoy the error handling in Go, even when it's verbose it's quite easy to  understand and straightforward. Do you use any best practices in your teams or in your projects? Do share them in the comments below!