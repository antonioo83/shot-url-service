package main

import (
	"context"
	"github.com/antonioo83/shot-url-service/config"
	"github.com/antonioo83/shot-url-service/internal/repositories/factory"
	"github.com/antonioo83/shot-url-service/internal/server"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
)

func main() {
	configSettings := config.GetConfigSettings()
	var pool *pgxpool.Pool
	context := context.Background()
	routeParameters :=
		server.RouteParameters{
			Config:             configSettings,
			ShotURLRepository:  factory.GetRepository(context, pool, configSettings),
			UserRepository:     factory.GetUserRepository(context, pool, configSettings),
			DatabaseRepository: factory.GetDatabaseRepository(configSettings),
		}
	handler := server.GetRouters(routeParameters)
	log.Fatal(http.ListenAndServe(configSettings.ServerAddress, handler))
}
