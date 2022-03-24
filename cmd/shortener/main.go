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
	var pool *pgxpool.Pool
	context := context.Background()
	configSettings := config.GetConfigSettings()
	if configSettings.IsUseDatabase {
		pool, _ = pgxpool.Connect(context, configSettings.DatabaseDsn)
		defer pool.Close()
	}

	databaseRepository := factory.GetDatabaseRepository(configSettings)
	shortUrlRepository := factory.GetRepository(context, pool, configSettings)
	userRepository := factory.GetUserRepository(context, pool, configSettings)
	if configSettings.IsUseDatabase {
		err := databaseInit(databaseRepository, pool, configSettings.FilepathToDBDump)
		if err != nil {
			log.Fatalln("can't load tables for the database:" + err.Error())
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
	log.Fatal(http.ListenAndServe(configSettings.ServerAddress, handler))
}

func databaseInit(repository interfaces.DatabaseRepository, connect *pgxpool.Pool, filepathToDBDump string) error {
	return repository.RunDump(context.Background(), connect, filepathToDBDump)
}
