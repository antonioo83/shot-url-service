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

func LoadModelsFromDatabase() bool {
	var model models.ShortURL
	_, err := filestore.LoadModels(localcache.Database, model, config.GetConfig())
	if err != nil {
		log.Fatal(err)

		return false
	}

	return true
}

func GetRouters() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r = getCreateShortURLRoute(r)
	r = getCreateJsonShortURLRoute(r)
	r = getOriginalURLRoute(r)

	return r
}

func getCreateShortURLRoute(r *chi.Mux) *chi.Mux {
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		originalURL, err := handlers.GetBody(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		getSavedShortUrlResponse(w, r, originalURL, func(w http.ResponseWriter, shotURL string) {
			handlers.GetCreateShortURLResponse(w, shotURL)
		})
	})

	return r
}

func getCreateJsonShortURLRoute(r *chi.Mux) *chi.Mux {
	r.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		originalURL, err := handlers.GetUrlParameter(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		getSavedShortUrlResponse(w, r, originalURL, func(w http.ResponseWriter, shotURL string) {
			handlers.GetCreateJsonShortURLResponse(w, shotURL)
		})
	})

	return r
}

func getSavedShortUrlResponse(w http.ResponseWriter, r *http.Request, originalURL string, responseFunc func(w http.ResponseWriter, shotURL string)) {
	shotURL, code, err := handlers.GetShortURL(originalURL, r, config.GetConfig().BaseUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if localcache.IsHasInDatabase(code) == true {
		responseFunc(w, shotURL)
		return
	}

	var shortURL models.ShortURL
	shortURL.Code = code
	shortURL.OriginalURL = originalURL
	shortURL.ShortURL = shotURL
	err = saveToStorage(shortURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseFunc(w, shotURL)
}

func saveToStorage(shortURL models.ShortURL) error {
	err := localcache.SaveURL(shortURL)
	if config.GetConfig().IsUseFileStore == true {
		err = filestore.SaveURL(shortURL, config.GetConfig())
	}

	return err
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

		handlers.GetOriginalURLResponse(w, model.OriginalURL)
	})

	return r
}
