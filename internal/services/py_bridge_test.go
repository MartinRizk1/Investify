package services

import (
	"os"
	"path/filepath"
	"testing"
)

// TestPythonBridgeInitialization tests that the Python bridge can be initialized
func TestPythonBridgeInitialization(t *testing.T) {
	bridge := NewPythonBridge()
	err := bridge.Initialize()

	// If Python is not installed, we'll skip the test
	if bridge.pythonExecutable == "" {
		t.Skip("Python executable not found, skipping test")
	}

	if err != nil {
		t.Errorf("Failed to initialize Python bridge: %v", err)
	}

	if !bridge.initialized {
		t.Error("Python bridge not marked as initialized after initialization")
	}
}

// TestGeneratePredictionScript tests that the prediction script can be generated
func TestGeneratePredictionScript(t *testing.T) {
	bridge := NewPythonBridge()

	// Set a test script directory
	testDir := filepath.Join(os.TempDir(), "investify_test")
	defer os.RemoveAll(testDir) // Clean up after test

	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	bridge.scriptDir = testDir
	err = bridge.generatePredictionScript()
	if err != nil {
		t.Errorf("Failed to generate prediction script: %v", err)
	}

	// Check if the script was created
	scriptPath := filepath.Join(testDir, "predict.py")
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		t.Error("Prediction script was not created")
	}
}

// TestPredictStockPrice tests the stock price prediction functionality
func TestPredictStockPrice(t *testing.T) {
	bridge := NewPythonBridge()

	// Skip if Python is not available
	if bridge.pythonExecutable == "" {
		t.Skip("Python executable not found, skipping test")
	}

	// Initialize the bridge
	err := bridge.Initialize()
	if err != nil {
		t.Fatalf("Failed to initialize Python bridge: %v", err)
	}

	// Create a test stock info
	stockInfo := &StockInfo{
		Ticker:      "AAPL",
		CompanyName: "Apple Inc.",
		Price:       150.0,
		Open:        149.0,
		High:        152.0,
		Low:         148.0,
		Change:      1.0,
		ChangePct:   "0.67%",
		Volume:      "10M",
		MarketCap:   "$2.5T",
	}

	// Make a prediction
	result, err := bridge.PredictStockPrice(stockInfo)
	
	// We need to handle the case where the model doesn't exist yet
	// This test doesn't assert the accuracy of predictions, just that the mechanism works
	if err != nil {
		if err.Error() == "python script error: Failed to load model: [Errno 2] No such file or directory: 'saved/AAPL_model.h5'" {
			t.Skip("Model not trained, skipping prediction test")
		} else {
			t.Errorf("Failed to predict stock price: %v", err)
		}
	}

	// Verify prediction result structure
	if result != nil {
		if result.PredictedPrice <= 0 {
			t.Error("Invalid predicted price")
		}

		if result.Confidence <= 0 || result.Confidence > 1 {
			t.Errorf("Invalid confidence level: %f", result.Confidence)
		}

		if result.Direction != "UP" && result.Direction != "DOWN" && result.Direction != "NEUTRAL" {
			t.Errorf("Invalid direction: %s", result.Direction)
		}

		if len(result.Factors) == 0 {
			t.Error("No factors provided in prediction result")
		}
	}
}
