package main

import (
	"github.com/antonioo83/shot-url-service/internal/server"
	"log"
	http "net/http"
)

func main() {
	// маршрутизация запросов обработчику
	http.HandleFunc("/", server.UrlHandler)
	// запуск сервера с адресом localhost, порт 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}
