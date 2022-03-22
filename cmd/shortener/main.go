package main

import (
	"github.com/antonioo83/shot-url-service/config"
	"github.com/antonioo83/shot-url-service/internal/repositories/factory"
	"github.com/antonioo83/shot-url-service/internal/server"
	"log"
	"net/http"
)

func main() {
	configSettings := config.GetConfigSettings()
	routeParameters :=
		server.RouteParameters{
			Config:             configSettings,
			ShotURLRepository:  factory.GetRepository(configSettings),
			UserRepository:     factory.GetUserRepository(configSettings),
			DatabaseRepository: factory.GetDatabaseRepository(configSettings),
		}
	handler := server.GetRouters(routeParameters)
	log.Fatal(http.ListenAndServe(configSettings.ServerAddress, handler))
}
