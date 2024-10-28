package main

import (
	"github.com/charmbracelet/lipgloss"
)

var appNameStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FF5F87")).
	Bold(true)

var errorAlertStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFFDF5")).
	Background(lipgloss.Color("#FF5F87")).
	Padding(0, 1)

var statusStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFFDF5")).
	Background(lipgloss.Color("#009933")).
	Padding(0, 1)

var errorInfoStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("250")).
	Padding(0, 1)

var listEnumeratorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).MarginRight(1)
var listItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).MarginRight(1)

var faint = lipgloss.NewStyle().Foreground(lipgloss.Color("250")).Faint(true)

func (m model) View() string {
	s := ""
	if m.state == editView {
		s += m.textarea.View() + "\n\n"
		s += "ctrl+s - save • esc - discard • q - quit\n"
	}

	if m.state == listView {
		s += "n - new note • q - quit\n"
	}

	return s
}
