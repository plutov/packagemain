package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type sessionState uint

const (
	listView sessionState = iota
	editView
)

type model struct {
	textarea textarea.Model
	state    sessionState
}

func NewModel() model {
	ta := textarea.New()
	ta.Blur()

	return model{
		textarea: ta,
		state:    listView,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "n":
			m.textarea.SetValue("")
			m.textarea.Focus()
			m.textarea.CursorEnd()
			m.state = editView
		case "q":
			return m, tea.Quit
		case "esc":
			if m.state == editView {
				m.textarea.Blur()
			}
		}
	}

	return m, tea.Batch(cmds...)
}
