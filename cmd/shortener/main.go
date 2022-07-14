package main

import (
	"context"
	"fmt"
	"github.com/antonioo83/shot-url-service/config"
	authFactory "github.com/antonioo83/shot-url-service/internal/handlers/auth/factory"
	"github.com/antonioo83/shot-url-service/internal/handlers/generators"
	"github.com/antonioo83/shot-url-service/internal/repositories/factory"
	"github.com/antonioo83/shot-url-service/internal/repositories/interfaces"
	"github.com/antonioo83/shot-url-service/internal/server"
	"github.com/antonioo83/shot-url-service/internal/services"
	"github.com/go-chi/jwtauth"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)
import _ "net/http/pprof"

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	fmt.Printf("Build version:%s\n", buildVersion)
	fmt.Printf("Build date:%s\n", buildDate)
	fmt.Printf("Build commit:%s\n", buildCommit)

	configFromFile, err := config.LoadConfigFile("config.json")
	if err != nil {
		fmt.Println("i can't load configuration file:" + err.Error())
	}
	config := config.GetConfigSettings(configFromFile)

	var tokenAuth *jwtauth.JWTAuth
	var pool *pgxpool.Pool
	ctx := context.Background()
	databaseRepository := factory.GetDatabaseRepository(config)
	if config.IsUseDatabase {
		pool, _ = pgxpool.Connect(ctx, config.DatabaseDsn)
		defer pool.Close()
		err := databaseInit(databaseRepository, pool, config.FilepathToDBDump)
		if err != nil {
			log.Fatal(err)
		}
	}
	userRepository := factory.GetUserRepository(ctx, pool, config)
	routeParameters :=
		server.RouteParameters{
			Config:             config,
			ShotURLRepository:  factory.GetRepository(ctx, pool, config),
			UserRepository:     userRepository,
			DatabaseRepository: databaseRepository,
			UserAuthHandler:    authFactory.NewAuthHandler(tokenAuth, userRepository, config),
			Generator:          generators.NewShortLinkDefaultGenerator(),
		}
	handler := server.GetRouters(routeParameters)

	var srv = http.Server{Addr: config.ServerAddress, Handler: handler}
	if config.EnableHTTPS {
		c := services.NewServerCertificate509Service(1658, "Yandex.Praktikum", "RU")
		if err := c.SaveCertificateAndPrivateKeyToFiles("cert.pem", "private.key"); err != nil {
			log.Fatalf("I can't save certificate and private key to files: %v", err)
		}
		if err := srv.ListenAndServeTLS("cert.pem", "private.key"); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	} else {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}

	// через этот канал сообщим основному потоку, что соединения закрыты
	idleConnsClosed := make(chan struct{})
	// канал для перенаправления прерываний
	sigint := make(chan os.Signal, 1)
	shutdownGracefully(ctx, &srv, idleConnsClosed, sigint)

	// ждём завершения процедуры graceful shutdown.
	<-idleConnsClosed
	// получили оповещение о завершении, освобождаем ресурсы перед выходом.
	fmt.Println("Server Shutdown gracefully")
	srv.Shutdown(ctx)
}

func shutdownGracefully(ctx context.Context, srv *http.Server, idleConnsClosed chan struct{}, sigint chan os.Signal) {
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-sigint
		// получили сигнал os.Interrupt, запускаем процедуру graceful shutdown
		if err := srv.Shutdown(ctx); err != nil {
			// ошибки закрытия Listener
			log.Printf("HTTP server Shutdown: %v", err)
		}
		// сообщаем основному потоку, что все сетевые соединения обработаны и закрыты
		close(idleConnsClosed)
	}()

}

func databaseInit(repository interfaces.DatabaseRepository, connect *pgxpool.Pool, filepathToDBDump string) error {
	return repository.RunDump(context.Background(), connect, filepathToDBDump)
}
