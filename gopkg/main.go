package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/plutov/gopkg/pkgsiteapi"
)

const apiBaseURL = "https://pkg.go.dev/v1beta"

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
	focus   string // input, results, detail
	loading bool

	item     searchItem
	versions *pkgsiteapi.PaginatedResponse
	symbols  *pkgsiteapi.PackageSymbols
}

func newModel(client *pkgsiteapi.ClientWithResponses) model {
	in := textinput.New()
	in.Placeholder = "Search Go packages"
	in.Focus()
	in.Width = 50

	vp := viewport.New(80, 20)
	vp.SetContent("Loading...")

	return model{input: in, vp: vp, client: client, focus: "input"}
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
	}

	if m.focus == "detail" {
		switch key.String() {
		case "esc", "backspace", "left", "h":
			m.focus = "results"
			return m, nil
		}
		var cmd tea.Cmd
		m.vp, cmd = m.vp.Update(msg)
		return m, cmd
	}

	switch key.String() {
	case "tab":
		if len(m.results) == 0 {
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

func (m model) View() string {
	if m.focus == "detail" {
		return strings.Join([]string{
			"gopkg — package details",
			m.item.Path,
			"",
			m.vp.View(),
			"",
			"esc: back • j/k: scroll • q: quit",
		}, "\n")
	}

	lines := []string{"gopkg — pkg.go.dev explorer", "", m.input.View(), ""}
	if m.loading {
		lines = append(lines, "Loading...")
	} else if len(m.results) == 0 {
		lines = append(lines, "Type query and press Enter.")
	} else {
		for i, it := range m.results {
			line := "  " + it.Path
			if i == m.sel {
				line = "> " + it.Path
			}
			if it.Version != "" {
				line += "  " + it.Version
			}
			lines = append(lines, line)
			if it.Synopsis != "" {
				lines = append(lines, "   "+it.Synopsis)
			}
			lines = append(lines, "")
		}
	}
	lines = append(lines, "tab: switch • enter: search/open • q: quit")
	return strings.Join(lines, "\n")
}

func (m *model) refresh() {
	if m.focus != "detail" || m.vp.Width == 0 || m.vp.Height == 0 {
		return
	}
	if m.loading {
		m.vp.SetContent(m.item.Path + "\n\nLoading...")
		return
	}

	lines := []string{m.item.Path}
	if m.item.Module != "" {
		lines = append(lines, m.item.Module)
	}
	if m.item.Synopsis != "" {
		lines = append(lines, "", m.item.Synopsis)
	}

	lines = append(lines, "", "Versions")
	for _, v := range *m.versions.Items {
		version, _ := v["version"].(string)
		commit, _ := v["commitTime"].(string)
		if t, err := time.Parse(time.RFC3339, commit); err == nil {
			commit = t.Format("2006-01-02")
		}
		line := "- " + version
		if commit != "" {
			line += "  " + commit
		}
		lines = append(lines, line)
	}

	lines = append(lines, "", "Symbols")
	for _, s := range *m.symbols.Symbols.Items {
		kind, _ := s["kind"].(string)
		name, _ := s["name"].(string)
		lines = append(lines, "- "+kind+" "+name)
		if synopsis, _ := s["synopsis"].(string); synopsis != "" {
			lines = append(lines, "  "+synopsis)
		}
	}

	m.vp.SetContent(strings.Join(lines, "\n"))
}

func searchCmd(client *pkgsiteapi.ClientWithResponses, q string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		limit := 10
		resp, err := client.GetSearchWithResponse(ctx, &pkgsiteapi.GetSearchParams{Q: &q, Limit: &limit})
		if err != nil {
			panic(err)
		}
		items := make([]searchItem, 0, len(*resp.JSON200.Items))
		for _, it := range *resp.JSON200.Items {
			path, _ := it["packagePath"].(string)
			module, _ := it["modulePath"].(string)
			version, _ := it["version"].(string)
			synopsis, _ := it["synopsis"].(string)
			items = append(items, searchItem{Path: path, Module: module, Version: version, Synopsis: synopsis})
		}
		return searchMsg{items}
	}
}

func detailCmd(client *pkgsiteapi.ClientWithResponses, item searchItem) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		path := item.Path
		if item.Module != "" {
			path = item.Module
		}

		limit := 12
		versions, err := client.GetVersionsWithResponse(ctx, path, &pkgsiteapi.GetVersionsParams{Limit: &limit})
		if err != nil {
			panic(err)
		}

		var module, version *string
		if item.Module != "" {
			module = &item.Module
		}
		if item.Version != "" {
			version = &item.Version
		}

		limit = 20
		symbols, err := client.GetSymbolsWithResponse(ctx, item.Path, &pkgsiteapi.GetSymbolsParams{
			Module:  module,
			Version: version,
			Limit:   &limit,
		})
		if err != nil {
			panic(err)
		}

		return detailMsg{path: item.Path, versions: versions.JSON200, symbols: symbols.JSON200}
	}
}

func main() {
	client, err := pkgsiteapi.NewClientWithResponses(apiBaseURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if _, err := tea.NewProgram(newModel(client), tea.WithAltScreen()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
