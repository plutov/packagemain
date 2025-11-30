package main

/*
#include <jemalloc/jemalloc.h>
*/
import "C"

import (
	"crypto/rand"
	"time"
	"unsafe"
)

func main() {
	go StartRSSMemoryMonitor()

	// simulateGC()
	simulateJemalloc()
}

func simulateGC() {
	for range 50 {
		size := 20 * 1024 * 1024
		data := make([]byte, size)

		rand.Read(data)
		time.Sleep(200 * time.Millisecond)

		data = nil
	}

	time.Sleep(2 * time.Second)
}

func simulateJemalloc() {
	for range 50 {

		size := 20 * 1024 * 1024

		ptr := C.malloc(C.size_t(size))
		if ptr == nil {
			panic("unable to allocate memory")
		}

		data := unsafe.Slice((*byte)(ptr), size)

		rand.Read(data)

		time.Sleep(500 * time.Millisecond)

		C.free(ptr)
		data = nil

		time.Sleep(500 * time.Millisecond)
	}

	time.Sleep(1 * time.Second)
}
