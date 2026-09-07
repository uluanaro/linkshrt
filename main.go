package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/uluanaro/linkshrt/internal/auth"
	"github.com/uluanaro/linkshrt/internal/handler"
	"github.com/uluanaro/linkshrt/internal/middleware"
	"github.com/uluanaro/linkshrt/internal/store"
	"golang.org/x/time/rate"
)

func main() {
	st := store.New()
	h := handler.New(st)

	router := chi.NewRouter()
	router.Use(chimw.Logger)
	router.Use(chimw.Recoverer)

	rl := middleware.NewRateLimiter(rate.Limit(2), 2)
	router.Use(rl.Middleware)

	router.Group(func(router chi.Router) {
		router.Use(auth.AuthMiddleware)
		router.Post("/shorten", h.Shorten)
	})
	router.Get("/{code}", h.Redirect)
	router.Post("/register", h.Register)
	router.Post("/login", h.Login)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  100 * time.Second,
	}
	go func() {
		log.Printf("Listening on port 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	log.Println("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped")
}
