package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

func main() {
	http.HandleFunc("/events", sseHandler)
	fmt.Println("server is running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("unable to start the server: %s", err.Error())
	}
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	// Set http headers required for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// You may need this locally for CORS requests
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create a channel for client disconnection
	clientGone := r.Context().Done()

	rc := http.NewResponseController(w)

	memTicker := time.NewTicker(time.Second)
	defer memTicker.Stop()

	cpuTicker := time.NewTicker(time.Second)
	defer cpuTicker.Stop()

	for {
		select {
		case <-clientGone:
			fmt.Println("client disconnected")
			return
		case <-memTicker.C:
			v, err := mem.VirtualMemory()
			if err != nil {
				return
			}

			// Send an event to the client
			if _, err := fmt.Fprintf(w, "event:mem\ndata:%s\n\n", fmt.Sprintf("Total: %v, Free: %v, Used: %.2f%%\n", v.Total, v.Free, v.UsedPercent)); err != nil {
				return
			}

			if err := rc.Flush(); err != nil {
				return
			}
		case <-cpuTicker.C:
			c, err := cpu.Times(false)
			if err != nil || len(c) == 0 {
				return
			}

			// Send an event to the client
			if _, err := fmt.Fprintf(w, "event:cpu\ndata:%s\n\n", fmt.Sprintf("User: %.2f, Sys: %.2f, Idle: %.2f\n", c[0].User, c[0].System, c[0].Idle)); err != nil {
				return
			}

			if err := rc.Flush(); err != nil {
				return
			}
		}
	}
}
