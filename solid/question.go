package survey

type Question interface {
	SetTitle(title string)
}

type QuestionWithOptions interface {
	Question
	AddOption(option string)
}

type TextInputQuestion struct {
	title string
}

func (q *TextInputQuestion) SetTitle(title string) {
	q.title = title
}

type DropdownQuestion struct {
	title   string
	options []string
}

func (q *DropdownQuestion) SetTitle(title string) {
	q.title = title
}

func (q *DropdownQuestion) AddOption(option string) {
	q.options = append(q.options, option)
}
