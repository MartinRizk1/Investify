package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// sanitizeSearchQuery sanitizes and validates the search input
func sanitizeSearchQuery(query string) (string, error) {
	// Validate input length to prevent abuse
	if len(query) > 50 {
		return "", fmt.Errorf("search query too long")
	}
	
	// Sanitize input - allow only alphanumeric and spaces
	re := regexp.MustCompile("[^a-zA-Z0-9 ]")
	cleanQuery := re.ReplaceAllString(query, "")
	
	// Normalize and trim
	return strings.ToUpper(strings.TrimSpace(cleanQuery)), nil
}

// SearchStock searches for a stock by company name or ticker
// This is a secure version with input validation
func SearchStockSecure(query string) (*StockInfo, error) {
	// Sanitize and validate input
	input, err := sanitizeSearchQuery(query)
	if err != nil {
		return nil, err
	}
	
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
