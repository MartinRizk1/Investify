package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// StockAPIHandler handles requests for stock data via API endpoint
func StockAPIHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ticker := vars["ticker"]
	
	if ticker == "" {
		http.Error(w, "Ticker symbol required", http.StatusBadRequest)
		return
	}
	
	// Get stock price data
	stockPrice, err := fetchRealtimePrice(ticker)
	if err != nil {
		log.Printf("Error fetching price for %s: %v", ticker, err)
		http.Error(w, "Error fetching stock data", http.StatusInternalServerError)
		return
	}

	// Get technical indicators
	technical, err := fetchTechnicalIndicators(ticker)
	if err != nil {
		log.Printf("Error fetching technical indicators for %s: %v", ticker, err)
		// Continue anyway, just without technical data
	}

	update := StockUpdate{
		Ticker:      ticker,
		Price:       stockPrice.Price,
		Change:      stockPrice.Change,
		ChangePct:   stockPrice.ChangePct,
		LastUpdated: formatTime(nil), // Current time formatted
		Technical:   technical,
	}

	// Set content type and send response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(update); err != nil {
		log.Printf("Error encoding stock update: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
