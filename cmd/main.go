package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/martinrizk/investify/internal/handlers"
)

func main() {
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/", handlers.HomeHandler).Methods("GET", "POST")
	r.HandleFunc("/health", handlers.HealthHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
