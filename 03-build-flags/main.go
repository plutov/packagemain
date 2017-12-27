package main

import (
	"fmt"
	"runtime"
)

var version string

func main() {
	fmt.Printf("OS: %s\nArchitecture: %s\n", runtime.GOOS, runtime.GOARCH)
}
