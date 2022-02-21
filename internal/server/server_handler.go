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

	r = getOriginalUrlRoute(r)
	r = getCreateShortUrlRoute(r)

	return r
}

func getCreateShortUrlRoute(r *chi.Mux) *chi.Mux {
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		originalUrl, err := handlers.GetBody(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		shotUrl, code, err := handlers.GetShortUrl(originalUrl, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var shortUrl models.ShortUrl
		shortUrl.Code = code
		shortUrl.OriginalUrl = originalUrl
		shortUrl.ShortUrl = shotUrl
		localcache.SaveUrl(shortUrl)

		w.WriteHeader(http.StatusCreated)
		logErr(w.Write([]byte(shotUrl)))
	})

	return r
}

func getOriginalUrlRoute(r *chi.Mux) *chi.Mux {
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

		w.Header().Set("Location", model.OriginalUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
		logErr(w.Write([]byte(model.OriginalUrl)))
	})

	return r
}

func logErr(n int, err error) int {
	if err != nil {
		log.Printf("Write failed: %v", err)
	}

	return n
}
