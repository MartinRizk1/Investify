package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// StockInfo represents stock information with AI analysis and ML predictions
type StockInfo struct {
	Ticker              string   `json:"ticker"`
	CompanyName         string   `json:"company_name"`
	Price               float64  `json:"price"`
	Change              float64  `json:"change"`
	ChangePct           string   `json:"change_pct"`
	Open                float64  `json:"open"`
	High                float64  `json:"high"`
	Low                 float64  `json:"low"`
	Volume              string   `json:"volume"`
	MarketCap           string   `json:"market_cap"`
	Recommendation      string   `json:"recommendation"`
	AIAnalysis          string   `json:"ai_analysis"`
	PredictedPrice      float64  `json:"predicted_price"`
	PredictionConfidence float64  `json:"prediction_confidence"`
	TrendDirection      string   `json:"trend_direction"`
	KeyFactors          []string `json:"key_factors"`
	DataAge             int64    `json:"data_age"` // Time in seconds since data was retrieved
}

type YahooResponse struct {
	QuoteResponse struct {
		Result []struct {
			Symbol                     string  `json:"symbol"`
			ShortName                  string  `json:"shortName"`
			LongName                   string  `json:"longName"`
			RegularMarketPrice         float64 `json:"regularMarketPrice"`
			RegularMarketChange        float64 `json:"regularMarketChange"`
			RegularMarketChangePercent float64 `json:"regularMarketChangePercent"`
			RegularMarketOpen          float64 `json:"regularMarketOpen"`
			RegularMarketDayHigh       float64 `json:"regularMarketDayHigh"`
			RegularMarketDayLow        float64 `json:"regularMarketDayLow"`
			RegularMarketVolume        int64   `json:"regularMarketVolume"`
			MarketCap                  int64   `json:"marketCap"`
		} `json:"result"`
		Error *struct {
			Code        string `json:"code"`
			Description string `json:"description"`
		} `json:"error,omitempty"`
	} `json:"quoteResponse"`
}

var (
	apiFailureCount  = 0
	lastApiCallTime  time.Time
	aiService        *AIService
	tfModelService   *TFModelService
)

// Cache system to reduce API calls
var stockCache = make(map[string]*CachedStock)

type CachedStock struct {
	Data      *StockInfo
	Timestamp time.Time
}

// Common ticker mappings for popular companies
var companyNameToTicker = map[string]string{
	"GOOGLE":     "GOOGL",
	"ALPHABET":   "GOOGL",
	"FACEBOOK":   "META",
	"TWITTER":    "TWTR",
	"X":          "TWTR", 
	"TESLA":      "TSLA",
	"APPLE":      "AAPL",
	"MICROSOFT":  "MSFT",
	"AMAZON":     "AMZN",
	"NETFLIX":    "NFLX",
	"UBER":       "UBER",
	"LYFT":       "LYFT",
	"NVIDIA":     "NVDA",
	"NVIDIA CORPORATION": "NVDA",
	"AMD":        "AMD",
	"ADVANCED MICRO DEVICES": "AMD",
	"INTEL":      "INTC",
	"BERKSHIRE":  "BRK-B",
	"BERKSHIRE HATHAWAY": "BRK-B",
	"JPMORGAN":   "JPM",
	"JP MORGAN":  "JPM",
	"WALMART":    "WMT",
	"COCA-COLA":  "KO",
	"COKE":       "KO",
	"MCDONALDS":  "MCD",
	"BOEING":     "BA",
	"DISNEY":     "DIS",
	"WALT DISNEY": "DIS",
	"PAYPAL":     "PYPL",
	"SALESFORCE": "CRM",
	"ZOOM":       "ZM",
	"SPOTIFY":    "SPOT",
	"SHOPIFY":    "SHOP",
	"SQUARE":     "SQ",
	"BLOCK":      "SQ",
	"IBM":        "IBM",
	"ORACLE":     "ORCL",
	"CISCO":      "CSCO",
	"ADOBE":      "ADBE",
	"QUALCOMM":   "QCOM",
	"TEXAS INSTRUMENTS": "TXN",
	"NINTENDO":   "NTDOY",
	"SONY":       "SONY",
	"MERCEDES":   "MBG.DE",
	"MERCEDES BENZ": "MBG.DE",
	"BMW":        "BMW.DE",
	"VOLKSWAGEN": "VOW3.DE",
	"TOYOTA":     "TM",
	"FORD":       "F",
	"GENERAL MOTORS": "GM",
	"GM":         "GM",
	"GAMESTOP":   "GME",
	"AMC":        "AMC",
	"ROBINHOOD":  "HOOD",
	"COINBASE":   "COIN",
	"SNAPCHAT":   "SNAP",
	"PINTEREST":  "PINS",
	"REDDIT":     "RDDT",
	"TIKTOK":     "BDNCE", // ByteDance
}

func init() {
	// Initialize AI service with OpenAI key from environment
	openAIKey := os.Getenv("OPENAI_API_KEY")
	aiService = NewAIService(openAIKey)
	
	// Initialize TensorFlow model service
	tfModelService = NewTFModelService()

	// Log service initialization
	log.Println("Stock services initialized. OpenAI API key present:", openAIKey != "")
}

