package services

import (
	"testing"
)

// TestTensorFlowPrediction tests the TensorFlow model prediction functionality
func TestTensorFlowPrediction(t *testing.T) {
	// Create the TensorFlow model service
	tfService := NewTFModelService()
	
	// Create a sample stock for prediction
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
	
	// Get prediction
	prediction, err := tfService.PredictStockMovement(stockInfo)
	if err != nil {
		t.Errorf("Failed to get prediction: %v", err)
	}
	
	// Verify prediction is reasonable
	if prediction.PredictedPrice <= 0 {
		t.Errorf("Expected positive predicted price, got %.2f", prediction.PredictedPrice)
	}
	
	if prediction.Confidence <= 0 || prediction.Confidence > 1.0 {
		t.Errorf("Confidence should be between 0 and 1, got %.2f", prediction.Confidence)
	}
	
	if prediction.Direction != "UP" && prediction.Direction != "DOWN" && prediction.Direction != "NEUTRAL" {
		t.Errorf("Direction should be UP, DOWN, or NEUTRAL, got %s", prediction.Direction)
	}
	
	// Make sure price change is within reasonable bounds (test max 10% limit)
	maxDiff := stockInfo.Price * 0.10
	actualDiff := prediction.PredictedPrice - stockInfo.Price
	if actualDiff > maxDiff {
		t.Errorf("Price change exceeds 10%% maximum allowed: %.2f", actualDiff)
	}
}

// TestErrorHandling tests error cases
func TestTFErrorHandling(t *testing.T) {
	tfService := NewTFModelService()
	
	// Test with nil stock
	_, err := tfService.PredictStockMovement(nil)
	if err == nil {
		t.Errorf("Expected error for nil stock data")
	}
	
	// Test with invalid price
	invalidStock := &StockInfo{
		Ticker:      "AAPL",
		CompanyName: "Apple Inc.",
		Price:       -10.0, // Invalid negative price
		Change:      2.5,
		Open:        149.0,
		High:        152.0,
		Low:         148.5,
	}
	
	_, err = tfService.PredictStockMovement(invalidStock)
	if err == nil {
		t.Errorf("Expected error for negative stock price")
	}
	
	// Test with zero OHLC data - should auto-repair where possible
	zeroDataStock := &StockInfo{
		Ticker:      "AAPL",
		CompanyName: "Apple Inc.",
		Price:       150.25,
		Change:      2.5,
		Open:        0, // Zero values
		High:        0,
		Low:         0,
	}
	
	prediction, err := tfService.PredictStockMovement(zeroDataStock)
	if err != nil {
		t.Errorf("Expected auto-repair of zero OHLC values, got error: %v", err)
	}
	
	if prediction == nil {
		t.Errorf("Expected valid prediction after auto-repair of data")
	}
}

// TestFeatureExtraction tests the feature extraction functionality
func TestFeatureExtraction(t *testing.T) {
	tfService := NewTFModelService()
	
	stockInfo := &StockInfo{
		Price: 100.0,
		High:  110.0,
		Low:   90.0,
		Open:  95.0,
		ChangePct: "5.0%",
	}
	
	features := tfService.extractFeatures(stockInfo)
	
	// Check that features map is populated
	if features == nil {
		t.Errorf("Expected non-nil features map")
	}
	
	// Check that features are calculated
	if len(features) == 0 {
		t.Errorf("Expected non-empty features map")
	}
}
