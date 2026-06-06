package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/plutov/gopkg/pkgsiteapi"
)

type focusMode int

const (
	focusInput focusMode = iota
	focusResults
	focusDetail
)

type Model struct {
	client *pkgsiteapi.ClientWithResponses

	input    textinput.Model
	viewport viewport.Model

	results  []pkgsiteapi.SearchResult
	currItem *pkgsiteapi.SearchResult
	versions *pkgsiteapi.PaginatedResponse
	symbols  *pkgsiteapi.PackageSymbols
	errMsg   string

	loading       bool
	currItemIndex int
	focus         focusMode
}

type searchMsg struct {
	items []pkgsiteapi.SearchResult
}

type detailMsg struct {
	path     string
	versions *pkgsiteapi.PaginatedResponse
	symbols  *pkgsiteapi.PackageSymbols
}

type errorMsg struct {
	text string
}

func newModel(client *pkgsiteapi.ClientWithResponses) Model {
	input := textinput.New()
	input.Placeholder = "Search Go packages"
	input.Focus()
	input.Width = 50

	vp := viewport.New(80, 20)
	vp.SetContent("Loading...")

	return Model{input: input, viewport: vp, client: client, focus: focusInput}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.refreshViewport()
		return m, nil
	case searchMsg:
		m.loading = false
		m.errMsg = ""
		m.results = msg.items
		m.currItemIndex = 0
		if len(m.results) > 0 {
			m.focus = focusResults
			m.input.Blur()
		}
		m.refreshViewport()
		return m, nil
	case detailMsg:
		m.loading = false
		m.errMsg = ""
		m.versions = msg.versions
		m.symbols = msg.symbols
		m.refreshViewport()
		return m, nil
	case errorMsg:
		m.loading = false
		m.errMsg = msg.text
		m.refreshViewport()
		return m, nil
	}

	key, ok := msg.(tea.KeyMsg)
	if !ok {
		if m.focus == focusInput {
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}

		if m.focus == focusDetail {
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}
		return m, nil
	}

	switch key.String() {
	case "q":
		return m, tea.Quit
	case "tab":
		if m.focus == focusInput {
			m.focus = focusResults
			m.input.Blur()
		} else {
			m.focus = focusInput
			m.input.Focus()
		}
		return m, nil
	}

	if m.focus == focusDetail {
		if key.String() == "esc" {
			m.focus = focusResults
			return m, nil
		}
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	if m.focus == focusResults {
		switch key.String() {
		case "k":
			if m.currItemIndex > 0 {
				m.currItemIndex--
			}
		case "j":
			if m.currItemIndex < len(m.results)-1 {
				m.currItemIndex++
			}
		case "enter":
			m.focus = focusDetail
			m.loading = true
			m.errMsg = ""
			m.currItem = &m.results[m.currItemIndex]
			m.versions = nil
			m.symbols = nil
			m.viewport.GotoTop()
			m.refreshViewport()
			return m, detailCmd(m.client, m.currItem)
		}
		return m, nil
	}

	if key.String() == "enter" {
		q := strings.TrimSpace(m.input.Value())
		if m.loading {
			return m, nil
		}
		m.loading = true
		m.errMsg = ""
		return m, searchCmd(m.client, q)
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}
