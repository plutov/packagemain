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

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	w.Header().Set("Access-Control-Allow-Origin", "*")

	memT := time.NewTicker(time.Second)
	defer memT.Stop()

	cpuT := time.NewTicker(time.Second)
	defer cpuT.Stop()

	clientGone := r.Context().Done()

	rc := http.NewResponseController(w)

	for {
		select {
		case <-clientGone:
			fmt.Println("client has disconnected")
			return
		case <-memT.C:
			m, err := mem.VirtualMemory()
			if err != nil {
				log.Printf("unable to get mem: %s", err.Error())
				return
			}

			if _, err := fmt.Fprintf(w, "event:mem\ndata:Total: %d, Used: %d, Perc: %.2f%%\n\n", m.Total, m.Used, m.UsedPercent); err != nil {
				log.Printf("unable to write: %s", err.Error())
				return
			}

			rc.Flush()
		case <-cpuT.C:
			c, err := cpu.Times(false)
			if err != nil {
				log.Printf("unable to get cpu: %s", err.Error())
				return
			}

			if _, err := fmt.Fprintf(w, "event:cpu\ndata:User: %.2f, Sys: %.2f, Idle: %.2f\n\n", c[0].User, c[0].System, c[0].Idle); err != nil {
				log.Printf("unable to write: %s", err.Error())
				return
			}

			rc.Flush()
		}
	}
}
