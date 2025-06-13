package services

import (
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

// TFModelService handles TensorFlow-based stock predictions
type TFModelService struct {
	modelReady bool
}

// StockPrediction represents a prediction made by the TensorFlow model
type StockPrediction struct {
	PredictedPrice float64
	Confidence     float64
	Direction      string // UP, DOWN, NEUTRAL
	Factors        []string
}

// NewTFModelService creates a new TensorFlow model service
func NewTFModelService() *TFModelService {
	// Initialize the Python bridge in the background
	go func() {
		bridge := GetPythonBridge()
		err := bridge.Initialize()
		if err != nil {
			log.Printf("Warning: Failed to initialize Python bridge for TensorFlow: %v", err)
			log.Printf("Using simulated predictions instead of actual TensorFlow model")
		} else {
			log.Printf("Successfully initialized Python bridge for TensorFlow models")
		}
	}()
	
	// We're ready to make predictions, either with the Python bridge or simulated
	return &TFModelService{
		modelReady: true,
	}
}

// PredictStockMovement predicts future stock movements based on current data
func (tf *TFModelService) PredictStockMovement(stock *StockInfo) (*StockPrediction, error) {
	if !tf.modelReady {
		return nil, fmt.Errorf("tensorflow model not initialized")
	}

	// More thorough validation of stock data
	if stock == nil {
		return nil, fmt.Errorf("stock data is nil")
	}
	
	if stock.Price <= 0 {
		return nil, fmt.Errorf("invalid stock price: %.2f", stock.Price)
	}
	
	// Validate other essential fields
	if stock.Open <= 0 || stock.High <= 0 || stock.Low <= 0 {
		log.Printf("Warning: Stock %s has potentially invalid OHLC data: Open=%.2f, High=%.2f, Low=%.2f", 
			stock.Ticker, stock.Open, stock.High, stock.Low)
		// Try to recover with available data
		if stock.Open <= 0 {
			stock.Open = stock.Price
		}
		if stock.High <= 0 {
			stock.High = math.Max(stock.Price, stock.Open) * 1.01 // Add small buffer
		}
		if stock.Low <= 0 {
			stock.Low = math.Min(stock.Price, stock.Open) * 0.99 // Add small buffer
		}
	}

	// Try to use Python bridge first if available
	bridge := GetPythonBridge()
	if bridge != nil {
		result, err := bridge.PredictStockPrice(stock)
		if err == nil && result != nil {
			// Convert Python bridge result to StockPrediction
			return &StockPrediction{
				PredictedPrice: result.PredictedPrice,
				Confidence:     result.Confidence,
				Direction:      result.Direction,
				Factors:        result.Factors,
			}, nil
		}
		// Log the error but continue with fallback prediction
		log.Printf("Python bridge prediction failed for %s: %v. Using fallback prediction.", 
			stock.Ticker, err)
	}

	// Fallback to simulated prediction
	log.Printf("Using simulated prediction for %s", stock.Ticker)
	
	// Extract features from the stock data
	features := tf.extractFeatures(stock)
	
	// Simulate a prediction with error handling
	var prediction *StockPrediction
	// Use defer/recover to handle any potential panic in prediction code
	func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic in stock prediction for %s: %v", stock.Ticker, r)
				prediction = nil
			}
		}()
		prediction = tf.simulatePrediction(stock, features)
	}()

	// Validate prediction
	if prediction == nil || prediction.PredictedPrice <= 0 || prediction.Confidence <= 0 {
		return nil, fmt.Errorf("failed to generate a valid prediction")
	}

	// Ensure predictions stay within realistic bounds (max 10% change)
	maxChange := stock.Price * 0.10
	if math.Abs(prediction.PredictedPrice - stock.Price) > maxChange {
		if prediction.PredictedPrice > stock.Price {
			prediction.PredictedPrice = stock.Price + maxChange
		} else {
			prediction.PredictedPrice = stock.Price - maxChange
		}
		
		// Adjust confidence downward for extreme predictions
		prediction.Confidence *= 0.85
		prediction.Factors = append(prediction.Factors, "Extreme prediction detected and adjusted")
	}

	return prediction, nil
}

// extractFeatures extracts relevant features from stock data
func (tf *TFModelService) extractFeatures(stock *StockInfo) map[string]float64 {
	// Calculate price relative to day's range
	dayRange := stock.High - stock.Low
	if dayRange == 0 {
		dayRange = 0.01 // Avoid division by zero
	}
	pricePosition := (stock.Price - stock.Low) / dayRange

	// Parse change percentage from string
	changePct := 0.0
	if pctStr := strings.TrimSuffix(stock.ChangePct, "%"); pctStr != "" {
		if val, err := strconv.ParseFloat(pctStr, 64); err == nil {
			changePct = val
		}
	}

	// Calculate volatility
	volatility := dayRange / stock.Price * 100
	
	// Market cap numeric value
	marketCapValue := 0.0
	if strings.HasPrefix(stock.MarketCap, "$") {
		mcStr := strings.TrimPrefix(stock.MarketCap, "$")
		if strings.HasSuffix(mcStr, "T") {
			if val, err := strconv.ParseFloat(strings.TrimSuffix(mcStr, "T"), 64); err == nil {
				marketCapValue = val * 1e12
			}
		} else if strings.HasSuffix(mcStr, "B") {
			if val, err := strconv.ParseFloat(strings.TrimSuffix(mcStr, "B"), 64); err == nil {
				marketCapValue = val * 1e9
			}
		} else if strings.HasSuffix(mcStr, "M") {
			if val, err := strconv.ParseFloat(strings.TrimSuffix(mcStr, "M"), 64); err == nil {
				marketCapValue = val * 1e6
			}
		}
	}
	
	return map[string]float64{
		"price_position": pricePosition,
		"change_pct":     changePct,
		"volatility":     volatility,
		"market_cap":     marketCapValue,
		"open_gap":       (stock.Open - stock.Low) / dayRange,
		"high_ratio":     stock.High / stock.Price,
		"low_ratio":      stock.Low / stock.Price,
	}
}

