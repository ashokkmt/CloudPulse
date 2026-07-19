package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"cloudpulse/backend/internal/handler"
	"cloudpulse/backend/internal/service"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	store := service.NewStore()
	router := handler.NewRouter(store)
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("CloudPulse backend listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
