package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

	// Apply CORS middleware to all routes
	r.Use(handlers.CorsMiddleware)
	
	// API routes for stock data
	r.HandleFunc("/api/health", handlers.APIHealthHandler).Methods("GET")
	
	// Initialize WebSocket handler and set up route for real-time stock updates
	handlers.InitWebSocketHandler()
	r.HandleFunc("/ws/stocks/{ticker}", handlers.HandleWebSocket)
	
	// API routes for polling fallback
	r.HandleFunc("/api/stocks/{ticker}", handlers.StockAPIHandler).Methods("GET")
	
	// Start the WebSocket broadcaster
	handlers.StartPriceUpdateBroadcaster()
	
	// Serve React app static files from frontend/build
	frontendBuildPath := filepath.Join("..", "frontend", "build")
	staticFileServer := http.FileServer(http.Dir(frontendBuildPath))
	
	// Special handling for the root path to serve index.html
	r.PathPrefix("/static/").Handler(http.StripPrefix("/", staticFileServer))
	r.PathPrefix("/manifest.json").Handler(staticFileServer)
	r.PathPrefix("/favicon.ico").Handler(staticFileServer)
	r.PathPrefix("/logo").Handler(staticFileServer)
	
	// Catch-all handler for React app - serves index.html for any unmatched routes
	// This enables client-side routing in React
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip API and WebSocket paths
		if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/ws/") {
			http.NotFound(w, r)
			return
		}
		
		indexPath := filepath.Join(frontendBuildPath, "index.html")
		http.ServeFile(w, r, indexPath)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"  // Changed to port 8084 to avoid conflict
	}

	log.Printf("Starting server on port %s (Python bridge: %v)", port, pythonAvailable)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
