package survey

import (
	"encoding/json"
	"io"
)

type Survey struct {
	Title     string
	Questions []string
}

func (s *Survey) GetTitle() string {
	return s.Title
}

func (s *Survey) Validate() bool {
	return len(s.Questions) > 0
}

// Repository is an interface that defines the methods that a concrete repository must implement
type Repository interface {
	Save(survey *Survey) error
}

type InMemoryRepository struct {
	surveys []*Survey
}

func (r *InMemoryRepository) Save(survey *Survey) error {
	r.surveys = append(r.surveys, survey)
	return nil
}

// Decouple the SaveSurvey function from the concrete Repository implementation
func SaveSurvey(survey *Survey, repo Repository) error {
	return repo.Save(survey)
}

type Exporter interface {
	Export(survey *Survey) error
}

type S3Exporter struct{}

func (e *S3Exporter) Export(survey *Survey) error {
	return nil
}

type GCSExporter struct{}

func (e *GCSExporter) Export(survey *Survey) error {
	return nil
}

func (s *Survey) Export(exporter Exporter) error {
	return exporter.Export(s)
}

// Use generic io.Writer instead of a concrete struct/file/etc.
func (s *Survey) Write(writer io.Writer) (int, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return 0, err
	}

	return writer.Write(b)
}

// Base interface
type Question interface {
	SetTitle()
}

// Has addtional method AddOption
type QuestionWithOptions interface {
	Question
	AddOption()
}

type TextQuestion struct {
	Title string
}

func (q *TextQuestion) SetTitle(title string) {
	q.Title = title
}

type DropdownQuestion struct {
	Title   string
	Options []string
}

func (q *DropdownQuestion) SetTitle(title string) {
	q.Title = title
}

func (q *DropdownQuestion) AddOption(option string) {
	q.Options = append(q.Options, option)
}

type SurveyManager struct {
	store Repository
}

// Pass the Repository interface instead of the concrete InMemoryRepository
func NewSurveyManager(store Repository) *SurveyManager {
	return &SurveyManager{
		store: store,
	}
}

func (m *SurveyManager) Save(survey *Survey) error {
	return m.store.Save(survey)
}
