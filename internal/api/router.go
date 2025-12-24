package api

import (
	"net/http"

	"go-std-todo/internal/handler"
)

func NewRouter(h *handler.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /todos", h.CreateTodo)
	mux.HandleFunc("GET /todos", h.GetTodos)
	mux.HandleFunc("GET /todos/", h.GetTodo)
	mux.HandleFunc("PUT /todos/", h.UpdateTodo)
	mux.HandleFunc("DELETE /todos/", h.DeleteTodo)
	return mux
}
