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
	grpc2 "github.com/antonioo83/shot-url-service/internal/server/grpc"
	pb "github.com/antonioo83/shot-url-service/internal/server/grpc/proto"
	"github.com/antonioo83/shot-url-service/internal/services"
	"github.com/go-chi/jwtauth"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"log"
	"net"
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
	cfg := config.GetConfigSettings(configFromFile)

	var tokenAuth *jwtauth.JWTAuth
	var pool *pgxpool.Pool
	ctx := context.Background()
	databaseRepository := factory.GetDatabaseRepository(cfg)
	if cfg.IsUseDatabase {
		pool, _ = pgxpool.Connect(ctx, cfg.DatabaseDsn)
		defer pool.Close()
		err := databaseInit(databaseRepository, pool, cfg.FilepathToDBDump)
		if err != nil {
			log.Fatal(err)
		}
	}
	userRepository := factory.GetUserRepository(ctx, pool, cfg)
	routeParameters :=
		server.RouteParameters{
			Config:             cfg,
			ShotURLRepository:  factory.GetRepository(ctx, pool, cfg),
			UserRepository:     userRepository,
			DatabaseRepository: databaseRepository,
			UserAuthHandler:    authFactory.NewAuthHandler(tokenAuth, userRepository, cfg),
			Generator:          generators.NewShortLinkDefaultGenerator(),
		}
	handler := server.GetRouters(routeParameters)

	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	if cfg.ServerType == config.HTTPServer {
		var srv = http.Server{Addr: cfg.ServerAddress, Handler: handler}
		runHTTPServer(cfg, srv)
		shutdownGracefullyHTTPServer(ctx, &srv, idleConnsClosed, sigint)
		<-idleConnsClosed
		fmt.Println("Server HTTP Shutdown gracefully")
		srv.Shutdown(ctx)
	} else if cfg.ServerType == config.GRPCServer {
		srv := grpc.NewServer()
		runGRPCServer(cfg, srv, routeParameters)
		shutdownGracefullyGRPCServer(srv, idleConnsClosed, sigint)
		<-idleConnsClosed
		fmt.Println("Server GRPC Shutdown gracefully")
		srv.GracefulStop()
	} else {
		log.Fatalf("Unknowned server type")
	}
}

func runHTTPServer(config config.Config, srv http.Server) {
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
}

func runGRPCServer(cfg config.Config, srv *grpc.Server, p server.RouteParameters) {
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		log.Fatal(err)
	}

	var s grpc2.ShortURLServer
	s.Config = p.Config
	s.ShotURLRepository = p.ShotURLRepository
	s.UserRepository = p.UserRepository
	s.DatabaseRepository = p.DatabaseRepository
	s.UserAuthHandler = p.UserAuthHandler
	s.Generator = p.Generator
	pb.RegisterShortURLServer(srv, &s)

	fmt.Println("Сервер gRPC начал работу")
	if err := srv.Serve(listen); err != nil {
		log.Fatal(err)
	}
}

func shutdownGracefullyHTTPServer(ctx context.Context, srv *http.Server, idleConnsClosed chan struct{}, sigint chan os.Signal) {
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

func shutdownGracefullyGRPCServer(srv *grpc.Server, idleConnsClosed chan struct{}, sigint chan os.Signal) {
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		<-sigint
		// получили сигнал os.Interrupt, запускаем процедуру graceful shutdown
		srv.GracefulStop()
		// сообщаем основному потоку, что все сетевые соединения обработаны и закрыты
		close(idleConnsClosed)
	}()
}

func databaseInit(repository interfaces.DatabaseRepository, connect *pgxpool.Pool, filepathToDBDump string) error {
	return repository.RunDump(context.Background(), connect, filepathToDBDump)
}
