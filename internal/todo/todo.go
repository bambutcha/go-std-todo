package todo

import (
	"errors"
	"sync"
)

type Todo struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type Store struct {
	mu     sync.RWMutex
	todos  map[int]*Todo
	nextID int
}

func NewStore() *Store {
	return &Store{
		todos:  make(map[int]*Todo),
		nextID: 1,
	}
}

func (s *Store) Create(todo *Todo) (*Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if todo.Title == "" {
		return nil, errors.New("title cannot be empty")
	}

	if todo.ID != 0 {
		if _, exists := s.todos[todo.ID]; exists {
			return nil, errors.New("todo with this ID already exists")
		}
		s.todos[todo.ID] = todo
		return todo, nil
	}

	todo.ID = s.nextID
	s.nextID++
	s.todos[todo.ID] = todo
	return todo, nil
}

func (s *Store) GetByID(id int) (*Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todo, exists := s.todos[id]
	if !exists {
		return nil, errors.New("todo not found")
	}
	return todo, nil
}

func (s *Store) GetAll() []*Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todos := make([]*Todo, 0, len(s.todos))
	for _, todo := range s.todos {
		todos = append(todos, todo)
	}
	return todos
}

func (s *Store) Update(id int, todo *Todo) (*Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.todos[id]; !exists {
		return nil, errors.New("todo not found")
	}

	if todo.Title == "" {
		return nil, errors.New("title cannot be empty")
	}

	todo.ID = id
	s.todos[id] = todo
	return todo, nil
}

func (s *Store) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.todos[id]; !exists {
		return errors.New("todo not found")
	}

	delete(s.todos, id)
	return nil
}
