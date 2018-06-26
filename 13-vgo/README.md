### Introduction to Versioned Go (vgo)

If you worked with Go for a while, you know that there is no official dependency manager yet. Go developers are using tools like `dep`, `glide`, `godeps`, etc. I personally used `dep` recently. But all of them are not official one.

But recently Russ Cox has proposed to update the `go` command to work with modules. One significant change is that ordinary build commands, like go build, go install, go run, and go test, will resolve new dependencies on demand. All it takes to use a package in a brand new module is to add an import to the Go source code and build the code.

The most significant change, though, is to remove GOPATH, you can clone repository, cd into it and build!

Russ built a prototype as a standalone command `vgo`. and this video shows what it's like to use vgo.

#### Install vgo

```
go get golang.org/x/vgo
```

### Simple import

Let's create simpla main.go file with uuid package imported:

```
package main

import (
	"fmt"

	"github.com/satori/go.uuid"
)

func main() {
	u := uuid.Must(uuid.NewV4())
	fmt.Println(u.String())
}
```

And run it using `vgo`:

```
vgo run main.go
vgo: creating new go.mod: module github.com/plutov/packagemain
vgo: resolving import "github.com/satori/go.uuid"
vgo: finding github.com/satori/go.uuid (latest)
vgo: adding github.com/satori/go.uuid v1.2.0
77331d9b-2a0b-4f99-b360-d01d48bd50b2
```

As you can see it created a `go.mod` file in the parent folder, in the repository root folder, and added there `uuid` dependency. `go get` is not required anymore.

go.mod:
```
module github.com/plutov/packagemain

require github.com/satori/go.uuid v1.2.0
```

If we run `vgo run` command again, it won't print anything else:
```
vgo run main.go
a59b42cc-ec06-40a3-97b7-de17c6d7cc4b
```

And we can list our modules with `vgo list -m`:
```
vgo list -m
```