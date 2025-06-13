package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// PythonBridge provides an interface to Python scripts for ML model inference
type PythonBridge struct {
	initialized      bool
	modelPath        string
	pythonExecutable string
	initMutex        sync.Mutex
	scriptDir        string
}

// PredictionResult represents the output from Python prediction model
type PredictionResult struct {
	PredictedPrice float64   `json:"predicted_price"`
	Confidence     float64   `json:"confidence"`
	Direction      string    `json:"direction"`
	Factors        []string  `json:"factors"`
	Error          string    `json:"error,omitempty"`
}

var defaultBridge *PythonBridge

// GetPythonBridge returns the shared PythonBridge instance
func GetPythonBridge() *PythonBridge {
	if defaultBridge == nil {
		defaultBridge = NewPythonBridge()
	}
	return defaultBridge
}

// NewPythonBridge creates a new Python bridge
func NewPythonBridge() *PythonBridge {
	return &PythonBridge{
		initialized:      false,
		modelPath:        "",
		pythonExecutable: detectPythonExecutable(),
		scriptDir:        detectScriptDirectory(),
	}
}

// Initialize initializes the Python bridge
func (pb *PythonBridge) Initialize() error {
	pb.initMutex.Lock()
	defer pb.initMutex.Unlock()

	if pb.initialized {
		return nil
	}

	log.Printf("Initializing Python bridge with executable: %s", pb.pythonExecutable)

	// Check if Python executable exists
	if pb.pythonExecutable == "" {
		return fmt.Errorf("Python executable not found")
	}

	// Check if the script directory exists
	if _, err := os.Stat(pb.scriptDir); os.IsNotExist(err) {
		return fmt.Errorf("script directory not found: %s", pb.scriptDir)
	}

	// Check for requirements.txt and ensure dependencies are installed
	requirementsPath := filepath.Join(pb.scriptDir, "requirements.txt")
	if _, err := os.Stat(requirementsPath); err == nil {
		log.Println("Installing Python dependencies...")
		cmd := exec.Command(pb.pythonExecutable, "-m", "pip", "install", "-r", requirementsPath)
		cmd.Dir = pb.scriptDir
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Failed to install dependencies: %v\nOutput: %s", err, string(output))
			return fmt.Errorf("failed to install Python dependencies: %v", err)
		}
		log.Println("Python dependencies installed successfully")
	}

	pb.initialized = true
	return nil
}

// PredictStockPrice uses the Python model to predict stock prices
func (pb *PythonBridge) PredictStockPrice(stock *StockInfo) (*PredictionResult, error) {
	if !pb.initialized {
		if err := pb.Initialize(); err != nil {
			return nil, err
		}
	}

	// Prepare the stock data for the Python model
	stockData := map[string]interface{}{
		"ticker":       stock.Ticker,
		"price":        stock.Price,
		"open":         stock.Open,
		"high":         stock.High,
		"low":          stock.Low,
		"change":       stock.Change,
		"change_pct":   strings.TrimSuffix(stock.ChangePct, "%"),
		"volume":       stock.Volume,
		"market_cap":   stock.MarketCap,
		"company_name": stock.CompanyName,
	}

	// Convert stock data to JSON
	stockJSON, err := json.Marshal(stockData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal stock data: %v", err)
	}

	// Execute the Python prediction script
	scriptPath := filepath.Join(pb.scriptDir, "predict.py")
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		// Generate the prediction script if it doesn't exist
		if err := pb.generatePredictionScript(); err != nil {
			return nil, err
		}
	}

	// Run the prediction script
	cmd := exec.Command(pb.pythonExecutable, scriptPath)
	cmd.Dir = pb.scriptDir
	cmd.Stdin = bytes.NewReader(stockJSON)
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err = cmd.Run()
	if err != nil {
		log.Printf("Python script error: %v\nStderr: %s", err, errOut.String())
		return nil, fmt.Errorf("failed to run prediction script: %v", err)
	}

	// Parse the result
	var result PredictionResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		log.Printf("Failed to parse prediction result: %v\nOutput: %s", err, out.String())
		return nil, fmt.Errorf("failed to parse prediction result: %v", err)
	}

	// Check for errors from Python script
	if result.Error != "" {
		return nil, fmt.Errorf("python script error: %s", result.Error)
	}

	return &result, nil
}

