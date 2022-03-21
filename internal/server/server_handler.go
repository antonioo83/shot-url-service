package server

import (
	"compress/flate"
	"github.com/antonioo83/shot-url-service/config"
	"github.com/antonioo83/shot-url-service/internal/handlers"
	"github.com/antonioo83/shot-url-service/internal/repositories/interfaces"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

func GetRouters(config config.Config, shotURLRepository interfaces.ShotURLRepository, userRepository interfaces.UserRepository) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	compressor := middleware.NewCompressor(flate.DefaultCompression)
	r.Use(compressor.Handler)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			next.ServeHTTP(w, r)
		})
	})

	r = getCreateShortURLRoute(r, config, shotURLRepository, userRepository)
	r = getCreateJSONShortURLRoute(r, config, shotURLRepository, userRepository)
	r = getOriginalURLRoute(r, shotURLRepository, userRepository)
	r = getUserUrlsRoute(r, shotURLRepository, userRepository)

	return r
}

func getCreateShortURLRoute(r *chi.Mux, config config.Config, repository interfaces.ShotURLRepository, userRepository interfaces.UserRepository) *chi.Mux {
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCreateShortURLResponse(w, r, config, repository, userRepository)
	})

	return r
}

func getCreateJSONShortURLRoute(r *chi.Mux, config config.Config, repository interfaces.ShotURLRepository, userRepository interfaces.UserRepository) *chi.Mux {
	r.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCreateJSONShortURLResponse(w, r, config, repository, userRepository)
	})

	return r
}

func getOriginalURLRoute(r *chi.Mux, repository interfaces.ShotURLRepository, userRepository interfaces.UserRepository) *chi.Mux {
	r.Get("/{code}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetOriginalURLResponse(w, r, repository)
	})

	return r
}

func getUserUrlsRoute(r *chi.Mux, shotURLRepository interfaces.ShotURLRepository, userRepository interfaces.UserRepository) *chi.Mux {
	r.Get("/api/user/urls", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUserURLsResponse(w, r, shotURLRepository, userRepository)
	})

	return r
}
