package main

import (
	"html/template"
	"log"
	"os"
	"strings"
)

// Custom template functions
var funcMap = template.FuncMap{
	"ToLower":  strings.ToLower,
	"contains": strings.Contains,
}

type PageData struct {
	StockInfo interface{}
	Error     string
	LastQuery string
}

func main() {
	// Test template parsing
	tmpl, err := template.New("").Funcs(funcMap).ParseFiles("templates/index.html")
	if err != nil {
		log.Fatalf("Template parsing error: %v", err)
	}
	
	log.Printf("Template parsed successfully!")
	log.Printf("Defined templates: %v", tmpl.DefinedTemplates())
	
	// Test template execution
	data := PageData{
		StockInfo: nil,
		Error:     "",
		LastQuery: "",
	}
	
	err = tmpl.ExecuteTemplate(os.Stdout, "index.html", data)
	if err != nil {
		log.Fatalf("Template execution error: %v", err)
	}
}
