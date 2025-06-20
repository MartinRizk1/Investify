package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// TestStockAPIHandler tests the API endpoint for fetching stock data
func TestStockAPIHandler(t *testing.T) {
	// Initialize the handler
	InitWebSocketHandler()

	// Create a test router with our handler
	r := mux.NewRouter()
	r.HandleFunc("/api/stocks/{ticker}", StockAPIHandler)

	// Create a test request
	req, err := http.NewRequest("GET", "/api/stocks/AAPL", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Parse the response
	var update StockUpdate
	if err := json.NewDecoder(rr.Body).Decode(&update); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Verify the response contains required fields
	if update.Ticker != "AAPL" {
		t.Errorf("Expected ticker AAPL, got %s", update.Ticker)
	}

	if update.Price <= 0 {
		t.Errorf("Expected positive price, got %f", update.Price)
	}

	if update.LastUpdated == "" {
		t.Error("Expected LastUpdated timestamp to be set")
	}

	// Verify technical indicators
	if update.Technical == nil {
		t.Error("Technical indicators data is missing from API response")
	} else {
		// Check for required technical indicators
		requiredFields := []string{"dates", "rsi", "macd", "macd_signal", "macd_histogram", "bollinger_middle", "bollinger_upper", "bollinger_lower"}
		for _, field := range requiredFields {
			if _, exists := update.Technical[field]; !exists {
				t.Errorf("Required technical indicator '%s' is missing", field)
			}
		}
	}
}
