package main

/*
#include <jemalloc/jemalloc.h>
*/
import "C"

import (
	"crypto/rand"
	"flag"
	"time"
	"unsafe"
)

var gc = flag.Bool("gc", false, "Use Go's GC for memory allocation instead of jemalloc")

func main() {
	flag.Parse()

	go StartRSSMemoryMonitor()

	if *gc {
		simulateGC()
	} else {
		simulateJemalloc()
	}
}

func simulateGC() {
	for range 50 {
		size := 20 * 1024 * 1024
		data := make([]byte, size)

		rand.Read(data)

		time.Sleep(500 * time.Millisecond)
		data = nil
		time.Sleep(500 * time.Millisecond)
	}

	time.Sleep(2 * time.Second)
}

func simulateJemalloc() {
	for range 50 {
		size := 20 * 1024 * 1024

		ptr := C.malloc(C.size_t(size))
		if ptr == nil {
			panic("no memory")
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
