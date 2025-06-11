package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// AIService handles AI-based stock recommendations
type AIService struct {
	openAIKey string
}

// OpenAIRequest represents the request to OpenAI API
type OpenAIRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse represents the response from OpenAI API
type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// NewAIService creates a new AI service instance
func NewAIService(openAIKey string) *AIService {
	return &AIService{
		openAIKey: openAIKey,
	}
}

// GetStockRecommendation generates an AI-based stock recommendation
func (ai *AIService) GetStockRecommendation(stock *StockInfo) (string, error) {
	// If no OpenAI key is provided, use rule-based recommendation
	if ai.openAIKey == "" {
		return ai.getRuleBasedRecommendation(stock), nil
	}

	// Create prompt for OpenAI
	prompt := fmt.Sprintf(`Analyze this stock and provide a recommendation (BUY, SELL, or HOLD) with a brief explanation:
	
Stock: %s (%s)
Current Price: $%.2f
Daily Change: $%.2f (%s)
Open: $%.2f
High: $%.2f
Low: $%.2f
Volume: %s
Market Cap: %s

Please provide a concise recommendation with reasoning based on the data provided.`,
		stock.CompanyName, stock.Ticker, stock.Price, stock.Change, stock.ChangePct,
		stock.Open, stock.High, stock.Low, stock.Volume, stock.MarketCap)

	// Make request to OpenAI
	reqBody := OpenAIRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a financial advisor providing stock recommendations based on market data.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens: 150,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return ai.getRuleBasedRecommendation(stock), nil
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return ai.getRuleBasedRecommendation(stock), nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ai.openAIKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("OpenAI API request failed: %v", err)
		// Check for specific network errors
		if strings.Contains(err.Error(), "timeout") {
			log.Printf("OpenAI API timeout - falling back to rule-based recommendation")
		} else if strings.Contains(err.Error(), "no such host") || strings.Contains(err.Error(), "lookup") {
			log.Printf("OpenAI API network connectivity issue - falling back to rule-based recommendation")
		}
		return ai.getRuleBasedRecommendation(stock), nil
	}
	defer resp.Body.Close()
	
	// Check status code
	if resp.StatusCode != http.StatusOK {
		log.Printf("OpenAI API returned non-200 status code: %d", resp.StatusCode)
		
		// Handle specific error codes
		switch resp.StatusCode {
		case http.StatusTooManyRequests:
			log.Printf("OpenAI rate limit exceeded - falling back to rule-based recommendation")
		case http.StatusUnauthorized:
			log.Printf("OpenAI API key invalid or expired - falling back to rule-based recommendation")
		default:
			log.Printf("OpenAI API error with status code %d - falling back to rule-based recommendation", resp.StatusCode)
		}
		
		return ai.getRuleBasedRecommendation(stock), nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read OpenAI API response: %v", err)
		return ai.getRuleBasedRecommendation(stock), nil
	}

	var openAIResp OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		log.Printf("Failed to parse OpenAI API response: %v", err)
		log.Printf("Response body: %s", string(body))
		return ai.getRuleBasedRecommendation(stock), nil
	}

	if len(openAIResp.Choices) > 0 {
		return openAIResp.Choices[0].Message.Content, nil
	}

	return ai.getRuleBasedRecommendation(stock), nil
}

// getRuleBasedRecommendation provides a fallback rule-based recommendation
func (ai *AIService) getRuleBasedRecommendation(stock *StockInfo) string {
	changeFloat := stock.Change
	changePct := calculateChangePercentage(stock.Change, stock.Price)

	// Simple technical analysis based on price movements
	dayRange := stock.High - stock.Low
	if dayRange == 0 {
		return "HOLD - Insufficient price data for analysis"
	}

	pricePosition := (stock.Price - stock.Low) / dayRange

	recommendation := "HOLD"
	reason := "Neutral market conditions"

	// Strong upward momentum
	if changePct > 5 && pricePosition > 0.8 {
		recommendation = "SELL"
		reason = "Strong gains suggest potential profit-taking opportunity"
	} else if changePct < -5 && pricePosition < 0.3 {
		recommendation = "BUY"
		reason = "Significant dip presents potential buying opportunity"
	} else if changePct > 2 && pricePosition > 0.6 {
		recommendation = "HOLD/BUY"
		reason = "Positive momentum with room for growth"
	} else if changePct < -2 && pricePosition < 0.4 {
		recommendation = "HOLD/BUY"
		reason = "Minor decline, potential value opportunity"
	} else if changeFloat > 0 && pricePosition > 0.5 {
		recommendation = "HOLD"
		reason = "Stable upward trend"
	} else if changeFloat < 0 && pricePosition < 0.5 {
		recommendation = "HOLD"
		reason = "Monitor for further developments"
	}

	return fmt.Sprintf("%s - %s", recommendation, reason)
}

// calculateChangePercentage calculates the percentage change
func calculateChangePercentage(change, price float64) float64 {
	if price == 0 {
		return 0
	}
	return (change / (price - change)) * 100
}
