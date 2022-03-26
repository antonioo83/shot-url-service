package main

import (
	"context"
	"github.com/antonioo83/shot-url-service/config"
	"github.com/antonioo83/shot-url-service/internal/repositories/database"
	"github.com/antonioo83/shot-url-service/internal/repositories/factory"
	"github.com/antonioo83/shot-url-service/internal/server"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
)

func main() {
	var pool *pgxpool.Pool
	context := context.Background()
	configSettings := config.GetConfigSettings()
	databaseRepository := factory.GetDatabaseRepository(configSettings)
	if configSettings.IsUseDatabase {
		pool, err := database.GetConnectToInitializedDB(context, configSettings, databaseRepository)
		if err != nil {
			log.Fatalf("can't connect and initialize databse %v\n", err)
		}
		defer pool.Close()
	}

	routeParameters :=
		server.RouteParameters{
			Config:             configSettings,
			ShotURLRepository:  factory.GetRepository(context, pool, configSettings),
			UserRepository:     factory.GetUserRepository(context, pool, configSettings),
			DatabaseRepository: databaseRepository,
		}
	handler := server.GetRouters(routeParameters)
	log.Fatal(http.ListenAndServe(configSettings.ServerAddress, handler))
}
