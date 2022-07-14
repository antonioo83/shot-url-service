package server

import (
	"github.com/go-chi/chi/v5"
)

func ExampleGetCreateShortURLRoute() {
	r := chi.NewRouter()
	p := RouteParameters{}
	r = GetCreateShortURLRoute(r, p.Config, p.ShotURLRepository, p.UserRepository, p.UserAuthHandler, p.Generator)
	// Output (parallel):
	// # Request
	// POST http://localhost:8080/
	// Content-Type: text/plain; charset=utf-8
	//
	// https://stackoverflow.com/questions/15240884/in-go
	//
	// # Response
	// HTTP/1.1 201 OK
	// Content-Type: text/plain; charset=utf-8
	//
	// http://localhost:8080/i1rYMHSU
}

func ExampleGetCreateJSONShortURLRoute() {
	r := chi.NewRouter()
	p := RouteParameters{}
	r = GetCreateJSONShortURLRoute(r, p.Config, p.ShotURLRepository, p.UserRepository, p.UserAuthHandler, p.Generator)
	// Output (parallel):
	// # Request
	// POST http://localhost:8080/api/shorten
	// Content-Type: application/json
	//
	// {"url": "https://stackoverflow.com/questions/15240884/in-go"}
	//
	// # Response
	// HTTP/1.1 201 OK
	// Content-Type: application/json
	//
	// {"result": "http://localhost:8080/9yBBZ3nW"}
}

func ExampleGetOriginalURLRoute() {
	r := chi.NewRouter()
	p := RouteParameters{}
	r = GetOriginalURLRoute(r, p.ShotURLRepository)
	// Output (parallel):
	// # Request
	// GET http://localhost:8080/code
	//
	// # Response
	// HTTP/1.1 307 OK
}

func ExampleGetUserUrlsRoute() {
	r := chi.NewRouter()
	p := RouteParameters{}
	r = GetUserUrlsRoute(r, p.ShotURLRepository, p.UserRepository, p.UserAuthHandler)
	// Output (parallel):
	// # Request
	// GET http://localhost:8080/api/user/urls
	//
	// # Response
	// HTTP/1.1 200 OK
	// Cookie: token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2NvZGUiOjF9.6dg8iJJI-CPFT6Bn4Dk3h0zITD4C7agRvUYJLszW5mI
	//
	// [{"short_url": "http://localhost:8080/i1rYMHSU","original_url": "https://stackoverflow.com/questions/15240884/in-go"}]
}

func ExampleGetDatabaseStatus() {
	r := chi.NewRouter()
	p := RouteParameters{}
	r = GetDatabaseStatus(r, p.DatabaseRepository)
	// Output (parallel):
	// # Request
	// GET http://localhost:8080/ping
	//
	// # Response
	// HTTP/1.1 200 OK
	// HTTP/1.1 500 OK
}

func ExampleGetCreateShortURLBatchRoute() {
	r := chi.NewRouter()
	p := RouteParameters{}
	r = GetCreateShortURLBatchRoute(r, p.Config, p.ShotURLRepository, p.UserRepository, p.UserAuthHandler, p.Generator)
	// Output (parallel):
	// # Request
	// POST http://localhost:8080/api/shorten/batch
	// Content-Type: application/json
	//
	// [{"correlation_id":"463186fc-72c8-4204-ae3c-48359c2f63bd","original_url":"http://rbc.ru/"}]
	//
	// # Response
	// HTTP/1.1 201 OK
	// Content-Type: application/json
	//
	// [{"correlation_id": "463186fc-72c8-4204-ae3c-48359c2f63bd","short_url": "http://localhost:8080/WhGuTjTu"}]
}
