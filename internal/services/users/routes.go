package users

import (
	"encoding/json"
	"net/http"
	"tiny-goclean/internal/types"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc types.UserService
}

func NewHandler(svc types.UserService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *chi.Mux) {
	r.Route("/api/v1/users", func(r chi.Router) {
		r.Get("/", h.searchAll)
		r.Post("/", h.create)
	})
}

func (h *Handler) searchAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.Search()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, users)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateUserPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	user, err := h.svc.CreateUser(payload)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
