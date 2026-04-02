package users

import (
	"encoding/json"
	"net/http"
	"tiny-goclean/internal/types"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(r *chi.Mux) {
	r.Route("/api/v1/users", func(r chi.Router) {
		r.Get("/", h.searchAll)
	})
}

func (h *Handler) searchAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	b, _ := json.MarshalIndent(map[string]string{"okoko from status api": "ok"}, " ", "")
	w.Write(b)
}
