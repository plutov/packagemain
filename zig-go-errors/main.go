package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	f, err := os.Open("main.zig")
	if err != nil {
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("name: %s\n", f.Name())
}
