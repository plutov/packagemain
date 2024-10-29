package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	listView uint = iota
	titleView
	bodyView
)

type model struct {
	store     Store
	state     uint
	textarea  textarea.Model
	textinput textinput.Model
	currNote  Note
	notes     []Note
	listIndex int
}

func NewModel(store Store) model {
	ta := textarea.New()
	ta.Blur()

	ti := textinput.New()
	ti.Blur()

	notes, _ := store.GetNotes()

	return model{
		store:     store,
		state:     listView,
		textarea:  ta,
		textinput: ti,
		notes:     notes,
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
				m.currNote = Note{}
				m.state = titleView
			case "up", "k":
				if m.listIndex > 0 {
					m.listIndex--
				}
			case "down", "j":
				if m.listIndex < len(m.notes)-1 {
					m.listIndex++
				}
			case "enter":
				m.currNote = m.notes[m.listIndex]
				m.state = bodyView
				m.textarea.SetValue(m.currNote.Body)
				m.textarea.Focus()
				m.textarea.CursorEnd()
			}

		// Title Input View key bindings
		case titleView:
			switch key {
			case "enter":
				if m.textinput.Value() != "" {
					m.currNote.Title = m.textinput.Value()
					m.textinput.Blur()

					m.state = bodyView
					m.textarea.SetValue("")
					m.textarea.Focus()
					m.textarea.CursorEnd()
				}
			case "esc":
				m.textinput.Blur()
				m.state = listView
			}

		// Body Textarea key bindings
		case bodyView:
			switch key {
			case "ctrl+s":
				m.currNote.Body = m.textarea.Value()
				m.textarea.Blur()

				m.store.SaveNote(m.currNote)
				var err error
				m.notes, err = m.store.GetNotes()
				if err != nil {
					return m, tea.Quit
				}

				m.state = listView
			case "esc":
				m.textarea.Blur()
				m.state = listView
			}
		}
	}

	return m, tea.Batch(cmds...)
}
