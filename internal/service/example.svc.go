package service

import (
	"fmt"
	"sync"

	"github.com/mhirii/rest-template/internal/dto"
)

// ExampleStore is a simple in-memory store for Example resources.
type ExampleStore struct {
	mu   sync.RWMutex
	data map[string]*dto.Example
}

// NewExampleStore creates a new ExampleStore.
func NewExampleStore() *ExampleStore {
	return &ExampleStore{data: make(map[string]*dto.Example)}
}

// Create adds a new Example to the store and returns it with a generated ID.
func (s *ExampleStore) Create(e dto.Example) dto.Example {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := generateID(len(s.data))
	e.ID = id
	s.data[id] = &e
	return e
}

// Replace sets the Example with the given ID, creating or overwriting it.
func (s *ExampleStore) Replace(id string, e dto.Example) dto.Example {
	s.mu.Lock()
	defer s.mu.Unlock()
	e.ID = id
	s.data[id] = &e
	return e
}

// Update updates fields of the Example with the given ID.
func (s *ExampleStore) Update(id string, name *string) (dto.Example, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.data[id]
	if !ok {
		return dto.Example{}, false
	}
	if name != nil {
		e.Name = *name
	}
	return *e, true
}

// Get returns the Example with the given ID.
func (s *ExampleStore) Get(id string) (dto.Example, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.data[id]
	if !ok {
		return dto.Example{}, false
	}
	return *e, true
}

// List returns all Examples in the store.
func (s *ExampleStore) List() []dto.Example {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]dto.Example, 0, len(s.data))
	for _, e := range s.data {
		list = append(list, *e)
	}
	return list
}

// Delete removes the Example with the given ID.
func (s *ExampleStore) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[id]; !ok {
		return false
	}
	delete(s.data, id)
	return true
}

// generateID creates a simple unique ID for demonstration.
func generateID(n int) string {
	return "ex-" + fmt.Sprint(n+1)
}
