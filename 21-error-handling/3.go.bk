package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func processFile(path string) error {
	if path == "" {
		return errors.New("path is empty")
	}

	f, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "unable to open a file")
	}

	if f.Name() != "test.txt" {
		return fmt.Errorf("invalid filename: %s", f.Name())
	}

	return nil
}

func main() {
	err := processFile("./foo")

	if err != nil {
		switch errors.Cause(err).(type) {
		case *os.PathError:
			fmt.Println("file doesn't exist")
		default:
			fmt.Printf("unknown error: %+v", err)
		}

		return
	}

	fmt.Println("ok")
}
