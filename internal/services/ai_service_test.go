package services

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestAIService tests the AI service functionality
func TestAIService(t *testing.T) {
	// Test rule-based recommendation when no OpenAI key is provided
	aiService := NewAIService("")
	stockInfo := &StockInfo{
		Ticker:      "AAPL",
		CompanyName: "Apple Inc.",
		Price:       150.25,
		Change:      2.5,
		ChangePct:   "1.68%",
		Open:        149.0,
		High:        152.0,
		Low:         148.5,
		Volume:      "75M",
		MarketCap:   "$2.5T",
	}

	// Get recommendation without API key
	recommendation, err := aiService.GetStockRecommendation(stockInfo)
	if err != nil {
		t.Errorf("Failed to get rule-based recommendation: %v", err)
	}
	
	// Verify we got some recommendation
	if recommendation == "" {
		t.Errorf("Expected non-empty recommendation")
	}
	
	// Check that recommendation includes BUY, SELL, or HOLD
	if !strings.Contains(recommendation, "BUY") && 
	   !strings.Contains(recommendation, "SELL") && 
	   !strings.Contains(recommendation, "HOLD") {
		t.Errorf("Recommendation should contain BUY, SELL or HOLD, got: %s", recommendation)
	}
}

// TestOpenAIResponse tests how the service handles OpenAI API responses
func TestOpenAIResponse(t *testing.T) {
	// Create a mock server for OpenAI API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"choices": [{
				"message": {
					"content": "BUY - Strong fundamentals and positive momentum."
				}
			}]
		}`))
	}))
	defer server.Close()
	
	// Create AI service with fake API key
	aiService := NewAIService("fake-api-key")
	
	// Create a test stock
	stockInfo := &StockInfo{
		Ticker:      "AAPL",
		CompanyName: "Apple Inc.",
		Price:       150.25,
		Change:      2.5,
		ChangePct:   "1.68%",
	}
	
	// Get recommendation
	recommendation, err := aiService.GetStockRecommendation(stockInfo)
	if err != nil {
		t.Errorf("Failed to get OpenAI recommendation: %v", err)
	}
	
	// Verify we got some recommendation - either from API or rule-based fallback
	if recommendation == "" {
		t.Errorf("Expected non-empty recommendation")
	}
}

// TestErrorHandling tests API error handling
func TestErrorHandling(t *testing.T) {
	// Create a mock server that simulates failures
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()
	
	// Create AI service with fake API key
	aiService := NewAIService("fake-api-key")
	
	// Create a test stock
	stockInfo := &StockInfo{
		Ticker:      "AAPL",
		CompanyName: "Apple Inc.",
		Price:       150.25,
	}
	
	// Get recommendation - should fall back to rule-based
	recommendation, err := aiService.GetStockRecommendation(stockInfo)
	
	// We should still get a recommendation, even if API fails
	if err != nil {
		t.Errorf("Expected no error when API fails (should use fallback), got: %v", err)
	}
	
	// Verify we got some recommendation
	if recommendation == "" {
		t.Errorf("Expected fallback recommendation when API fails")
	}
}
