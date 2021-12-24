package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ale0x78ey/yandex-practicum-go-developer-devops/api/rest"
)

const (
	host = "0.0.0.0"
	port = "80"
)

func main() {
	restServer := rest.NewServer()
	if restServer == nil {
		log.Fatal("Server failed: REST Server wasn't created")
	}

	router := restServer.NewRouter()
	if router == nil {
		log.Fatal("Server failed: Router wasn't created")
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: router,
	}

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
