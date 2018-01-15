### How to use "go build" flags

Go is very famous for it's fast and flexible compilation. It's almost impossible to write a program that will take long time to compile. One more cool thing is cross compilation in Go, from any platform you can easily compile an executable for another platform. Go now comes with support for all architectures built in. That means cross compiling is as simple as setting the right `GOOS` and `GOARCH` environment variables before building.

Now let's test this out. Let’s compile this program for an Apple MacBook. To do so, we simply set two environment variables: $GOOS, which is the target operating system, and $GOARCH, which is the target processor. Then we run go build as normal:

```
GOOS=darwin GOARCH=386 go build main.go
```

Note that the main executable only runs on OS X, and cannot be run on Linux or Windows. On the other hand, if we wanted to compile for Microsoft Windows, we’d simply set GOOS=windows and GOARCH=386.

```
GOOS=windows GOARCH=386 go build main.go
```

The result is .exe file.

Now let's see what flags do we have in "go build" command. If you are curious about the Go toolchain, or using a cross-C compiler and wondering about flags passed to the external compiler, or suspicious about a linker bug, use -x to see all the invocations.

```
go build -x main.go
```

I often use -ldflags option when build Go programs. This will pass arguments on each go tool link invocation. It can be useful to set version of your app during the build and not save it in the code.

```
var version string
```

```
go build -ldflags "-X main.version=v0.0.1" main.go
```

-X Sets the value of the string variable in package var. Now during the build we can pass git hash there for example.

Go includes debugging information into binary for GDB. It makes executable file a bit bigger. If you distribute the binary to customers, and they are not supposed to debug it, you can turn it off with flags -w -s.

Let's compare the size. I will compile a big application I have, which contains a lof of dependencies, so we can see the real difference.

```
cd logpacker
go build -o daemon_debug daemon/daemon.go
go build -ldflags="-w -s" -o daemon_no_debug daemon/daemon.go
du -h daemon_debug
du -h daemon_no_debug
strings daemon_debug | wc -l
strings daemon_no_debug | wc -l
```

File size is 3.6MB smaller and has less information, however it works the same.

go build -gcflags used to pass flags to the Go compiler. go tool compile -help lists all the flags that can be passed to the compiler.

```
go tool compile -help
```

Also, Go by default doing some optimizations with your code, for example inlining. When a package is compiled, any small function that is suitable for inlining is marked and then compiled as usual.

```
func getVersion() string {
	return version
}
```

To see what Go does with inlining you can use -m compiler flag:

```
go build -gcflags=-m main.go
```

That's it, enjoy compiling your programs and see you later!

Resources:
https://dave.cheney.net/2012/10/07/notes-on-exploring-the-compiler-flags-in-the-go-compiler-suite
http://pliutau.com/go-tools-are-awesome/
http://pliutau.com/optimize-go-binary-size/
https://dave.cheney.net/2015/08/22/cross-compilation-with-go-1-5