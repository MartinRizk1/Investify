package handlers

import (
	"strings"
	"testing"
)

func init() {
	// Enable test mode for all handler tests
	IsTestMode = true
}

// TestIsValidStockQuery tests the stock query validation
func TestIsValidStockQuery(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  bool
	}{
		{
			name:  "Valid stock ticker",
			query: "AAPL",
			want:  true,
		},
		{
			name:  "Valid company name",
			query: "Apple Inc",
			want:  true,
		},
		{
			name:  "Valid with special chars",
			query: "Berkshire Hathaway-B",
			want:  true,
		},
		{
			name:  "Invalid with script tags",
			query: "<script>alert('xss')</script>",
			want:  false,
		},
		{
			name:  "Too long input",
			query: strings.Repeat("A", 101),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidStockQuery(tt.query); got != tt.want {
				t.Errorf("IsValidStockQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
