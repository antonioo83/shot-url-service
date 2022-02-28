package server

import (
	"github.com/antonioo83/shot-url-service/internal/handlers"
	"github.com/antonioo83/shot-url-service/internal/models"
	"github.com/antonioo83/shot-url-service/internal/repositories/localcache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"time"
)

func GetRouters() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r = getOriginalURLRoute(r)
	r = getCreateShortURLRoute(r)
	r = getCreateJsonShortURLRoute(r)

	return r
}

func getCreateJsonShortURLRoute(r *chi.Mux) *chi.Mux {
	r.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		originalURL, err := handlers.GetUrlParameter(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		shotURL, code, err := handlers.GetShortURL(originalURL, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var shortURL models.ShortURL
		shortURL.Code = code
		shortURL.OriginalURL = originalURL
		shortURL.ShortURL = shotURL
		localcache.SaveURL(shortURL)

		w.WriteHeader(http.StatusCreated)
		logErr(w.Write(handlers.GetJsonResponse("result", shotURL)))
	})

	return r
}

func getCreateShortURLRoute(r *chi.Mux) *chi.Mux {
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		originalURL, err := handlers.GetBody(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		shotURL, code, err := handlers.GetShortURL(originalURL, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var shortURL models.ShortURL
		shortURL.Code = code
		shortURL.OriginalURL = originalURL
		shortURL.ShortURL = shotURL
		localcache.SaveURL(shortURL)

		w.WriteHeader(http.StatusCreated)
		logErr(w.Write([]byte(shotURL)))
	})

	return r
}

func getOriginalURLRoute(r *chi.Mux) *chi.Mux {
	r.Get("/{httpStatus}", func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "httpStatus")
		if code == "" {
			http.Error(w, "httpStatus param is missed", http.StatusBadRequest)
			return
		}
		model, err := localcache.FindByCode(code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", model.OriginalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		logErr(w.Write([]byte(model.OriginalURL)))
	})

	return r
}

func logErr(n int, err error) int {
	if err != nil {
		log.Printf("Write failed: %v", err)
	}

	return n
}
