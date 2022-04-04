package main

import (
	"context"
	"github.com/antonioo83/shot-url-service/config"
	authFactory "github.com/antonioo83/shot-url-service/internal/handlers/auth/factory"
	"github.com/antonioo83/shot-url-service/internal/repositories/factory"
	"github.com/antonioo83/shot-url-service/internal/repositories/interfaces"
	"github.com/antonioo83/shot-url-service/internal/server"
	"github.com/go-chi/jwtauth"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
)

func main() {
	config := config.GetConfigSettings()
	var tokenAuth *jwtauth.JWTAuth
	var pool *pgxpool.Pool
	context := context.Background()
	databaseRepository := factory.GetDatabaseRepository(config)
	if config.IsUseDatabase {
		pool, _ = pgxpool.Connect(context, config.DatabaseDsn)
		defer pool.Close()
		err := databaseInit(databaseRepository, pool, config.FilepathToDBDump)
		if err != nil {
			log.Fatal(err)
		}
	}
	userRepository := factory.GetUserRepository(context, pool, config)
	routeParameters :=
		server.RouteParameters{
			Config:             config,
			ShotURLRepository:  factory.GetRepository(context, pool, config),
			UserRepository:     userRepository,
			DatabaseRepository: databaseRepository,
			UserAuthHandler:    authFactory.NewAuthHandler(tokenAuth, userRepository, config),
		}
	handler := server.GetRouters(routeParameters)
	log.Fatal(http.ListenAndServe(config.ServerAddress, handler))
}

func databaseInit(repository interfaces.DatabaseRepository, connect *pgxpool.Pool, filepathToDBDump string) error {
	return repository.RunDump(context.Background(), connect, filepathToDBDump)
}
