package main

import (
	"github.com/antonioo83/shot-url-service/internal/server"
	"log"
	"net/http"
)

func main() {
	log.Fatal(http.ListenAndServe(":8080", server.GetRouters()))
}
