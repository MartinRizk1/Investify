package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/martinrizk/investify/internal/handlers"
	"github.com/martinrizk/investify/internal/services"
)

func main() {
	// Check Python bridge availability
	bridge := services.GetPythonBridge()
	pythonAvailable := false
	
	if err := bridge.Initialize(); err == nil {
		pythonAvailable = true
		log.Println("Python bridge initialized successfully - TensorFlow models will be used if available")
	} else {
		log.Printf("Python bridge initialization failed: %v", err)
		log.Println("Using Go fallback for predictions - this is expected if Python/TensorFlow is not installed")
	}

	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/", handlers.HomeHandler).Methods("GET", "POST")
	r.HandleFunc("/health", handlers.HealthHandler).Methods("GET")
	
	// Static file routes
	r.HandleFunc("/scripts.js", handlers.StaticFileHandler)
	r.HandleFunc("/styles.css", handlers.StaticFileHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s (Python bridge: %v)", port, pythonAvailable)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
