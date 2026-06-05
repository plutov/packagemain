package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/plutov/gopkg/pkgsiteapi"
)

type searchItem struct {
	Path     string
	Module   string
	Version  string
	Synopsis string
}

type searchMsg struct{ items []searchItem }

type detailMsg struct {
	path     string
	versions *pkgsiteapi.PaginatedResponse
	symbols  *pkgsiteapi.PackageSymbols
}

type model struct {
	input   textinput.Model
	vp      viewport.Model
	client  *pkgsiteapi.ClientWithResponses
	results []searchItem
	sel     int
	focus   string
	loading bool

	item     searchItem
	versions *pkgsiteapi.PaginatedResponse
	symbols  *pkgsiteapi.PackageSymbols
}

func newModel(client *pkgsiteapi.ClientWithResponses) model {
	input := textinput.New()
	input.Placeholder = "Search Go packages"
	input.Focus()
	input.Width = 50

	vp := viewport.New(80, 20)
	vp.SetContent("Loading...")

	return model{input: input, vp: vp, client: client, focus: "input"}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tea.SetWindowTitle("gopkg"))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.input.Width = msg.Width - 4
		if m.input.Width < 20 {
			m.input.Width = 20
		}
		if m.input.Width > 80 {
			m.input.Width = 80
		}
		m.vp.Width = msg.Width - 4
		if m.vp.Width < 20 {
			m.vp.Width = 20
		}
		m.vp.Height = msg.Height - 6
		if m.vp.Height < 8 {
			m.vp.Height = 8
		}
		m.refresh()
		return m, nil
	case searchMsg:
		m.loading = false
		m.results = msg.items
		m.sel = 0
		if len(m.results) > 0 {
			m.focus = "results"
			m.input.Blur()
		}
		return m, nil
	case detailMsg:
		if msg.path != m.item.Path {
			return m, nil
		}
		m.loading = false
		m.versions = msg.versions
		m.symbols = msg.symbols
		m.refresh()
		return m, nil
	}

	key, ok := msg.(tea.KeyMsg)
	if !ok {
		if m.focus == "input" {
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}
		if m.focus == "detail" {
			var cmd tea.Cmd
			m.vp, cmd = m.vp.Update(msg)
			return m, cmd
		}
		return m, nil
	}

	switch key.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "tab":
		if len(m.results) == 0 || m.focus == "detail" {
			return m, nil
		}
		if m.focus == "input" {
			m.focus = "results"
			m.input.Blur()
		} else {
			m.focus = "input"
			m.input.Focus()
		}
		return m, nil
	}

	if m.focus == "detail" {
		if key.String() == "esc" || key.String() == "backspace" || key.String() == "left" || key.String() == "h" {
			m.focus = "results"
			return m, nil
		}
		var cmd tea.Cmd
		m.vp, cmd = m.vp.Update(msg)
		return m, cmd
	}

	if m.focus == "results" {
		switch key.String() {
		case "up", "k":
			if m.sel > 0 {
				m.sel--
			}
		case "down", "j":
			if m.sel < len(m.results)-1 {
				m.sel++
			}
		case "enter", "o", "right", "l":
			if len(m.results) == 0 || m.loading {
				return m, nil
			}
			m.focus = "detail"
			m.loading = true
			m.item = m.results[m.sel]
			m.versions = nil
			m.symbols = nil
			m.vp.GotoTop()
			m.refresh()
			return m, detailCmd(m.client, m.item)
		}
		return m, nil
	}

	if key.String() == "enter" {
		q := strings.TrimSpace(m.input.Value())
		if q == "" || m.loading {
			return m, nil
		}
		m.loading = true
		return m, searchCmd(m.client, q)
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}
