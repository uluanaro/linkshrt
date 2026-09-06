package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/uluanaro/linkshrt/internal/store"
)

type Handler struct {
	store *store.Store
}

func New(s *store.Store) *Handler {
	return &Handler{
		store: s,
	}
}

type shortenRequest struct {
	URL string `json:"url" validate:"required,url"`
}
type shortenResponse struct {
	Short string `json:"short"`
}

var validate = validator.New()

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	err = validate.Struct(req)
	if err != nil {
		http.Error(w, "validation failed", http.StatusBadRequest)
		return
	}
	code := h.store.Save(req.URL)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(shortenResponse{Short: "http://localhost:8080/" + code})
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	url, ok := h.store.Get(code)
	if !ok {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}