// simulatePrediction simulates a TensorFlow prediction
func (tf *TFModelService) simulatePrediction(stock *StockInfo, features map[string]float64) *StockPrediction {
	// In a real implementation, this would use the TensorFlow model to make predictions
	// For now, we'll use a rule-based system with some randomness to simulate predictions
	
	// Calculate base prediction
	momentumFactor := features["change_pct"] * 0.1
	positionFactor := (features["price_position"] - 0.5) * -0.2 // Mean reversion
	volatilityFactor := features["volatility"] * 0.05
	
	// Market cap factor - larger companies tend to be more stable
	marketCapFactor := 0.0
	if features["market_cap"] > 1e11 { // $100B+
		marketCapFactor = -0.1 // Less volatile
	} else if features["market_cap"] < 1e9 { // Less than $1B
		marketCapFactor = 0.2 // More volatile
	}
	
	// Calculate predicted change percentage
	predictedChangePct := momentumFactor + positionFactor + volatilityFactor + marketCapFactor
	
	// Current day simulation
	intraday := time.Now().Hour() < 16 // Before market close
	if intraday {
		// If it's intraday, factor in the current position in daily range
		if features["price_position"] > 0.8 {
			predictedChangePct -= 0.3 // More likely to revert if already near high
		} else if features["price_position"] < 0.2 {
			predictedChangePct += 0.3 // More likely to bounce if already near low
		}
	}
	
	// Add some controlled randomness
	randomFactor := (math.Sin(float64(time.Now().UnixNano())) * 0.5)
	predictedChangePct += randomFactor
	
	// Calculate predicted price with more precision
	predictedPrice := stock.Price * (1 + predictedChangePct/100)
	
	// Round to 2 decimal places for better display
	predictedPrice = math.Round(predictedPrice*100) / 100
	
	// Determine direction with clearer thresholds
	direction := "NEUTRAL"
	if predictedChangePct > 1.0 {
		direction = "UP"
	} else if predictedChangePct < -1.0 {
		direction = "DOWN"
	}
	
	// Calculate confidence based on consistency of signals
	signals := []float64{momentumFactor, positionFactor, volatilityFactor, marketCapFactor, randomFactor}
	confidence := tf.calculateConfidence(signals, predictedChangePct)
	
	// Identify key factors driving the prediction
	factors := tf.identifyKeyFactors(stock, features, predictedChangePct)
	
	return &StockPrediction{
		PredictedPrice: predictedPrice,
		Confidence:     confidence,
		Direction:      direction,
		Factors:        factors,
	}
}

// calculateConfidence calculates a confidence score for the prediction
func (tf *TFModelService) calculateConfidence(signals []float64, prediction float64) float64 {
	// Count how many signals agree with the prediction direction
	agreementCount := 0
	for _, signal := range signals {
		if (prediction > 0 && signal > 0) || (prediction < 0 && signal < 0) {
			agreementCount++
		}
	}
	
	// Base confidence on agreement percentage
	baseConfidence := float64(agreementCount) / float64(len(signals))
	
	// Adjust for prediction magnitude - very large predictions are less confident
	magnitudePenalty := math.Min(math.Abs(prediction)/10, 0.3)
	
	// Ensure minimum confidence is 35% and maximum is 95%
	return math.Max(0.35, math.Min(0.95, baseConfidence-magnitudePenalty))
}

// identifyKeyFactors identifies key factors that influenced the prediction
func (tf *TFModelService) identifyKeyFactors(stock *StockInfo, features map[string]float64, prediction float64) []string {
	var factors []string
	
	// Create a list of factor descriptions and their importance
	// Calculate market cap factor
	var marketCapFactor float64
	if features["market_cap"] > 1e11 {
		marketCapFactor = -0.1
	} else if features["market_cap"] < 1e9 {
		marketCapFactor = 0.2
	}

	factorImportances := []struct {
		description string
		importance  float64
	}{
		{"Price momentum", features["change_pct"] * 0.1},
		{"Position in day's range", (features["price_position"] - 0.5) * -0.2},
		{"Volatility", features["volatility"] * 0.05},
		{"Market capitalization", marketCapFactor},
		{"Open price position", features["open_gap"] * 0.1},
	}
	
	// Sort by absolute importance
	sort.Slice(factorImportances, func(i, j int) bool {
		return math.Abs(factorImportances[i].importance) > math.Abs(factorImportances[j].importance)
	})
	
	// Take top 3 factors
	for i, factor := range factorImportances {
		if i >= 3 {
			break
		}
		
		// Skip factors with near-zero importance
		if math.Abs(factor.importance) < 0.05 {
			continue
		}
		
		// Create description based on factor importance
		var description string
		if factor.importance > 0 {
			description = fmt.Sprintf("%s appears positive", factor.description)
		} else {
			description = fmt.Sprintf("%s appears negative", factor.description)
		}
		
		factors = append(factors, description)
	}
	
	// Add market context
	if stock.Price > stock.Open {
		factors = append(factors, "Price is above opening")
	} else if stock.Price < stock.Open {
		factors = append(factors, "Price is below opening")
	}
	
	return factors
}
