package server

import (
	"compress/flate"
	"github.com/antonioo83/shot-url-service/config"
	"github.com/antonioo83/shot-url-service/internal/handlers"
	authInterfaces "github.com/antonioo83/shot-url-service/internal/handlers/auth/interfaces"
	"github.com/antonioo83/shot-url-service/internal/repositories/interfaces"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

type RouteParameters struct {
	Config             config.Config
	ShotURLRepository  interfaces.ShotURLRepository
	UserRepository     interfaces.UserRepository
	DatabaseRepository interfaces.DatabaseRepository
	UserAuthHandler    authInterfaces.UserAuthHandler
}

func GetRouters(p RouteParameters) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	compressor := middleware.NewCompressor(flate.DefaultCompression)
	r.Use(compressor.Handler)

	r = getCreateShortURLRoute(r, p.Config, p.ShotURLRepository, p.UserRepository, p.UserAuthHandler)
	r = getCreateJSONShortURLRoute(r, p.Config, p.ShotURLRepository, p.UserRepository, p.UserAuthHandler)
	r = getOriginalURLRoute(r, p.ShotURLRepository)
	r = getUserUrlsRoute(r, p.ShotURLRepository, p.UserRepository, p.UserAuthHandler)
	r = getDatabaseStatus(r, p.DatabaseRepository)
	r = getCreateShortURLBatchRoute(r, p.Config, p.ShotURLRepository, p.UserRepository, p.UserAuthHandler)
	r = getDeleteShortURLRoute(r, p.Config, p.ShotURLRepository, p.UserAuthHandler)

	return r
}

func getCreateShortURLRoute(r *chi.Mux, config config.Config, repository interfaces.ShotURLRepository, userRepository interfaces.UserRepository,
	userAuthHandler authInterfaces.UserAuthHandler) *chi.Mux {
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCreateShortURLResponse(w, r, config, repository, userRepository, userAuthHandler)
	})

	return r
}

func getCreateJSONShortURLRoute(r *chi.Mux, config config.Config, repository interfaces.ShotURLRepository, userRepository interfaces.UserRepository,
	userAuthHandler authInterfaces.UserAuthHandler) *chi.Mux {
	r.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCreateJSONShortURLResponse(w, r, config, repository, userRepository, userAuthHandler)
	})

	return r
}

func getOriginalURLRoute(r *chi.Mux, repository interfaces.ShotURLRepository) *chi.Mux {
	r.Get("/{code}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetOriginalURLResponse(w, r, repository)
	})

	return r
}

func getUserUrlsRoute(r *chi.Mux, shotURLRepository interfaces.ShotURLRepository, userRepository interfaces.UserRepository,
	userAuthHandler authInterfaces.UserAuthHandler) *chi.Mux {
	r.Get("/api/user/urls", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUserURLsResponse(w, r, shotURLRepository, userRepository, userAuthHandler)
	})

	return r
}

func getDatabaseStatus(r *chi.Mux, databaseRepository interfaces.DatabaseRepository) *chi.Mux {
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetDBStatusResponse(w, databaseRepository)
	})

	return r
}

func getCreateShortURLBatchRoute(r *chi.Mux, config config.Config, repository interfaces.ShotURLRepository, userRepository interfaces.UserRepository,
	userAuthHandler authInterfaces.UserAuthHandler) *chi.Mux {
	r.Post("/api/shorten/batch", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCreateShortURLBatchResponse(w, r, config, repository, userRepository, userAuthHandler)
	})

	return r
}

func getDeleteShortURLRoute(r *chi.Mux, config config.Config, repository interfaces.ShotURLRepository, userAuthHandler authInterfaces.UserAuthHandler) *chi.Mux {
	jobCh := make(chan handlers.ShotURLDelete)
	runDeleteShortURLWorker(jobCh, repository, config.DeleteShotURL.WorkersCount)

	r.Delete("/api/user/urls", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetDeleteShortURLResponse(w, r, config, repository, userAuthHandler, jobCh)
	})

	return r
}

func runDeleteShortURLWorker(jobCh chan handlers.ShotURLDelete, repository interfaces.ShotURLRepository, workersCount int) {
	for i := 0; i < workersCount; i++ {
		go func() {
			for shotURLDelete := range jobCh {
				repository.Delete(shotURLDelete.UserCode, shotURLDelete.Codes)
			}
		}()
	}
}
