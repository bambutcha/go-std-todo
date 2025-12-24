package todo

import (
	"testing"
)

func TestStore_Create(t *testing.T) {
	store := NewStore()

	t.Run("successful creation", func(t *testing.T) {
		todo := &Todo{
			Title:       "Test Todo",
			Description: "Test Description",
			Completed:   false,
		}

		created, err := store.Create(todo)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if created.ID == 0 {
			t.Error("Expected ID to be assigned, got 0")
		}

		if created.Title != "Test Todo" {
			t.Errorf("Expected title 'Test Todo', got '%s'", created.Title)
		}
	})

	t.Run("validation error - empty title", func(t *testing.T) {
		todo := &Todo{
			Title:       "",
			Description: "Test Description",
		}

		_, err := store.Create(todo)
		if err == nil {
			t.Error("Expected validation error for empty title")
		}

		if err.Error() != "title cannot be empty" {
			t.Errorf("Expected error 'title cannot be empty', got '%s'", err.Error())
		}
	})

	t.Run("duplicate ID error", func(t *testing.T) {
		store := NewStore()
		todo1 := &Todo{
			ID:    1,
			Title: "First Todo",
		}

		todo2 := &Todo{
			ID:    1,
			Title: "Second Todo",
		}

		_, err := store.Create(todo1)
		if err != nil {
			t.Fatalf("Expected no error for first todo, got %v", err)
		}

		_, err = store.Create(todo2)
		if err == nil {
			t.Error("Expected error for duplicate ID")
		}

		if err.Error() != "todo with this ID already exists" {
			t.Errorf("Expected error 'todo with this ID already exists', got '%s'", err.Error())
		}
	})
}

func TestStore_GetByID(t *testing.T) {
	store := NewStore()

	todo := &Todo{
		Title: "Test Todo",
	}
	created, _ := store.Create(todo)

	t.Run("successful retrieval", func(t *testing.T) {
		retrieved, err := store.GetByID(created.ID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if retrieved.ID != created.ID {
			t.Errorf("Expected ID %d, got %d", created.ID, retrieved.ID)
		}
	})

	t.Run("not found error", func(t *testing.T) {
		_, err := store.GetByID(999)
		if err == nil {
			t.Error("Expected error for non-existent ID")
		}

		if err.Error() != "todo not found" {
			t.Errorf("Expected error 'todo not found', got '%s'", err.Error())
		}
	})
}

func TestStore_GetAll(t *testing.T) {
	store := NewStore()

	todo1 := &Todo{Title: "Todo 1"}
	todo2 := &Todo{Title: "Todo 2"}

	store.Create(todo1)
	store.Create(todo2)

	todos := store.GetAll()
	if len(todos) != 2 {
		t.Errorf("Expected 2 todos, got %d", len(todos))
	}
}

func TestStore_Update(t *testing.T) {
	store := NewStore()

	todo := &Todo{
		Title: "Original Title",
	}
	created, _ := store.Create(todo)

	t.Run("successful update", func(t *testing.T) {
		updatedTodo := &Todo{
			Title:       "Updated Title",
			Description: "Updated Description",
			Completed:   true,
		}

		updated, err := store.Update(created.ID, updatedTodo)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if updated.Title != "Updated Title" {
			t.Errorf("Expected title 'Updated Title', got '%s'", updated.Title)
		}

		if updated.ID != created.ID {
			t.Errorf("Expected ID to remain %d, got %d", created.ID, updated.ID)
		}
	})

	t.Run("not found error", func(t *testing.T) {
		updatedTodo := &Todo{Title: "New Title"}
		_, err := store.Update(999, updatedTodo)
		if err == nil {
			t.Error("Expected error for non-existent ID")
		}

		if err.Error() != "todo not found" {
			t.Errorf("Expected error 'todo not found', got '%s'", err.Error())
		}
	})

	t.Run("validation error - empty title", func(t *testing.T) {
		updatedTodo := &Todo{Title: ""}
		_, err := store.Update(created.ID, updatedTodo)
		if err == nil {
			t.Error("Expected validation error for empty title")
		}

		if err.Error() != "title cannot be empty" {
			t.Errorf("Expected error 'title cannot be empty', got '%s'", err.Error())
		}
	})
}

func TestStore_Delete(t *testing.T) {
	store := NewStore()

	todo := &Todo{Title: "Test Todo"}
	created, _ := store.Create(todo)

	t.Run("successful deletion", func(t *testing.T) {
		err := store.Delete(created.ID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		_, err = store.GetByID(created.ID)
		if err == nil {
			t.Error("Expected todo to be deleted")
		}
	})

	t.Run("not found error", func(t *testing.T) {
		err := store.Delete(999)
		if err == nil {
			t.Error("Expected error for non-existent ID")
		}

		if err.Error() != "todo not found" {
			t.Errorf("Expected error 'todo not found', got '%s'", err.Error())
		}
	})
}
