package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"go-std-todo/internal/todo"
)

type Handler struct {
	store *todo.Store
}

func NewHandler(store *todo.Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)

	var t todo.Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if t.Title == "" {
		http.Error(w, "Title cannot be empty", http.StatusBadRequest)
		return
	}

	select {
	case <-ctx.Done():
		http.Error(w, "Request timeout", http.StatusRequestTimeout)
		return
	default:
	}

	created, err := h.store.Create(&t)
	if err != nil {
		if err.Error() == "todo with this ID already exists" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *Handler) GetTodos(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)

	select {
	case <-ctx.Done():
		http.Error(w, "Request timeout", http.StatusRequestTimeout)
		return
	default:
	}

	todos := h.store.GetAll()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func (h *Handler) GetTodo(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)

	id, err := extractID(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	select {
	case <-ctx.Done():
		http.Error(w, "Request timeout", http.StatusRequestTimeout)
		return
	default:
	}

	t, err := h.store.GetByID(id)
	if err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func (h *Handler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)

	id, err := extractID(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var t todo.Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if t.Title == "" {
		http.Error(w, "Title cannot be empty", http.StatusBadRequest)
		return
	}

	select {
	case <-ctx.Done():
		http.Error(w, "Request timeout", http.StatusRequestTimeout)
		return
	default:
	}

	updated, err := h.store.Update(id, &t)
	if err != nil {
		if err.Error() == "todo not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func (h *Handler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	log.Printf("[%s] %s %s", r.Method, r.URL.Path, r.RemoteAddr)

	id, err := extractID(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	select {
	case <-ctx.Done():
		http.Error(w, "Request timeout", http.StatusRequestTimeout)
		return
	default:
	}

	if err := h.store.Delete(id); err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func extractID(path string) (int, error) {
	if len(path) < 8 || path[:7] != "/todos/" {
		return 0, &strconv.NumError{Func: "Atoi", Num: path, Err: strconv.ErrSyntax}
	}
	idStr := path[7:]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return id, nil
}
