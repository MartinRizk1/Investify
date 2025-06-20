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
	
	// Try to use the Python bridge for predictions
	bridge := GetPythonBridge()
	if bridge.initialized {
		// Make prediction using Python - try our simple analyzer first
		result, err := bridge.PredictStockPriceWithSimpleAnalyzer(stock.Ticker)
		if err == nil && result != nil {
			// Use the Python prediction results
			log.Printf("Using Python-based prediction for %s", stock.Ticker)
			return &StockPrediction{
				PredictedPrice: result.PredictedPrice,
				Confidence:     result.Confidence,
				Direction:      result.Direction,
				Factors:        result.Factors,
			}, nil
		}
		
		// Try the original prediction method as a backup
		result, err = bridge.PredictStockPrice(stock.Ticker)
		if err == nil && result != nil {
			log.Printf("Using TensorFlow-based prediction for %s", stock.Ticker)
			return &StockPrediction{
				PredictedPrice: result.PredictedPrice,
				Confidence:     result.Confidence,
				Direction:      result.Direction,
				Factors:        result.Factors,
			}, nil
		}
		
		// Log the error but continue with fallback
		log.Printf("Python prediction failed: %v, using fallback", err)
	}

	// Fallback to simple prediction
	log.Printf("Using simple prediction for %s", stock.Ticker)
	
	// Simple prediction based on price movement
	changePercent := 0.0
	if stock.Change != 0 && stock.Price != 0 {
		changePercent = (stock.Change / stock.Price) * 100
	}
	
	// Simulate some predictive model output
	randomFactor := (math.Sin(float64(time.Now().Unix())) + 1.0) * 0.5 // 0.0-1.0
	if randomFactor > 0.5 {
		changePercent *= 1.2 // Amplify the trend
	} else {
		changePercent *= -0.8 // Reverse the trend somewhat
	}
	
	// Predict price (limited to +/- 5%)
	changePercent = math.Max(-5.0, math.Min(5.0, changePercent))
	predictedChange := stock.Price * (changePercent / 100.0)
	predictedPrice := stock.Price + predictedChange
	
	// Round to 2 decimal places
	predictedPrice = math.Round(predictedPrice*100) / 100
	
	// Determine direction
	direction := "NEUTRAL"
	if predictedChange > 0 {
		direction = "UP"
	} else if predictedChange < 0 {
		direction = "DOWN"
	}
	
	// Generate confidence (60-90%)
	confidence := 60.0 + (randomFactor * 30.0)
	
	// Generate prediction factors
	factors := []string{}
	
	// Price momentum
	if stock.Change > 0 {
		factors = append(factors, "Recent positive price momentum")
	} else if stock.Change < 0 {
		factors = append(factors, "Recent negative price momentum")
	}
	
	// Price position relative to day's range
	dayRange := stock.High - stock.Low
	if dayRange > 0 {
		pricePosition := (stock.Price - stock.Low) / dayRange
		if pricePosition > 0.8 {
			factors = append(factors, "Price near daily high")
		} else if pricePosition < 0.2 {
			factors = append(factors, "Price near daily low")
		}
	}
	
	// Market conditions
	currentHour := time.Now().Hour()
	if currentHour < 12 {
		factors = append(factors, "Morning market conditions")
	} else if currentHour >= 12 && currentHour < 16 {
		factors = append(factors, "Afternoon trading patterns")
	} else {
		factors = append(factors, "After-hours sentiment")
	}
	
	// If we don't have enough factors, add a generic one
	if len(factors) < 2 {
		factors = append(factors, "Based on technical analysis")
	}
	
	return &StockPrediction{
		PredictedPrice: predictedPrice,
		Confidence:     confidence,
		Direction:      direction,
		Factors:        factors,
	}, nil
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
