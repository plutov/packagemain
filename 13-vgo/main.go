package main

import (
	"fmt"

	"github.com/satori/go.uuid"
)

func main() {
	fmt.Println(uuid.NewV4().String())
}
