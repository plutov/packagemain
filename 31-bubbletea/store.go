package main

type Note struct {
	ID    uint
	Title string
	Body  string
}

type Store interface {
	GetNotes() ([]Note, error)
	SaveNote(note Note) error
}

type InMemStore struct {
	notes []Note
}

func (s *InMemStore) GetNotes() ([]Note, error) {
	return s.notes, nil
}

func (s *InMemStore) SaveNote(note Note) error {
	s.notes = append(s.notes, note)
	return nil
}
