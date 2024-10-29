package main

func (m model) View() string {
	s := "NOTES APP\n\n"

	if m.state == titleView {
		s += "Note title:\n\n"
		s += m.textinput.View() + "\n\n"
		s += "enter - save • esc - discard"
	}
	if m.state == bodyView {
		s += "Note:\n\n"
		s += "TODO\n\n"
		s += "ctrl+s - save • esc - discard"
	}
	if m.state == listView {
		s += "n - new note • q - quit"
	}

	return s
}
