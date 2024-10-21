package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

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

    ping := time.NewTicker(time.Second)
    defer ping.Stop()

    for {
        select {
        case <-clientGone:
            fmt.Println("client disconnected")
            return
        case <-ping.C:
            // Send an event to the client
            now := time.Now().UTC().Unix()
            eventType := "ping"
            if now % 2 == 0 {
            	eventType = "data"
            }
            if _, err := fmt.Fprintf(w, "event:%s\nid:%d\ndata:the time is %d\n\n", eventType, now, now); err != nil {
                return
            }

            if err := rc.Flush(); err != nil {
                return
            }
        }
    }
}

func main() {
    http.HandleFunc("/events", sseHandler)
    fmt.Println("server is running on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("unable to start the server: %s", err.Error())
    }
}
