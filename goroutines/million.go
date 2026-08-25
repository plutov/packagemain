package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	const n = 1_000

	for range n {
		go func() {
			select {}
		}()
	}

	time.Sleep(time.Second)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("goroutines: %d\n", runtime.NumGoroutine())
	fmt.Printf("heap:       %.2f MB\n", float64(m.HeapAlloc)/1024/1024)
	fmt.Printf("sys memory: %.2f MB\n", float64(m.Sys)/1024/1024)

	select {}
}
