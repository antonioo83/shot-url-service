package server

import (
	"github.com/antonioo83/shot-url-service/config"
	"github.com/antonioo83/shot-url-service/internal/handlers"
	"github.com/antonioo83/shot-url-service/internal/models"
	"github.com/antonioo83/shot-url-service/internal/repositories/filestore"
	"github.com/antonioo83/shot-url-service/internal/repositories/localcache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"time"
)

func LoadModelsFromDatabase(config config.Config) bool {
	var model models.ShortURL
	_, err := filestore.LoadModels(localcache.Database, model, config)
	if err != nil {
		log.Fatal(err)

		return false
	}

	return true
}

func GetRouters(config config.Config) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r = getCreateShortURLRoute(r, config)
	r = getCreateJSONShortURLRoute(r, config)
	r = getOriginalURLRoute(r)

	return r
}

func getCreateShortURLRoute(r *chi.Mux, config config.Config) *chi.Mux {
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		originalURL, err := handlers.GetBody(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		getSavedShortURLResponse(w, r, config, originalURL, func(w http.ResponseWriter, shotURL string) {
			handlers.GetCreateShortURLResponse(w, shotURL)
		})
	})

	return r
}

func getCreateJSONShortURLRoute(r *chi.Mux, config config.Config) *chi.Mux {
	r.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		originalURL, err := handlers.GetOriginalURLFromBody(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		getSavedShortURLResponse(w, r, config, originalURL, func(w http.ResponseWriter, shotURL string) {
			handlers.GetCreateJSONShortURLResponse(w, shotURL)
		})
	})

	return r
}

func getSavedShortURLResponse(w http.ResponseWriter, r *http.Request, config config.Config, originalURL string, responseFunc func(w http.ResponseWriter, shotURL string)) {
	shotURL, code, err := handlers.GetShortURL(originalURL, r, config.BaseURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if localcache.IsInDatabase(code) {
		responseFunc(w, shotURL)
		return
	}

	var shortURL models.ShortURL
	shortURL.Code = code
	shortURL.OriginalURL = originalURL
	shortURL.ShortURL = shotURL
	err = saveToStorage(config, shortURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseFunc(w, shotURL)
}

func saveToStorage(config config.Config, shortURL models.ShortURL) error {
	err := localcache.SaveURL(shortURL)
	if config.IsUseFileStore {
		err = filestore.SaveURL(shortURL, config)
	}

	return err
}

func getOriginalURLRoute(r *chi.Mux) *chi.Mux {
	r.Get("/{code}", func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")
		if code == "" {
			http.Error(w, "httpStatus param is missed", http.StatusBadRequest)
			return
		}
		model, err := localcache.FindByCode(code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		handlers.GetOriginalURLResponse(w, model.OriginalURL)
	})

	return r
}
