package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
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
			templates, err = template.New("").Funcs(funcMap).ParseGlob(path)
			if err == nil {
				// Found templates, break the loop
				templateErr = nil
				break
			}
			templateErr = err
		}
		
		// If we still have an error after trying all paths
		if templateErr != nil {
			log.Printf("Error loading templates from paths: %v", templateErr)
			log.Fatalf("Could not find templates in any expected location")
		}
	} else {
		// Create a simple template for tests
		templates, _ = template.New("index.html").Funcs(funcMap).Parse(`
			<html><body>{{.Error}}{{if .StockInfo}}{{.StockInfo.CompanyName}}{{end}}</body></html>
		`)
		log.Printf("Using simplified templates for tests")
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	var data PageData
	
	if r.Method == http.MethodPost {
		query := strings.TrimSpace(r.FormValue("ticker"))
		data.LastQuery = query
		
		if query == "" {
			data.Error = "Please enter a company name or ticker symbol."
		} else {
			stockInfo, err := services.SearchStock(query)
			if err != nil {
				data.Error = fmt.Sprintf("%v", err)
			} else {
				data.StockInfo = stockInfo
			}
		}
	}

	// Try to render index-fixed.html first, fallback to index.html if not available
	err := templates.ExecuteTemplate(w, "index-fixed.html", data)
	if err != nil {
		log.Printf("Trying to use index-fixed.html: %v, falling back to index.html", err)
		if err := templates.ExecuteTemplate(w, "index.html", data); err != nil {
			log.Printf("Template error: %v", err)
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
		}
	}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
