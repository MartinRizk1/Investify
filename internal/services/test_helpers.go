package services

import (
	"time"
)

// Exported functions to support testing

// FormatVolume formats a volume number into a human-readable string
func FormatVolume(volume int64) string {
	return formatVolume(volume)
}

// FormatMarketCap formats a market cap number into a human-readable string
func FormatMarketCap(marketCap int64) string {
	return formatMarketCap(marketCap)
}

// GetRuleBasedRecommendation provides a rule-based stock recommendation
func GetRuleBasedRecommendation(stock *StockInfo) string {
	if aiService != nil {
		return aiService.getRuleBasedRecommendation(stock)
	}
	
	// Fallback implementation
	if stock.Change > 0 {
		return "BUY - Stock shows positive momentum"
	} else {
		return "HOLD - Stock shows negative momentum"
	}
}

// PredictStockMovement predicts stock price movement using TF model
func PredictStockMovement(stock *StockInfo) (*StockPrediction, error) {
	if tfModelService != nil {
		return tfModelService.PredictStockMovement(stock)
	}
	return nil, nil
}

// CacheStockInfo adds a stock to the cache
func CacheStockInfo(key string, stock *StockInfo) {
	stockCache[key] = &CachedStock{
		Data:      stock,
		Timestamp: now(),
	}
}

// GetCachedStock retrieves a stock from cache
func GetCachedStock(key string) *StockInfo {
	if cached, ok := stockCache[key]; ok {
		cached.Data.DataAge = int64(now().Sub(cached.Timestamp).Seconds())
		return cached.Data
	}
	return nil
}

// Wrapper for time.Now() to make testing easier
func now() time.Time {
	return time.Now()
}
