package handlers

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/martinrizk/investify/internal/services"
)

// StockUpdate represents a stock price update that is sent over websocket
type StockUpdate struct {
	Ticker        string                 `json:"ticker"`
	Price         float64                `json:"price"`
	Change        float64                `json:"change"`
	ChangePct     string                 `json:"change_pct"`
	LastUpdated   string                 `json:"last_updated"`
	Technical     map[string]interface{} `json:"technical,omitempty"`
}

// formatTime returns a formatted time string
// If t is nil, the current time is used
func formatTime(t *time.Time) string {
	if t == nil {
		now := time.Now()
		t = &now
	}
	return t.Format(time.RFC3339)
}

// StockPrice represents the price data for a stock
type StockPrice struct {
	Price     float64 `json:"price"`
	Change    float64 `json:"change"`
	ChangePct string  `json:"change_pct"`
}

// Global variables
var (
	pythonBridge *services.PythonBridge
)

// InitWebSocketHandler initializes required services for WebSocket handler
func InitWebSocketHandler() {
	pythonBridge = services.GetPythonBridge()
	rand.Seed(time.Now().UnixNano()) // Initialize random seed
}

var (
	// Websocket upgrader with CORS support
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for development, restrict in production
		},
	}

	// Store active connections
	clients      = make(map[*websocket.Conn]string) // map[connection]ticker
	clientsMutex sync.Mutex
)

// HandleWebSocket upgrades an HTTP connection to WebSocket
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Extract ticker from URL (format: /ws/stocks/{ticker})
	ticker := r.URL.Path[len("/ws/stocks/"):]
	if ticker == "" {
		http.Error(w, "Ticker symbol required", http.StatusBadRequest)
		return
	}

	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to websocket: %v", err)
		return
	}

	// Register client
	clientsMutex.Lock()
	clients[conn] = ticker
	clientsMutex.Unlock()

	log.Printf("WebSocket client connected for ticker: %s", ticker)

	// Start goroutine to handle WebSocket connection
	go handleConnection(conn, ticker)
}

// handleConnection processes messages from the WebSocket connection
func handleConnection(conn *websocket.Conn, ticker string) {
	defer func() {
		// Unregister client on disconnect
		clientsMutex.Lock()
		delete(clients, conn)
		clientsMutex.Unlock()
		
		conn.Close()
		log.Printf("WebSocket client disconnected for ticker: %s", ticker)
	}()

	// Send initial update
	sendStockUpdate(conn, ticker)

	// Handle WebSocket messages (not used yet, but could be used for client requests)
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
	}
}

// sendStockUpdate sends stock data to the WebSocket client
func sendStockUpdate(conn *websocket.Conn, ticker string) {
	// Get stock price data
	stockPrice, err := fetchRealtimePrice(ticker)
	if err != nil {
		log.Printf("Error fetching price for %s: %v", ticker, err)
		return
	}

	// Get technical indicators
	technical, err := fetchTechnicalIndicators(ticker)
	if err != nil {
		log.Printf("Error fetching technical indicators for %s: %v", ticker, err)
		// Continue anyway, just without technical data
	}

	update := StockUpdate{
		Ticker:      ticker,
		Price:       stockPrice.Price,
		Change:      stockPrice.Change,
		ChangePct:   stockPrice.ChangePct,
		LastUpdated: formatTime(nil),
		Technical:   technical,
	}

	// Send the update
	if err := conn.WriteJSON(update); err != nil {
		log.Printf("Error sending stock update: %v", err)
	}
}

// broadcastPriceUpdates periodically sends price updates to all connected clients
func StartPriceUpdateBroadcaster() {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker.C {
			clientsMutex.Lock()
			for conn, symbol := range clients {
				go sendStockUpdate(conn, symbol)
			}
			clientsMutex.Unlock()
		}
	}()
}