// Helper function to detect Python executable
func detectPythonExecutable() string {
	// Common Python executable names
	pythonCandidates := []string{
		"python3", "python", "python3.10", "python3.9", "python3.8", "python3.7",
	}

	for _, candidate := range pythonCandidates {
		cmd := exec.Command(candidate, "--version")
		if err := cmd.Run(); err == nil {
			log.Printf("Found Python executable: %s", candidate)
			return candidate
		}
	}

	// Path to Python in common virtual environment locations
	venvPaths := []string{
		"venv/bin/python",
		".venv/bin/python",
		"env/bin/python",
	}

	for _, path := range venvPaths {
		if _, err := os.Stat(path); err == nil {
			absPath, err := filepath.Abs(path)
			if err == nil {
				log.Printf("Found Python executable in virtual environment: %s", absPath)
				return absPath
			}
		}
	}

	log.Printf("Warning: Python executable not found")
	return ""
}

// Helper function to detect the script directory
func detectScriptDirectory() string {
	// Try to find the models directory relative to the executable
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		candidates := []string{
			filepath.Join(exeDir, "models"),
			filepath.Join(exeDir, "..", "models"),
			filepath.Join(exeDir, "..", "..", "models"),
		}

		for _, candidate := range candidates {
			if _, err := os.Stat(candidate); err == nil {
				return candidate
			}
		}
	}

	// Try to find the models directory relative to the current working directory
	workingDir, err := os.Getwd()
	if err == nil {
		candidates := []string{
			filepath.Join(workingDir, "models"),
			filepath.Join(workingDir, "..", "models"),
		}

		for _, candidate := range candidates {
			if _, err := os.Stat(candidate); err == nil {
				return candidate
			}
		}
	}

	// Default to the models directory in the project root
	return "./models"
}

