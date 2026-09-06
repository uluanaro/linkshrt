package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/uluanaro/linkshrt/internal/auth"
	"github.com/uluanaro/linkshrt/internal/store"
	"golang.org/x/crypto/bcrypt"
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

type authRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type loginResponse struct {
	Token string `json:"token"`
}

type userResponse struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	login := req.Login
	password := req.Password
	err = validate.Struct(req)
	if err != nil {
		http.Error(w, "validation failed", http.StatusBadRequest)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to generate password", http.StatusInternalServerError)
		return
	}
	user, err := h.store.SaveUser(login, string(hash))
	if err != nil {
		http.Error(w, "failed to save user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userResponse{ID: user.ID, Login: user.Login})
}
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	login := req.Login
	password := req.Password
	err = validate.Struct(req)
	if err != nil {
		http.Error(w, "validation failed", http.StatusBadRequest)
		return
	}
	user, err := h.store.GetUser(login)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginResponse{Token: token})
}
