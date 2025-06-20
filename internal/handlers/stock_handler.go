package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
)

// IsTestMode is used to check if code is running in test mode
var IsTestMode bool

// IsValidStockQuery validates a stock ticker or company name input
// to prevent injection attacks and ensure valid input
func IsValidStockQuery(query string) bool {
	// Allow letters, numbers, spaces, dots, and some special characters commonly used in company names
	// Limit length to prevent abuse
	if len(query) > 100 {
		return false
	}

	// Using regexp to validate the input pattern
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9\s\.\-&']+$`)
	return validPattern.MatchString(query)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
