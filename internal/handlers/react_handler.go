package handlers

import (
	"encoding/json"
	"net/http"
)

// ReactAppHandler serves the React frontend application
func ReactAppHandler(w http.ResponseWriter, r *http.Request) {
	// This is a simple redirect handler for React app
	// The actual serving of React files is done in main.go
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// APIHealthHandler returns API health status
func APIHealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	
	response := map[string]interface{}{
		"status": "healthy",
		"version": "1.0.0",
		"features": map[string]bool{
			"ai_predictions": true,
			"real_time_data": true,
			"technical_indicators": true,
		},
	}
	
	json.NewEncoder(w).Encode(response)
}
