package main

import (
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/plutov/gopkg/pkgsiteapi"
)

const apiBaseURL = "https://pkg.go.dev/v1beta"

func main() {
	client, err := pkgsiteapi.NewClientWithResponses(apiBaseURL)
	if err != nil {
		slog.Error("Error creating API client", "error", err)
		os.Exit(1)
	}

	if _, err := tea.NewProgram(newModel(client), tea.WithAltScreen()).Run(); err != nil {
		slog.Error("Error running program", "error", err)
		os.Exit(1)
	}
}
