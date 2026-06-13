package main

import (
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/plutov/gopkg/pkgsiteapi"
)

const apiBaseUrl = "https://pkg.go.dev/v1beta"

func main() {
	client, err := pkgsiteapi.NewClientWithResponses(apiBaseUrl)
	if err != nil {
		slog.Error("Unable to create an API client", "err", err)
		os.Exit(1)
	}

	_, err = tea.NewProgram(newModel(client), tea.WithAltScreen()).Run()
	if err != nil {
		slog.Error("Unable to run Tea program", "err", err)
		os.Exit(1)
	}
}
