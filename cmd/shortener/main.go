package main

import (
	"context"
	"fmt"
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
import _ "net/http/pprof"

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

// Review: Кажется, что абстрактная фабрика здесь избытчна. Инициализация зависимостей будет происходить у нас по коду только в одном месте - при старте приложения,
// где нам гораздое важнее прозрачность того, что у нас в итоге наинициализируется, нежели удобство интерфейса.
// давай уберем все абстрактные фабрики из main
//
// Answer: factories are used here as well as in tests and benchmarks. I think this solution makes code more readable and understandable.
// And I don't use duplicate code in other packages. But if it's important, I'll make it. Thank you for your remark!
func main() {
	fmt.Printf("Build version:%s\n", buildVersion)
	fmt.Printf("Build date:%s\n", buildDate)
	fmt.Printf("Build commit:%s\n", buildCommit)

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