// fetchRealtimePrice gets the real-time price for a stock ticker
// In a production app, this would call a financial API
func fetchRealtimePrice(ticker string) (*StockPrice, error) {
	// In a real app, we'd fetch from an API here
	// For now, simulate a slightly random price based on ticker
	
	// Base price depends on ticker for variety
	var basePrice float64
	var change float64
	
	switch ticker {
	case "AAPL":
		basePrice = 180.0 + (rand.Float64() * 5.0 - 2.5)
		change = 0.75
	case "MSFT":
		basePrice = 350.0 + (rand.Float64() * 7.5 - 3.75)
		change = 1.2
	case "GOOGL":
		basePrice = 130.0 + (rand.Float64() * 4.0 - 2.0)
		change = 0.5
	case "AMZN":
		basePrice = 140.0 + (rand.Float64() * 5.0 - 2.5)
		change = 0.65
	case "META":
		basePrice = 310.0 + (rand.Float64() * 6.0 - 3.0)
		change = 1.1
	case "TSLA":
		basePrice = 230.0 + (rand.Float64() * 10.0 - 5.0)
		change = 2.5
	default:
		basePrice = 100.0 + (rand.Float64() * 25.0)
		change = 0.5
	}
	
	// Add small random variation to simulate real-time changes
	priceChange := (rand.Float64() - 0.5) * 0.5 // Random value between -0.25 and +0.25
	newPrice := basePrice + priceChange
	
	// Calculate new change
	newChange := change + priceChange
	newChangePct := fmt.Sprintf("%.2f%%", (newChange/newPrice)*100)
	
	return &StockPrice{
		Price:     newPrice,
		Change:    newChange,
		ChangePct: newChangePct,
	}, nil
}

// fetchTechnicalIndicators gets technical indicator data for a ticker
// In a production app, this would call our Python analyzer
func fetchTechnicalIndicators(ticker string) (map[string]interface{}, error) {
	// In a real app, call the Python analyzer here
	
	// Initialize the Python bridge if needed
	if pythonBridge == nil {
		pythonBridge = services.GetPythonBridge()
	}
	
	// Try to get technical data from Python analyzer
	result, err := pythonBridge.PredictStockPriceWithSimpleAnalyzer(ticker)
	if err == nil && result != nil && result.Technical != nil {
		return result.Technical, nil
	}
	
	// If that fails, return simulated data
	return createSimulatedTechnicalData(), nil
}

// Create simulated technical indicator data
func createSimulatedTechnicalData() map[string]interface{} {
	// Generate dates for the last 20 days
	dates := make([]string, 20)
	for i := 19; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		dates[19-i] = date.Format("2006-01-02")
	}
	
	// Generate RSI data (oscillates between 30-70)
	rsiData := make([]float64, 20)
	for i := 0; i < 20; i++ {
		value := 50.0 + math.Sin(float64(i)*0.5)*20.0
		rsiData[i] = math.Round(value*10) / 10
	}
	
	// Generate MACD data
	macdLine := make([]float64, 20)
	signalLine := make([]float64, 20)
	histogram := make([]float64, 20)
	
	for i := 0; i < 20; i++ {
		macd := math.Sin(float64(i)*0.3) * 2.0
		signal := math.Sin((float64(i)-3.0)*0.3) * 2.0
		
		macdLine[i] = math.Round(macd*100) / 100
		signalLine[i] = math.Round(signal*100) / 100
		histogram[i] = math.Round((macd-signal)*100) / 100
	}
	
	// Generate Bollinger Bands data
	middleBand := make([]float64, 20)
	upperBand := make([]float64, 20)
	lowerBand := make([]float64, 20)
	basePrice := 150.0
	
	for i := 0; i < 20; i++ {
		price := basePrice + math.Sin(float64(i)*0.3)*15.0
		volatility := 5.0 + math.Abs(math.Sin(float64(i)*0.5))*10.0
		
		middleBand[i] = math.Round(price*100) / 100
		upperBand[i] = math.Round((price+volatility)*100) / 100
		lowerBand[i] = math.Round((price-volatility)*100) / 100
	}
	
	return map[string]interface{}{
		"dates":            dates,
		"rsi":              rsiData,
		"macd":             macdLine,
		"macd_signal":      signalLine,
		"macd_histogram":   histogram,
		"bollinger_middle": middleBand,
		"bollinger_upper":  upperBand,
		"bollinger_lower":  lowerBand,
	}
}
