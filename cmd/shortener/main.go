package main

import (
	"context"
	"github.com/antonioo83/shot-url-service/config"
	"github.com/antonioo83/shot-url-service/internal/repositories/factory"
	"github.com/antonioo83/shot-url-service/internal/repositories/interfaces"
	"github.com/antonioo83/shot-url-service/internal/server"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
)

func main() {
	configSettings := config.GetConfigSettings()
	var pool *pgxpool.Pool
	context := context.Background()
	databaseRepository := factory.GetDatabaseRepository(configSettings)
	if configSettings.IsUseDatabase {
		pool, _ = pgxpool.Connect(context, configSettings.DatabaseDsn)
		defer pool.Close()
		err := databaseInit(databaseRepository, pool, configSettings.FilepathToDBDump)
		if err != nil {
			log.Fatal(err)
		}
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

func databaseInit(repository interfaces.DatabaseRepository, connect *pgxpool.Pool, filepathToDBDump string) error {
	return repository.RunDump(context.Background(), connect, filepathToDBDump)
}