// SearchStock searches for a stock by company name or ticker
func SearchStock(query string) (*StockInfo, error) {
	// Normalize input
	input := strings.ToUpper(strings.TrimSpace(query))
	
	// Try to match company name to ticker
	ticker := input
	if mappedTicker, ok := companyNameToTicker[input]; ok {
		ticker = mappedTicker
		log.Printf("Mapped company name '%s' to ticker '%s'", input, ticker)
	}
	
	// Check cache first
	if cached, ok := stockCache[ticker]; ok {
		// If cache is less than 5 minutes old, use it
		if time.Since(cached.Timestamp) < 5*time.Minute {
			log.Printf("Using cached data for %s (age: %v)", ticker, time.Since(cached.Timestamp))
			cached.Data.DataAge = int64(time.Since(cached.Timestamp).Seconds())
			return cached.Data, nil
		}
		log.Printf("Cached data for %s expired, fetching fresh data", ticker)
	}
	
	return FetchStockInfo(ticker)
}

func FetchStockInfo(ticker string) (*StockInfo, error) {
	// Validate and normalize ticker
	ticker = strings.ToUpper(strings.TrimSpace(ticker))
	if ticker == "" {
		return nil, fmt.Errorf("please enter a valid ticker symbol or company name")
	}
	
	// Implement rate limiting and backoff if we've had recent API failures
	if apiFailureCount > 3 && time.Since(lastApiCallTime) < 5*time.Minute {
		return nil, fmt.Errorf("yahoo Finance API is currently experiencing issues - please try again in a few minutes")
	}
	
	log.Printf("Fetching stock data for ticker: %s", ticker)
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v7/finance/quote?symbols=%s", ticker)
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Add user agent to avoid blocking
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	
	// Update the last API call time
	lastApiCallTime = time.Now()

	resp, err := client.Do(req)
	if err != nil {
		// Increment failure count
		apiFailureCount++
		log.Printf("Yahoo Finance API request failed: %v", err)
		
		// Classify network errors for better user feedback
		if strings.Contains(err.Error(), "timeout") {
			return nil, fmt.Errorf("the Yahoo Finance API is currently slow or unavailable - please try again later")
		}
		if strings.Contains(err.Error(), "no such host") || strings.Contains(err.Error(), "lookup") {
			return nil, fmt.Errorf("network connection issue: unable to reach Yahoo Finance servers - please check your internet connection")
		}
		if strings.Contains(err.Error(), "refused") || strings.Contains(err.Error(), "reset") {
			return nil, fmt.Errorf("connection to Yahoo Finance failed - the service may be temporarily unavailable")
		}
		
		// Add exponential backoff based on failure count
		backoffTime := time.Duration(math.Min(float64(apiFailureCount*apiFailureCount), 15)) * time.Minute
		log.Printf("Setting backoff time to %v after %d failures", backoffTime, apiFailureCount)
		
		return nil, fmt.Errorf("network error connecting to Yahoo Finance: %v - please try again later", err)
	}
	defer resp.Body.Close()
	
	// Check HTTP status code with detailed responses
	if resp.StatusCode != http.StatusOK {
		apiFailureCount++
		log.Printf("Yahoo Finance API returned non-200 status code: %d", resp.StatusCode)
		
		switch resp.StatusCode {
		case http.StatusTooManyRequests:
			return nil, fmt.Errorf("rate limit exceeded - our application is making too many requests to yahoo finance. Please try again in a few minutes")
		case http.StatusServiceUnavailable, http.StatusGatewayTimeout:
			return nil, fmt.Errorf("yahoo finance service is temporarily unavailable. Please try again later")
		case http.StatusForbidden, http.StatusUnauthorized:
			return nil, fmt.Errorf("access to yahoo finance API is restricted. This could be temporary - please try again later")
		case http.StatusNotFound:
			return nil, fmt.Errorf("ticker symbol '%s' not found. Please check the spelling and try again", ticker)
		default:
			return nil, fmt.Errorf("yahoo finance API error (code %d). Please try again later", resp.StatusCode)
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read Yahoo Finance API response: %v", err)
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var yahooResp YahooResponse
	if err := json.Unmarshal(body, &yahooResp); err != nil {
		log.Printf("Failed to parse Yahoo Finance API response: %v", err)
		log.Printf("Response body: %s", string(body))
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}
	
	// Check if Yahoo returned an API error
	if yahooResp.QuoteResponse.Error != nil {
		apiFailureCount++
		log.Printf("Yahoo Finance API error: %s", yahooResp.QuoteResponse.Error.Description)
		return nil, fmt.Errorf("yahoo API error: %s", yahooResp.QuoteResponse.Error.Description)
	}

	// Reset API failure count on success
	apiFailureCount = 0

	if len(yahooResp.QuoteResponse.Result) == 0 {
		log.Printf("No stock data found for ticker: %s", ticker)
		
		// Look for similar company names or tickers that might match
		var suggestions []string
		
		// Try to find similar companies
		for company, symbol := range companyNameToTicker {
			// Check if company name contains the search term
			if strings.Contains(company, ticker) || strings.Contains(ticker, company) {
				suggestions = append(suggestions, fmt.Sprintf("%s (%s)", company, symbol))
				if len(suggestions) >= 5 {
					break
				}
			}
		}
		
		if len(suggestions) > 0 {
			log.Printf("Suggesting alternatives for %s: %v", ticker, suggestions)
			return nil, fmt.Errorf("no data found for '%s'. Did you mean: %s?", ticker, strings.Join(suggestions, ", "))
		}
		
		return nil, fmt.Errorf("no data found for '%s'. Try using a valid ticker symbol like AAPL (Apple) or MSFT (Microsoft)", ticker)
	}

	quote := yahooResp.QuoteResponse.Result[0]
	log.Printf("Retrieved data for %s (%s): $%.2f (%+.2f%%)", 
		getCompanyName(quote.LongName, quote.ShortName, quote.Symbol),
		quote.Symbol,
		quote.RegularMarketPrice,
		quote.RegularMarketChangePercent)

	// Validate critical data points
	if quote.RegularMarketPrice <= 0 {
		return nil, fmt.Errorf("invalid price data received from Yahoo Finance API")
	}

	stockInfo := &StockInfo{
		Ticker:      quote.Symbol,
		CompanyName: getCompanyName(quote.LongName, quote.ShortName, quote.Symbol),
		Price:       quote.RegularMarketPrice,
		Change:      quote.RegularMarketChange,
		ChangePct:   fmt.Sprintf("%.2f%%", quote.RegularMarketChangePercent),
		Open:        quote.RegularMarketOpen,
		High:        quote.RegularMarketDayHigh,
		Low:         quote.RegularMarketDayLow,
		Volume:      formatVolume(quote.RegularMarketVolume),
		MarketCap:   formatMarketCap(quote.MarketCap),
		DataAge:     0, // Just fetched
	}

	// Get AI recommendation from OpenAI
	recommendation, err := aiService.GetStockRecommendation(stockInfo)
	if err != nil {
		log.Printf("Failed to get AI recommendation: %v", err)
		stockInfo.Recommendation = "HOLD - Unable to generate recommendation"
	} else {
		log.Printf("AI recommendation for %s: %s", stockInfo.Ticker, recommendation)
		stockInfo.Recommendation = recommendation
	}
	
	// Get TensorFlow predictions
	if tfModelService != nil && stockInfo.Price > 0 {
		prediction, err := tfModelService.PredictStockMovement(stockInfo)
		if err == nil && prediction != nil {
			stockInfo.PredictedPrice = prediction.PredictedPrice
			stockInfo.PredictionConfidence = prediction.Confidence * 100
			stockInfo.TrendDirection = prediction.Direction
			stockInfo.KeyFactors = prediction.Factors
			
			// Generate AI analysis based on TF prediction
			priceDiff := ((prediction.PredictedPrice - stockInfo.Price) / stockInfo.Price) * 100
			directionText := "remain stable"
			if priceDiff > 1.0 {
				directionText = "rise"
			} else if priceDiff < -1.0 {
				directionText = "fall"
			}
			
			stockInfo.AIAnalysis = fmt.Sprintf("AI model predicts the stock will %s with %.1f%% confidence. Predicted price: $%.2f (%.1f%%)",
				directionText, 
				prediction.Confidence * 100,
				prediction.PredictedPrice,
				priceDiff)
				
			log.Printf("ML prediction for %s: %s, confidence: %.1f%%, direction: %s", 
				stockInfo.Ticker, 
				stockInfo.AIAnalysis,
				stockInfo.PredictionConfidence,
				stockInfo.TrendDirection)
		} else if err != nil {
			log.Printf("ML prediction error: %v", err)
		}
	}

	// Cache the results
	stockCache[ticker] = &CachedStock{
		Data:      stockInfo,
		Timestamp: time.Now(),
	}

	return stockInfo, nil
}

func getCompanyName(longName, shortName, symbol string) string {
	if longName != "" {
		return longName
	}
	if shortName != "" {
		return shortName
	}
	return symbol
}

func formatVolume(volume int64) string {
	if volume >= 1e9 {
		return fmt.Sprintf("%.2fB", float64(volume)/1e9)
	} else if volume >= 1e6 {
		return fmt.Sprintf("%.2fM", float64(volume)/1e6)
	} else if volume >= 1e3 {
		return fmt.Sprintf("%.2fK", float64(volume)/1e3)
	}
	return strconv.FormatInt(volume, 10)
}

func formatMarketCap(marketCap int64) string {
	if marketCap >= 1e12 {
		return fmt.Sprintf("$%.2fT", float64(marketCap)/1e12)
	} else if marketCap >= 1e9 {
		return fmt.Sprintf("$%.2fB", float64(marketCap)/1e9)
	} else if marketCap >= 1e6 {
		return fmt.Sprintf("$%.2fM", float64(marketCap)/1e6)
	}
	return fmt.Sprintf("$%d", marketCap)
}
