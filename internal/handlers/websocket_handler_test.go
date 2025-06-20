package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// TestWebSocketConnection tests the WebSocket connection and initial data
func TestWebSocketConnection(t *testing.T) {
	// Initialize the handler
	InitWebSocketHandler()

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/ws/stocks/") {
			HandleWebSocket(w, r)
		} else {
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	// Convert http to ws URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/stocks/AAPL"

	// Connect to WebSocket
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer ws.Close()

	// Wait for and read the initial message
	ws.SetReadDeadline(time.Now().Add(5 * time.Second))
	_, message, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read WebSocket message: %v", err)
	}

	// Parse the JSON message
	var update StockUpdate
	if err := json.Unmarshal(message, &update); err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}
	
	// Verify the update contains technical indicators data
	if update.Technical == nil {
		t.Error("Technical indicators data is missing from WebSocket update")
	} else {
		// Check for required technical indicators
		requiredFields := []string{"dates", "rsi", "macd", "macd_signal", "macd_histogram", "bollinger_middle", "bollinger_upper", "bollinger_lower"}
		for _, field := range requiredFields {
			if _, exists := update.Technical[field]; !exists {
				t.Errorf("Required technical indicator '%s' is missing", field)
			}
		}
	}

	// Verify the message contents
	if update.Ticker != "AAPL" {
		t.Errorf("Expected ticker AAPL, got %s", update.Ticker)
	}
	if update.Price <= 0 {
		t.Errorf("Expected positive price, got %f", update.Price)
	}
	if update.LastUpdated == "" {
		t.Error("Expected LastUpdated timestamp to be set")
	}

	// Check for technical indicators
	if update.Technical == nil {
		t.Error("Expected technical indicators to be present")
	} else {
		// Verify key technical indicators
		if rsi, ok := update.Technical["rsi"]; !ok {
			t.Error("RSI data missing from technical indicators")
		} else if rsiSlice, isSlice := rsi.([]interface{}); !isSlice || len(rsiSlice) == 0 {
			t.Error("RSI data should be a non-empty slice")
		}

		if macd, ok := update.Technical["macd"]; !ok {
			t.Error("MACD data missing from technical indicators")
		} else if macdSlice, isSlice := macd.([]interface{}); !isSlice || len(macdSlice) == 0 {
			t.Error("MACD data should be a non-empty slice")
		}

		if bbMiddle, ok := update.Technical["bollinger_middle"]; !ok {
			t.Error("Bollinger middle band data missing from technical indicators")
		} else if bbSlice, isSlice := bbMiddle.([]interface{}); !isSlice || len(bbSlice) == 0 {
			t.Error("Bollinger band data should be a non-empty slice")
		}
	}
}

// TestWebSocketBroadcaster tests the broadcaster functionality
func TestWebSocketBroadcaster(t *testing.T) {
	// This is a simple test to verify the broadcaster starts without panicking
	// For a real test, we'd need to wait for broadcasts and check them
	
	// Initialize the handler
	InitWebSocketHandler()
	
	// Start the broadcaster and wait a moment
	done := make(chan bool)
	go func() {
		StartPriceUpdateBroadcaster()
		time.Sleep(2 * time.Second)
		done <- true
	}()
	
	select {
	case <-done:
		// Broadcaster started successfully
	case <-time.After(5 * time.Second):
		t.Error("Timeout waiting for broadcaster")
	}
}
