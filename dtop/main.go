package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

//go:embed index.html
var assets embed.FS

var page = template.Must(template.New("index.html").Funcs(template.FuncMap{"bar": terminalBar}).ParseFS(assets, "index.html"))

type App struct{}

func (a App) home(w http.ResponseWriter, r *http.Request) {
	m := collectMetrics()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := page.ExecuteTemplate(w, "index.html", m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a App) events(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			m := collectMetrics()
			fmt.Fprintf(w, "event: datastar-patch-elements\ndata: elements %s\ndata: mode outer\n\n", metricsFragment(m))
			flusher.Flush()
		}
	}
}

func main() {
	a := App{}
	mux := http.NewServeMux()
	mux.HandleFunc("/", a.home)
	mux.HandleFunc("/events", a.events)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
