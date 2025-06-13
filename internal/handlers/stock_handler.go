package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/martinrizk/investify/internal/services"
)

type PageData struct {
	StockInfo *services.StockInfo
	Error     string
	LastQuery string
}

var templates *template.Template

// Custom template functions
var funcMap = template.FuncMap{
	"ToLower":  strings.ToLower,
	"contains": strings.Contains,
	"printTrendArrow": func(trend string) string {
		trend = strings.ToUpper(trend)
		if trend == "UP" {
			return "↑"
		} else if trend == "DOWN" {
			return "↓"
		}
		return "→"
	},
	"printTrendClass": func(trend string) string {
		trend = strings.ToUpper(trend)
		if trend == "UP" {
			return "positive"
		} else if trend == "DOWN" {
			return "negative"
		}
		return "neutral"
	},
	"printTrendIcon": func(trend string) string {
		trend = strings.ToUpper(trend)
		if trend == "UP" {
			return "📈"
		} else if trend == "DOWN" {
			return "📉"
		}
		return "📊"
	},
	"getPercentInRange": func(value, min, max float64) float64 {
		if max == min {
			return 50.0 // Avoid division by zero
		}
		percent := (value - min) / (max - min) * 100
		// Clamp between 0-100
		if percent < 0 {
			percent = 0
		} else if percent > 100 {
			percent = 100
		}
		return percent
	},
}

// IsTestMode is used to check if code is running in test mode
var IsTestMode bool

func init() {
	var err error
	// Add the custom functions to the template
	if !IsTestMode {
		// Try different paths to find templates
		templatePaths := []string{
			"templates/*.html",                   // Run from project root
			"../templates/*.html",                // Run from internal directory
			"../../templates/*.html",             // Run from internal/handlers directory
		}
		
		var templateErr error
		for _, path := range templatePaths {
			log.Printf("Trying template path: %s", path)
			templates, err = template.New("").Funcs(funcMap).ParseGlob(path)
			if err == nil {
				// Found templates, break the loop
				log.Printf("Successfully loaded templates from: %s", path)
				log.Printf("Loaded templates: %v", templates.DefinedTemplates())
				templateErr = nil
				break
			}
			log.Printf("Failed to load from %s: %v", path, err)
			templateErr = err
		}
		
		// If we still have an error after trying all paths
		if templateErr != nil {
			log.Printf("Error loading templates from paths: %v", templateErr)
			log.Fatalf("Could not find templates in any expected location")
		}
	}
	// Note: Test mode templates are initialized in the handler
}

// isValidStockQuery validates a stock ticker or company name input
// to prevent injection attacks and ensure valid input
func isValidStockQuery(query string) bool {
	// Allow letters, numbers, spaces, dots, and some special characters commonly used in company names
	// Limit length to prevent abuse
	if len(query) > 100 {
		return false
	}

	// Using regexp to validate the input pattern
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9\s\.\-&']+$`)
	return validPattern.MatchString(query)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	var data PageData
	
	// Debug logging
	log.Printf("HomeHandler called with method: %s", r.Method)
	log.Printf("Templates initialized: %v", templates != nil)
	if templates != nil {
		log.Printf("Available templates: %v", templates.DefinedTemplates())
	}
	
	// Ensure template is initialized for test mode
	if IsTestMode && templates == nil {
		tmpl, err := template.New("index.html").Funcs(funcMap).Parse(`<html><body><h1>Investify Test Mode</h1>{{if .Error}}<div class="error">{{.Error}}</div>{{end}}{{if .StockInfo}}<div class="stock-info"><h2>{{.StockInfo.Ticker}}</h2><p>{{.StockInfo.CompanyName}}</p><p>Price: ${{printf "%.2f" .StockInfo.Price}}</p></div>{{else}}<div class="search-form"><form method="POST"><input type="text" name="ticker" placeholder="Enter ticker" value="{{.LastQuery}}"><button type="submit">Search</button></form></div>{{end}}</body></html>`)
		if err != nil {
			http.Error(w, "Template creation error", http.StatusInternalServerError)
			return
		}
		templates = tmpl
	}
	
	if r.Method == http.MethodPost {
		query := strings.TrimSpace(r.FormValue("ticker"))
		data.LastQuery = query
		
		// Validate input - only allow alphanumeric characters and some special chars
		if query == "" {
			data.Error = "Please enter a company name or ticker symbol."
		} else {
			// Use the secure version of SearchStock that validates input
			stockInfo, err := services.SearchStockSecure(query)
			if err != nil {
				data.Error = fmt.Sprintf("%v", err)
			} else {
				data.StockInfo = stockInfo
			}
		}
	}

	// Render the main template
	var buf bytes.Buffer
	log.Printf("About to render template. Data: %+v", data)
	
	if IsTestMode {
		err := templates.Execute(&buf, data)
		if err != nil {
			log.Printf("Template execution error in test mode: %v", err)
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
	} else {
		err := templates.ExecuteTemplate(&buf, "index.html", data)
		if err != nil {
			log.Printf("Template execution error: %v", err)
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
	}
	
	log.Printf("Template rendered successfully. Buffer length: %d", buf.Len())
	
	// Write to response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Add security headers
	w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net https://cdnjs.cloudflare.com; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; img-src 'self' data:")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
	w.Write(buf.Bytes())
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// StaticFileHandler serves static files like JavaScript and CSS
func StaticFileHandler(w http.ResponseWriter, r *http.Request) {
	// Extract filename from path
	filename := filepath.Base(r.URL.Path)
	
	// Determine file extension
	ext := filepath.Ext(filename)
	
	// Set content type based on file extension
	switch ext {
	case ".js":
		w.Header().Set("Content-Type", "application/javascript")
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	default:
		w.Header().Set("Content-Type", "text/plain")
	}
	
	// Try different paths to find the static file
	paths := []string{
		filepath.Join("templates", filename),
		filepath.Join("../templates", filename),
		filepath.Join("../../templates", filename),
	}
	
	// Try to read and serve the file
	for _, path := range paths {
		if data, err := os.ReadFile(path); err == nil {
			log.Printf("Serving static file: %s", path)
			w.Write(data)
			return
		}
	}
	
	// If file not found, return 404
	http.NotFound(w, r)
}
