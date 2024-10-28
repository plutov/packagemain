package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	s := new(InMemStore)
	m := NewModel(s)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("unable to run: %v", err)
		os.Exit(1)
	}
}
