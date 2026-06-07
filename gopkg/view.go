package main

import "strings"

func (m Model) View() string {
	lines := []string{"gopkg — pkg.go.dev explorer"}
	if m.errMsg != "" {
		lines = append(lines, "Error: "+m.errMsg, "")
	} else {
		lines = append(lines, "")
	}

	if m.focus == focusDetail {
		lines = append(lines, "", m.viewport.View(), "", "esc: back • j/k: scroll • ctrl+c: quit")
		return strings.Join(lines, "\n")
	}

	lines = append(lines, m.input.View(), "")

	if m.loading {
		lines = append(lines, "Loading...")
	} else if len(m.results) > 0 {
		for i, it := range m.results {
			line := "  " + it.PackagePath
			if i == m.currItemIndex {
				line = "> " + it.PackagePath
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
	lines = append(lines, "tab: switch • enter: search/open • ctrl+c: quit")
	return strings.Join(lines, "\n")
}

func (m *Model) refreshViewport() {
	if m.focus != focusDetail || m.currItem == nil {
		return
	}
	if m.loading {
		m.viewport.SetContent(m.currItem.PackagePath + "\n\nLoading...")
		return
	}

	lines := []string{m.currItem.PackagePath}
	if m.currItem.Synopsis != "" {
		lines = append(lines, "", m.currItem.Synopsis)
	}

	versions := m.versions.VersionResults()
	lines = append(lines, "", "Versions")
	for _, v := range versions {
		line := "- " + v.Version
		if v.CommitTime != nil {
			line += "  " + v.CommitTime.Format("2006-01-02")
		}
		lines = append(lines, line)
	}

	symbols := m.symbols.SymbolResults()
	lines = append(lines, "", "Symbols")
	for _, s := range symbols {
		lines = append(lines, "- "+s.Kind+" "+s.Name)
		if s.Synopsis != "" {
			lines = append(lines, "  "+s.Synopsis)
		}
	}

	m.viewport.SetContent(strings.Join(lines, "\n"))
}
