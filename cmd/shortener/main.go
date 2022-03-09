package main

import (
	"github.com/antonioo83/shot-url-service/config"
	"github.com/antonioo83/shot-url-service/internal/server"
	"log"
	"net/http"
)

func main() {
	configSettings := config.GetConfigSettings()
	server.LoadModelsFromDatabase(configSettings)
	log.Fatal(http.ListenAndServe(configSettings.ServerAddress, server.GetRouters(configSettings)))
}
