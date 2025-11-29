package main

import (
	"crypto/rand"
	"time"
)

func main() {
	StartRSSMemoryMonitor()

	for range 50 {
		size := 20 * 1024 * 1024
		data := make([]byte, size)

		rand.Read(data)
		time.Sleep(200 * time.Millisecond)

		data = nil
	}

	time.Sleep(1 * time.Second)
}
