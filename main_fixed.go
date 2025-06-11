package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type StockInfo struct {
	Ticker         string  `json:"ticker"`
	CompanyName    string  `json:"company_name"`
	Price          float64 `json:"price"`
	Change         float64 `json:"change"`
	ChangePct      string  `json:"change_pct"`
	Open           float64 `json:"open"`
	High           float64 `json:"high"`
	Low            float64 `json:"low"`
	Volume         string  `json:"volume"`
	MarketCap      string  `json:"market_cap"`
	Recommendation string  `json:"recommendation"`
}

type PageData struct {
	Info      *StockInfo `json:"info"`
	Error     string     `json:"error"`
	LastQuery string     `json:"last_query"`
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
	} `json:"quoteResponse"`
}

var templates *template.Template

func main() {
	var err error
	templates, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Printf("Template error: %v, using fallback", err)
		templates = template.Must(template.New("index").Parse(fallbackTemplate))
	}

	r := mux.NewRouter()
	r.HandleFunc("/", handleHome).Methods("GET", "POST")
	r.HandleFunc("/health", handleHealth).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Starting MoneyMaker on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	var data PageData

	if r.Method == "POST" {
		ticker := strings.TrimSpace(strings.ToUpper(r.FormValue("ticker")))
		data.LastQuery = ticker
		
		if ticker == "" {
			data.Error = "Please enter a stock ticker symbol."
		} else {
			info, err := fetchStock(ticker)
			if err != nil {
				data.Error = fmt.Sprintf("Error fetching %s: %v", ticker, err)
			} else {
				data.Info = info
			}
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func fetchStock(ticker string) (*StockInfo, error) {
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v7/finance/quote?symbols=%s", ticker)
	
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var yahooResp YahooResponse
	if err := json.Unmarshal(body, &yahooResp); err != nil {
		return nil, err
	}
	
	if len(yahooResp.QuoteResponse.Result) == 0 {
		return nil, fmt.Errorf("no data found for ticker %s", ticker)
	}
	
	quote := yahooResp.QuoteResponse.Result[0]
	recommendation := generateRecommendation(quote.RegularMarketChangePercent)
	
	return &StockInfo{
		Ticker:         quote.Symbol,
		CompanyName:    getCompanyName(quote.LongName, quote.ShortName, quote.Symbol),
		Price:          quote.RegularMarketPrice,
		Change:         quote.RegularMarketChange,
		ChangePct:      fmt.Sprintf("%.2f%%", quote.RegularMarketChangePercent),
		Open:           quote.RegularMarketOpen,
		High:           quote.RegularMarketDayHigh,
		Low:            quote.RegularMarketDayLow,
		Volume:         formatVolume(quote.RegularMarketVolume),
		MarketCap:      formatMarketCap(quote.MarketCap),
		Recommendation: recommendation,
	}, nil
}

func generateRecommendation(changePct float64) string {
	if changePct > 5 {
		return "SELL - High gains, consider taking profits"
	} else if changePct < -5 {
		return "BUY - Significant dip, potential opportunity"
	} else if changePct > 2 {
		return "HOLD - Positive momentum"
	} else if changePct < -2 {
		return "HOLD - Minor decline, watch closely"
	}
	return "HOLD - Neutral conditions"
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

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

const fallbackTemplate = `<!DOCTYPE html>
<html>
<head>
    <title>MoneyMaker</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background: #f5f5f5; }
        .container { max-width: 600px; margin: 0 auto; background: white; padding: 30px; border-radius: 10px; }
        .error { background: #ffebee; color: #c62828; padding: 15px; border-radius: 5px; margin: 15px 0; }
        input { width: 100%; padding: 12px; font-size: 16px; border: 1px solid #ddd; border-radius: 5px; margin: 10px 0; }
        button { background: #1976d2; color: white; padding: 12px 24px; border: none; border-radius: 5px; font-size: 16px; cursor: pointer; }
        button:hover { background: #1565c0; }
        .stock-info { background: #f8f9fa; padding: 20px; border-radius: 5px; margin: 20px 0; }
        .positive { color: #2e7d32; }
        .negative { color: #d32f2f; }
    </style>
</head>
<body>
    <div class="container">
        <h1>💰 MoneyMaker</h1>
        <p>Simple Stock Analysis</p>
        
        {{if .Error}}
        <div class="error">{{.Error}}</div>
        {{end}}
        
        <form method="POST">
            <input type="text" name="ticker" placeholder="Enter stock ticker (e.g., AAPL, TSLA)" value="{{.LastQuery}}" required>
            <button type="submit">Get Stock Data</button>
        </form>
        
        {{if .Info}}
        <div class="stock-info">
            <h2>{{.Info.CompanyName}} ({{.Info.Ticker}})</h2>
            <p><strong>Price:</strong> ${{printf "%.2f" .Info.Price}}</p>
            <p><strong>Change:</strong> <span class="{{if ge .Info.Change 0}}positive{{else}}negative{{end}}">{{if ge .Info.Change 0}}+{{end}}${{printf "%.2f" .Info.Change}} ({{.Info.ChangePct}})</span></p>
            <p><strong>Open:</strong> ${{printf "%.2f" .Info.Open}}</p>
            <p><strong>High:</strong> ${{printf "%.2f" .Info.High}}</p>
            <p><strong>Low:</strong> ${{printf "%.2f" .Info.Low}}</p>
            <p><strong>Volume:</strong> {{.Info.Volume}}</p>
            <p><strong>Market Cap:</strong> {{.Info.MarketCap}}</p>
            <p><strong>Recommendation:</strong> {{.Info.Recommendation}}</p>
        </div>
        {{end}}
    </div>
</body>
</html>`
