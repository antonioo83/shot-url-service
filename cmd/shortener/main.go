package main

import (
	"github.com/antonioo83/shot-url-service/internal/server"
	"log"
	http "net/http"
)

func main() {
	log.Fatal(http.ListenAndServe(":8080", server.GetRouters()))
}
