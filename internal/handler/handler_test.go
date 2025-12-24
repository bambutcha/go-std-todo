package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-std-todo/internal/todo"
)

func TestHandler_CreateTodo(t *testing.T) {
	store := todo.NewStore()
	h := NewHandler(store)

	t.Run("successful creation", func(t *testing.T) {
		body := map[string]interface{}{
			"title":       "Test Todo",
			"description": "Test Description",
			"completed":   false,
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/todos", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.CreateTodo(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
		}

		var response todo.Todo
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.Title != "Test Todo" {
			t.Errorf("Expected title 'Test Todo', got '%s'", response.Title)
		}

		if response.ID == 0 {
			t.Error("Expected ID to be assigned, got 0")
		}
	})

	t.Run("validation error - empty title", func(t *testing.T) {
		body := map[string]interface{}{
			"title":       "",
			"description": "Test Description",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/todos", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.CreateTodo(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("duplicate ID error", func(t *testing.T) {
		body1 := map[string]interface{}{
			"id":    1,
			"title": "First Todo",
		}
		jsonBody1, _ := json.Marshal(body1)

		body2 := map[string]interface{}{
			"id":    1,
			"title": "Second Todo",
		}
		jsonBody2, _ := json.Marshal(body2)

		req1 := httptest.NewRequest("POST", "/todos", bytes.NewBuffer(jsonBody1))
		req1.Header.Set("Content-Type", "application/json")
		w1 := httptest.NewRecorder()
		h.CreateTodo(w1, req1)

		req2 := httptest.NewRequest("POST", "/todos", bytes.NewBuffer(jsonBody2))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		h.CreateTodo(w2, req2)

		if w2.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w2.Code)
		}
	})
}

func TestHandler_GetTodos(t *testing.T) {
	store := todo.NewStore()
	h := NewHandler(store)

	store.Create(&todo.Todo{Title: "Todo 1"})
	store.Create(&todo.Todo{Title: "Todo 2"})

	req := httptest.NewRequest("GET", "/todos", nil)
	w := httptest.NewRecorder()

	h.GetTodos(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var todos []todo.Todo
	if err := json.NewDecoder(w.Body).Decode(&todos); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(todos) != 2 {
		t.Errorf("Expected 2 todos, got %d", len(todos))
	}
}

func TestHandler_GetTodo(t *testing.T) {
	store := todo.NewStore()
	h := NewHandler(store)

	created, _ := store.Create(&todo.Todo{Title: "Test Todo"})

	t.Run("successful retrieval", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/todos/1", nil)
		w := httptest.NewRecorder()

		h.GetTodo(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response todo.Todo
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.ID != created.ID {
			t.Errorf("Expected ID %d, got %d", created.ID, response.ID)
		}
	})

	t.Run("not found error", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/todos/999", nil)
		w := httptest.NewRecorder()

		h.GetTodo(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestHandler_UpdateTodo(t *testing.T) {
	store := todo.NewStore()
	h := NewHandler(store)

	store.Create(&todo.Todo{Title: "Original Title"})

	t.Run("successful update", func(t *testing.T) {
		body := map[string]interface{}{
			"title":       "Updated Title",
			"description": "Updated Description",
			"completed":   true,
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("PUT", "/todos/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.UpdateTodo(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response todo.Todo
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.Title != "Updated Title" {
			t.Errorf("Expected title 'Updated Title', got '%s'", response.Title)
		}
	})

	t.Run("not found error", func(t *testing.T) {
		body := map[string]interface{}{"title": "New Title"}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("PUT", "/todos/999", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.UpdateTodo(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("validation error - empty title", func(t *testing.T) {
		body := map[string]interface{}{"title": ""}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("PUT", "/todos/1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.UpdateTodo(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestHandler_DeleteTodo(t *testing.T) {
	store := todo.NewStore()
	h := NewHandler(store)

	created, _ := store.Create(&todo.Todo{Title: "Test Todo"})

	t.Run("successful deletion", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/todos/1", nil)
		w := httptest.NewRecorder()

		h.DeleteTodo(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
		}

		_, err := store.GetByID(created.ID)
		if err == nil {
			t.Error("Expected todo to be deleted")
		}
	})

	t.Run("not found error", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/todos/999", nil)
		w := httptest.NewRecorder()

		h.DeleteTodo(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}
