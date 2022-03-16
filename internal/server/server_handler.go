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

func GetRouters(config config.Config, repository interfaces.ShotURLRepository) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	compressor := middleware.NewCompressor(flate.DefaultCompression)
	r.Use(compressor.Handler)

	r = getCreateShortURLRoute(r, config, repository)
	r = getCreateJSONShortURLRoute(r, config, repository)
	r = getOriginalURLRoute(r, repository)

	return r
}

func getCreateShortURLRoute(r *chi.Mux, config config.Config, repository interfaces.ShotURLRepository) *chi.Mux {
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCreateShortURLResponse(w, r, config, repository)
	})

	return r
}

func getCreateJSONShortURLRoute(r *chi.Mux, config config.Config, repository interfaces.ShotURLRepository) *chi.Mux {
	r.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCreateJSONShortURLResponse(w, r, config, repository)
	})

	return r
}

func getOriginalURLRoute(r *chi.Mux, repository interfaces.ShotURLRepository) *chi.Mux {
	r.Get("/{code}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetOriginalURLResponse(w, r, repository)
	})

	return r
}
