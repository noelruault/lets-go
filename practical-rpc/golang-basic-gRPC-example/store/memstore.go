package store

import (
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

type InMemStore struct {
	Films map[string]Film
}

type Film struct {
	ID          string
	Title       string
	Director    string
	Producer    string
	ReleaseDate time.Time
}

func defaultFilms() map[string]Film {
	return map[string]Film{
		"c849f921-5488-4750-b989-70ce002eb572": {
			ID:          "c849f921-5488-4750-b989-70ce002eb572",
			Title:       "A New Hope",
			Director:    "George Lucas",
			Producer:    "Gary Kurtz, Rick McCallum",
			ReleaseDate: time.Date(1977, time.Month(5), 25, 0, 0, 0, 0, time.Local),
		},
		"ff686c7d-2e09-4920-81e4-50e28731094c": {
			ID:          "ff686c7d-2e09-4920-81e4-50e28731094c",
			Title:       "The Empire Strikes Back",
			Director:    "Irvin Kershner",
			Producer:    "Gary Kurtz, Rick McCallum",
			ReleaseDate: time.Date(1980, time.Month(5), 17, 0, 0, 0, 0, time.Local),
		},
		"35066997-1cb2-47d4-9ccf-d45d11f7c3ee": {
			ID:          "35066997-1cb2-47d4-9ccf-d45d11f7c3ee",
			Title:       "Return of the Jedi",
			Director:    "Richard Marquand",
			Producer:    "Howard G. Kazanjian, George Lucas, Rick McCallum",
			ReleaseDate: time.Date(1983, time.Month(5), 25, 0, 0, 0, 0, time.Local),
		},
	}
}

func NewInMemStore() *InMemStore {
	return &InMemStore{
		Films: defaultFilms(),
	}
}

// Create
func (s *InMemStore) Create(f Film) Film {

	f.ID = uuid.NewV4().String()
	// Check if ID already exists.
	if _, ok := s.Films[f.ID]; !ok {
		s.Films[f.ID] = f
		return f
	}

	return Film{} // TODO(nruault): Send here an error too?
}

// List
func (s *InMemStore) List() map[string]Film {
	return s.Films
}

// Retrieve
func (s *InMemStore) Get(id string) (Film, error) {
	v, ok := s.Films[id]
	if !ok {
		return Film{}, errors.New("not_found") // 404 not found
		// https://github.com/grpc/grpc/blob/master/doc/statuscodes.md
	}
	return v, nil
}

// Update
func (s *InMemStore) Update(f Film) (Film, error) {

	// Check if film contains ID
	if f.ID == "" {
		return Film{}, errors.New("bad_request") // 400 bad_request
	}

	_, ok := s.Films[f.ID]
	if !ok {
		return Film{}, errors.New("not_found") // 404 not_found
	}

	s.Films[f.ID] = f // Change values to the requested ones

	return s.Films[f.ID], nil
}

// Delete
func (s *InMemStore) Delete(id string) error {
	if _, ok := s.Films[id]; !ok {
		return errors.New("not_found") // 404 not found
	}
	delete(s.Films, id)

	// if _, ok := s.Films[id]; !ok {
	// 	fmt.Println("REMOVED: " + id)
	// }

	return nil
}
