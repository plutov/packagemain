package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	listView uint = iota
	titleView
	bodyView
)

type model struct {
	state     uint
	textinput textinput.Model
}

func NewModel() model {
	ti := textinput.New()
	ti.Blur()

	return model{
		state:     listView,
		textinput: ti,
	}
}

func (m model) Init() tea.Cmd {
	// TODO: can add some initial I/O
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	var cmd tea.Cmd
	m.textinput, cmd = m.textinput.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch m.state {
		// List View key bindings
		case listView:
			switch key {
			case "q":
				return m, tea.Quit
			case "n":
				m.textinput.SetValue("")
				m.textinput.Focus()
				m.state = titleView
			}

		// Title Input View key bindings
		case titleView:
			switch key {
			case "enter":
				// TODO: save
			case "esc":
				m.textinput.Blur()
				m.state = listView
			}

		// Body Textarea key bindings
		case bodyView:
			switch key {
			case "ctrl+s":
				// TODO: save
				m.state = listView
			case "esc":
				m.state = listView
			}
		}
	}

	return m, tea.Batch(cmds...)
}
