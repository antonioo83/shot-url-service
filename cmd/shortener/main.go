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
	databaseRepository := factory.GetDatabaseRepository(configSettings)

	var pool *pgxpool.Pool
	context := context.Background()
	if configSettings.IsUseDatabase {
		pool, _ = pgxpool.Connect(context, configSettings.DatabaseDsn) //databaseRepository.Connect(context)
		defer pool.Close()
	}

	shortUrlRepository := factory.GetRepository(context, pool, configSettings)
	userRepository := factory.GetUserRepository(context, pool, configSettings)
	if configSettings.IsUseDatabase {
		err := databaseInit(databaseRepository, pool, configSettings.FilepathToDBDump)
		if err != nil {
			log.Fatalln(err)
		}
	}
	routeParameters :=
		server.RouteParameters{
			Config:             configSettings,
			ShotURLRepository:  shortUrlRepository,
			UserRepository:     userRepository,
			DatabaseRepository: databaseRepository,
		}
	handler := server.GetRouters(routeParameters)
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func databaseInit(repository interfaces.DatabaseRepository, connect *pgxpool.Pool, filepathToDBDump string) error {
	return repository.RunDump(context.Background(), connect, filepathToDBDump)
}
