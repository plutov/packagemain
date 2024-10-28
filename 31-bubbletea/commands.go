package main

import (
	"strings"

	"github.com/plutov/ultrafocus/hosts"
)

type command struct {
	Run  func(m model) model
	Name string
	Desc string
}

var commandFocusOn = command{
	Name: "focus on",
	Desc: "Start focus window.",
	Run: func(m model) model {
		if err := hosts.WriteDomainsToHostsFile(m.domains, hosts.FocusStatusOn); err != nil {
			m.fatalErr = err
			return m
		}

		m.status = hosts.FocusStatusOn
		return m
	},
}

var commandFocusOff = command{
	Name: "focus off",
	Desc: "Stop focus window.",
	Run: func(m model) model {
		if err := hosts.WriteDomainsToHostsFile(m.domains, hosts.FocusStatusOff); err != nil {
			m.fatalErr = err
			return m
		}

		m.status = hosts.FocusStatusOff
		return m
	},
}

var commandConfigureBlacklist = command{
	Name: "blacklist",
	Desc: "Configure blacklist.",
	Run: func(m model) model {
		m.state = blacklistView
		m.textarea.SetValue(strings.Join(m.domains, "\n"))
		m.textarea.Focus()
		m.textarea.CursorEnd()
		return m
	},
}
