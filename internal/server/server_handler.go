// Package server This package is intended for service configuration.
package server

import (
	"compress/flate"
	"github.com/antonioo83/shot-url-service/config"
	"github.com/antonioo83/shot-url-service/internal/handlers"
	authInterfaces "github.com/antonioo83/shot-url-service/internal/handlers/auth/interfaces"
	genInterfaces "github.com/antonioo83/shot-url-service/internal/handlers/generators/interfaces"
	"github.com/antonioo83/shot-url-service/internal/repositories/interfaces"
	"github.com/antonioo83/shot-url-service/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"net/http/pprof"
	"time"
)

type RouteParameters struct {
	Config             config.Config
	ShotURLRepository  interfaces.ShotURLRepository
	UserRepository     interfaces.UserRepository
	DatabaseRepository interfaces.DatabaseRepository
	UserAuthHandler    authInterfaces.UserAuthHandler
	Generator          genInterfaces.ShortLinkGenerator
}

// GetRouters Returns all available routers.
func GetRouters(p RouteParameters) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	compressor := middleware.NewCompressor(flate.DefaultCompression)
	r.Use(compressor.Handler)

	r = GetCreateShortURLRoute(r, p.Config, p.ShotURLRepository, p.UserRepository, p.UserAuthHandler, p.Generator)
	r = GetCreateJSONShortURLRoute(r, p.Config, p.ShotURLRepository, p.UserRepository, p.UserAuthHandler, p.Generator)
	r = GetOriginalURLRoute(r, p.ShotURLRepository)
	r = GetUserUrlsRoute(r, p.ShotURLRepository, p.UserRepository, p.UserAuthHandler)
	r = GetDatabaseStatus(r, p.DatabaseRepository)
	r = GetCreateShortURLBatchRoute(r, p.Config, p.ShotURLRepository, p.UserRepository, p.UserAuthHandler, p.Generator)
	r = GetDeleteShortURLRoute(r, p.Config, p.ShotURLRepository, p.UserAuthHandler)
	r = GetStatsRoute(r, p.Config, p.ShotURLRepository, p.UserRepository)
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/allocs", pprof.Index)
	r.HandleFunc("/debug/pprof/block", pprof.Index)
	r.HandleFunc("/debug/pprof/goroutine", pprof.Index)
	r.HandleFunc("/debug/pprof/heap", pprof.Index)
	r.HandleFunc("/debug/pprof/mutex", pprof.Index)
	r.HandleFunc("/debug/pprof/threadcreate", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)

	return r
}

// GetCreateShortURLRoute Returns a route to create short url.
//
// POST http://localhost:8080/
func GetCreateShortURLRoute(r *chi.Mux, config config.Config, repository interfaces.ShotURLRepository, userRepository interfaces.UserRepository,
	userAuthHandler authInterfaces.UserAuthHandler, generator genInterfaces.ShortLinkGenerator) *chi.Mux {
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCreateShortURLResponse(w, r, config, repository, userRepository, userAuthHandler, generator)
	})

	return r
}

// GetCreateJSONShortURLRoute Returns a route to get short url in json format.
//
// POST http://localhost:8080/api/shorten
func GetCreateJSONShortURLRoute(r *chi.Mux, config config.Config, repository interfaces.ShotURLRepository, userRepository interfaces.UserRepository,
	userAuthHandler authInterfaces.UserAuthHandler, generator genInterfaces.ShortLinkGenerator) *chi.Mux {
	r.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCreateJSONShortURLResponse(w, r, config, repository, userRepository, userAuthHandler, generator)
	})

	return r
}

// GetOriginalURLRoute Returns a route to get short url by code.
//
// GET http://localhost:8080/code
func GetOriginalURLRoute(r *chi.Mux, repository interfaces.ShotURLRepository) *chi.Mux {
	r.Get("/{code}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetOriginalURLResponse(w, r, repository)
	})

	return r
}

// GetUserUrlsRoute Returns a route to get short url as array.
//
// GET http://localhost:8080/api/user/urls
func GetUserUrlsRoute(r *chi.Mux, shotURLRepository interfaces.ShotURLRepository, userRepository interfaces.UserRepository,
	userAuthHandler authInterfaces.UserAuthHandler) *chi.Mux {
	r.Get("/api/user/urls", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetUserURLsResponse(w, r, shotURLRepository, userRepository, userAuthHandler)
	})

	return r
}

// GetStatsRoute returns a route to get count of shot urls and count of users saved in the databases.
//
// GET http://localhost:8080/api/internal/stats
func GetStatsRoute(r *chi.Mux, config config.Config, shotURLRepository interfaces.ShotURLRepository, userRepository interfaces.UserRepository) *chi.Mux {
	r.Get("/api/internal/stats", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetStatsResponse(w, r, config, shotURLRepository, userRepository)
	})

	return r
}

// GetDatabaseStatus Returns a route to get database status.
//
// GET http://localhost:8080/ping
func GetDatabaseStatus(r *chi.Mux, databaseRepository interfaces.DatabaseRepository) *chi.Mux {
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetDBStatusResponse(w, databaseRepository)
	})

	return r
}

// GetCreateShortURLBatchRoute Returns a route to create an array of short URLs.
//
// POST http://localhost:8080/api/shorten/batch
func GetCreateShortURLBatchRoute(r *chi.Mux, config config.Config, repository interfaces.ShotURLRepository, userRepository interfaces.UserRepository,
	userAuthHandler authInterfaces.UserAuthHandler, generator genInterfaces.ShortLinkGenerator) *chi.Mux {
	r.Post("/api/shorten/batch", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetCreateShortURLBatchResponse(w, r, config, repository, userRepository, userAuthHandler, generator)
	})

	return r
}

// GetDeleteShortURLRoute Returns a route to delete an array of short URLs.
//
// DELETE http://localhost:8080/api/user/urls
func GetDeleteShortURLRoute(r *chi.Mux, config config.Config, repository interfaces.ShotURLRepository, userAuthHandler authInterfaces.UserAuthHandler) *chi.Mux {
	jobCh := make(chan services.ShotURLDelete)
	services.RunDeleteShortURLWorker(jobCh, repository, config.DeleteShotURL.WorkersCount)

	r.Delete("/api/user/urls", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetDeleteShortURLResponse(w, r, config, repository, userAuthHandler, jobCh)
	})

	return r
}
