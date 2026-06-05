package main

import (
	"strings"
	"time"
)

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
