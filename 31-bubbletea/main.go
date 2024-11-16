package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	store := new(Store)
	if err := store.Init(); err != nil {
		log.Fatalf("unable to init store: %v", err)
	}

	m := NewModel(store)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatalf("unable to run tui: %v", err)
	}
}