// generatePredictionScript creates the Python prediction script
func (pb *PythonBridge) generatePredictionScript() error {
	scriptContent := `#!/usr/bin/env python3
import sys
import json
import numpy as np
import joblib
import os
import tensorflow as tf
from datetime import datetime

def load_model(ticker):
    """Load the trained model for a ticker."""
    model_path = f"saved/{ticker}_model.h5"
    scaler_path = f"saved/{ticker}_scaler.pkl"
    
    # Check if model exists for this ticker
    if not os.path.exists(model_path):
        # Try a default model
        model_path = "saved/DEFAULT_model.h5"
        scaler_path = "saved/DEFAULT_scaler.pkl"
        
        if not os.path.exists(model_path):
            return None, None
    
    try:
        model = tf.keras.models.load_model(model_path)
        scaler = joblib.load(scaler_path)
        return model, scaler
    except Exception as e:
        print(json.dumps({
            "error": f"Failed to load model: {str(e)}"
        }))
        return None, None

def predict(stock_data):
    """Make a prediction for the stock."""
    ticker = stock_data.get("ticker", "DEFAULT")
    
    # Try to load the model for this ticker
    model, scaler = load_model(ticker)
    
    # If no model exists, use a simple rule-based prediction
    if model is None:
        return rule_based_prediction(stock_data)
        
    try:
        # Extract features for prediction
        price = float(stock_data.get("price", 0))
        open_price = float(stock_data.get("open", price))
        high = float(stock_data.get("high", price * 1.01))
        low = float(stock_data.get("low", price * 0.99))
        
        # Create a feature vector (this needs to match your model's expected input)
        # This is a simplified example - your actual model may need different features
        features = np.array([[price, open_price, high, low]])
        
        # Scale the features
        scaled_features = scaler.transform(features)
        
        # Make prediction
        prediction = model.predict(scaled_features)
        
        # Convert prediction back to price
        predicted_price = float(scaler.inverse_transform(prediction)[0][0])
        
        # Calculate other metrics
        price_diff = ((predicted_price - price) / price) * 100
        confidence = min(0.95, max(0.55, 0.85 - abs(price_diff) * 0.01))
        
        direction = "NEUTRAL"
        if price_diff > 1.0:
            direction = "UP"
        elif price_diff < -1.0:
            direction = "DOWN"
            
        # Generate factors
        factors = generate_factors(stock_data, price_diff)
        
        return {
            "predicted_price": round(predicted_price, 2),
            "confidence": confidence,
            "direction": direction,
            "factors": factors
        }
        
    except Exception as e:
        # Fallback to rule-based prediction on error
        print(f"Model prediction error: {str(e)}", file=sys.stderr)
        return rule_based_prediction(stock_data)

def rule_based_prediction(stock_data):
    """Simple rule-based prediction as fallback."""
    price = float(stock_data.get("price", 0))
    if price <= 0:
        return {
            "error": "Invalid price data"
        }
        
    # Parse change percentage
    change_pct = stock_data.get("change_pct", "0")
    if isinstance(change_pct, str) and "%" in change_pct:
        change_pct = change_pct.replace("%", "")
    change_pct = float(change_pct)
    
    # Simple momentum-based prediction
    momentum = change_pct * 0.1
    
    # Add some randomness
    import random
    random_factor = (random.random() - 0.5) * 0.5
    
    # Calculate predicted change percentage
    predicted_change_pct = momentum + random_factor
    
    # Calculate predicted price
    predicted_price = price * (1 + predicted_change_pct / 100)
    predicted_price = round(predicted_price, 2)
    
    # Determine direction
    direction = "NEUTRAL"
    if predicted_change_pct > 1.0:
        direction = "UP"
    elif predicted_change_pct < -1.0:
        direction = "DOWN"
    
    # Calculate confidence (higher for smaller predictions)
    confidence = 0.7 - min(0.2, abs(predicted_change_pct) * 0.02)
    
    # Generate factors
    factors = generate_factors(stock_data, predicted_change_pct)
    
    return {
        "predicted_price": predicted_price,
        "confidence": confidence,
        "direction": direction,
        "factors": factors
    }

def generate_factors(stock_data, predicted_change_pct):
    """Generate key factors that influence the prediction."""
    factors = []
    
    # Price momentum
    price = float(stock_data.get("price", 0))
    open_price = float(stock_data.get("open", price))
    
    if price > open_price:
        factors.append("Price is above opening level")
    elif price < open_price:
        factors.append("Price is below opening level")
    
    # Market cap based factor
    market_cap = stock_data.get("market_cap", "")
    if "T" in market_cap:
        factors.append("Large market cap indicates stability")
    elif "B" in market_cap and not "T" in market_cap:
        if float(market_cap.replace("$", "").replace("B", "")) > 10:
            factors.append("Medium-large market cap suggests moderate volatility")
        else:
            factors.append("Medium market cap may lead to higher volatility")
    else:
        factors.append("Smaller capitalization suggests higher volatility")
    
    # Add prediction direction factor
    if predicted_change_pct > 0:
        factors.append("Technical indicators suggest positive momentum")
    else:
        factors.append("Technical indicators suggest negative pressure")
    
    return factors[:3]  # Return top 3 factors

def main():
    """Main entry point for the prediction script."""
    try:
        # Read stock data from stdin
        stock_data = json.load(sys.stdin)
        
        # Make prediction
        result = predict(stock_data)
        
        # Output result as JSON
        print(json.dumps(result))
        
    except Exception as e:
        print(json.dumps({
            "error": str(e)
        }))

if __name__ == "__main__":
    main()
`

	scriptPath := filepath.Join(pb.scriptDir, "predict.py")
	
	// Create script file with executable permissions
	err := os.WriteFile(scriptPath, []byte(scriptContent), 0755)
	if err != nil {
		return fmt.Errorf("failed to create prediction script: %v", err)
	}
	
	log.Printf("Generated prediction script at %s", scriptPath)
	return nil
}
