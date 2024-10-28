package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

func GetTextareaModel() textarea.Model {
	ti := textarea.New()
	tiFocusedStyle := textarea.Style{
		Base:             lipgloss.NewStyle(),
		CursorLine:       lipgloss.NewStyle().Background(lipgloss.Color("0")),
		CursorLineNumber: lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		EndOfBuffer:      lipgloss.NewStyle().Foreground(lipgloss.Color("0")),
		LineNumber:       lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
		Placeholder:      lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		Prompt:           lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
		Text:             lipgloss.NewStyle(),
	}
	tiBlurredStyle := textarea.Style{
		Base:             lipgloss.NewStyle(),
		CursorLine:       lipgloss.NewStyle(),
		CursorLineNumber: lipgloss.NewStyle(),
		EndOfBuffer:      lipgloss.NewStyle(),
		LineNumber:       lipgloss.NewStyle(),
		Placeholder:      lipgloss.NewStyle(),
		Prompt:           lipgloss.NewStyle(),
		Text:             lipgloss.NewStyle(),
	}
	ti.FocusedStyle = tiFocusedStyle
	ti.BlurredStyle = tiBlurredStyle
	ti.Blur()

	return ti
}

func GetInputModel() textinput.Model {
	ti := textinput.New()
	ti.CharLimit = 3
	ti.Width = 20
	ti.Blur()

	return ti
}
