package survey

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

// ===========

type Repository interface {
	Save(s *Survey) error
}

type InMemoryRepository struct {
	surveys []*Survey
}

func (r *InMemoryRepository) Save(s *Survey) error {
	r.surveys = append(r.surveys, s)
	return nil
}

func SaveSurvey(s *Survey, r Repository) error {
	return r.Save(s)
}

// ===========

type SurveyManager struct {
	store Repository
}

func NewSurveyManager(store Repository) SurveyManager {
	return SurveyManager{
		store: store,
	}
}
