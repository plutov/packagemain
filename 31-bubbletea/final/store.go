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
	for i, n := range s.notes {
		if n.ID == note.ID {
			s.notes[i] = note
			return nil
		}
	}

	note.ID = uint(len(s.notes) + 1)
	s.notes = append(s.notes, note)
	return nil
}
