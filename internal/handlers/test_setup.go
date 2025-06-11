package handlers

import (
	"html/template"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// Mock the templates before running tests
	mockTemplates()
	
	// Run all tests
	result := m.Run()
	
	// Exit with the test result code
	os.Exit(result)
}

// mockTemplates overrides the templates variable for testing
func mockTemplates() {
	// Create a simple template with required structure for tests
	tmpl, err := template.New("index.html").Funcs(funcMap).Parse(`
		<html>
		<head><title>Stock Analyzer</title></head>
		<body>
			<h1>Stock Analyzer</h1>
			<form>
				<input type="text" placeholder="Enter company name or ticker" />
			</form>
			{{if .Error}}
				<div>{{.Error}}</div>
			{{end}}
			{{if .StockInfo}}
				<div>{{.StockInfo.CompanyName}}</div>
			{{end}}
		</body>
		</html>
	`)
	
	if err != nil {
		// Don't fail the test setup, just log the error
		// Tests will fail later if template is needed
		return
	}
	
	// Replace the existing templates
	templates = tmpl
	
	// Wait a moment to ensure templates are initialized
	time.Sleep(10 * time.Millisecond)
}
