package main

import (
	"github.com/charmbracelet/lipgloss"
)

var appNameStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#000000")).
	Background(lipgloss.Color("86")).
	Padding(0, 1)

var faint = lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Faint(true)

func (m model) View() string {
	s := appNameStyle.Render("Notes App") + "\n\n"

	if m.state == editView {
		s += m.textarea.View() + "\n\n"
		s += faint.Render("ctrl+s - save • esc - discard • q - quit")
	}

	if m.state == listView {
		s += faint.Render("n - new note • q - quit")
	}

	return s
}
