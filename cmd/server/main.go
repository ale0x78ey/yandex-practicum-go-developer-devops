package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	host = "0.0.0.0"
	port = "80"
)

type serverHandler struct {
}

func (s *serverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		// TODO: read all from r?
		return
	}
	log.Printf("serverHandler.ServeHTTP %v", r.URL)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func main() {
	serverHandler := &serverHandler{}
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: serverHandler,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
