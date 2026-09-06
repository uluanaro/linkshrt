package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/uluanaro/linkshrt/internal/store"
)

func TestShorten(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{"valid", `{"url":"https://google.com"}`, http.StatusCreated},
		{"invalid", `{"url":"not-a-url}`, http.StatusBadRequest},
		{"empty", `{"url":""}`, http.StatusBadRequest},
		{"bad json", `abracadabra`, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New(store.New())
			req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()
			h.Shorten(rec, req)
			if rec.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}
		})
	}
}
func TestRedirect(t *testing.T) {
	s := store.New()
	code := s.Save("https://example.com")
	h := New(s)

	r := chi.NewRouter()
	r.Get("/{code}", h.Redirect)

	t.Run("existing code", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/"+code, nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusFound {
			t.Errorf("expected status %d, got %d", http.StatusFound, rec.Code)
		}
		if loc := rec.Header().Get("Location"); loc != "https://example.com" {
			t.Errorf("Location = %q, want https://example.com", loc)
		}
	})
	t.Run("invalid code", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/nonexisten"+code, nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Errorf("status = %d, want 404", rec.Code)
		}
	})
}
