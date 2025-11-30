package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

func StartRSSMemoryMonitor() {
	pid := int32(os.Getpid())
	p, err := process.NewProcess(pid)
	if err != nil {
		panic(err)
	}

	fmt.Println("Memory:")

	for {
		memInfo, _ := p.MemoryInfo()
		rssMB := int(memInfo.RSS / 1024 / 1024)

		barLength := 100
		filled := min(rssMB, barLength)

		bar := strings.Repeat("#", filled) + strings.Repeat(" ", barLength-filled)

		fmt.Printf("[%s] %3dMB\n", bar, rssMB)

		time.Sleep(100 * time.Millisecond)
	}
}
