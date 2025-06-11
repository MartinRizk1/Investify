# 🚀 Investify - AI-Powered Stock Analysis Platform

A sophisticated Go-based web application that provides real-time stock information, interactive price charts, and dual AI-powered investment recommendations using Yahoo Finance API data. The platform combines OpenAI's language model for qualitative analysis and a TensorFlow-simulated model for quantitative predictions.

## Features

- **Real-time Stock Data**: Fetches live stock prices, market metrics, and company information from Yahoo Finance
- **Dual AI Analysis**:
  - OpenAI GPT for qualitative market insights and recommendations
  - TensorFlow-based model for quantitative price predictions and trend analysis
- **Interactive Charts**: Visual price movement representation with day's range
- **Beautiful UI**: Modern, responsive design with gradient themes and intuitive layout
- **Smart Search**: Company name mapping to correct ticker symbols
- **Technical Indicators**: Price position analysis, volatility metrics, and momentum tracking

## Prerequisites

- Go 1.20 or higher
- (Optional) OpenAI API key for advanced AI recommendations

## Installation & Setup

1. Clone the repository
```bash
git clone https://github.com/yourusername/investify.git
cd investify
```

2. Install dependencies
```bash
make deps
```

3. Build the application
```bash
make build
```

4. Set OpenAI API key (optional, for enhanced AI recommendations)
```bash
export OPENAI_API_KEY=your_api_key_here
```

5. Run the application
```bash
./investify
# Or use: make run
```

6. Access the application
Open your browser and navigate to http://localhost:8080

## Running Tests

Run all tests:
```bash
make test
```

Run tests with coverage analysis:
```bash
make test-coverage
```

## Development

Format code:
```bash
make fmt
```

Run static analysis:
```bash
make vet
```

## Error Handling

The application implements robust error handling for:
- Network failures with Yahoo Finance API
- Rate limiting scenarios with exponential backoff
- OpenAI API connectivity issues with fallback recommendations
- Invalid stock data with auto-repair mechanisms

1. **Clone the repository:**

   ```bash
   git clone <repository-url>
   cd Investify
   ```

2. **Install dependencies:**

   ```bash
   go mod tidy
   ```

3. **Set up OpenAI API Key (Optional):**

   ```bash
   export OPENAI_API_KEY="your-openai-api-key-here"
   ```

   If you don't have an OpenAI API key, the application will still work with rule-based recommendations.

4. **Run the application:**

   ```bash
   go run cmd/main.go
   ```

5. **Open your browser:**
   Navigate to `http://localhost:8080`

## Usage

1. Enter a stock ticker symbol (e.g., AAPL, TSLA, GOOGL, MSFT)
2. Click "🔍 Analyze Stock"
3. View comprehensive stock information including:
   - Current price and daily changes
   - Market metrics (open, high, low, volume, market cap)
   - Interactive price chart
   - AI-powered investment recommendation

## API Keys Required

### OpenAI API Key (Optional but Recommended)

To get advanced AI-powered recommendations, you'll need an OpenAI API key:

1. Go to [OpenAI API](https://platform.openai.com/api-keys)
2. Create an account or sign in
3. Generate a new API key
4. Set it as an environment variable:
   ```bash
   export OPENAI_API_KEY="sk-your-key-here"
   ```

**Note**: Without an OpenAI API key, the application will use rule-based recommendations that analyze:

- Price momentum and trends
- Daily price position within the trading range
- Volume and volatility indicators

## Project Structure

```
Investify/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── handlers/
│   │   └── stock_handler.go    # HTTP handlers
│   └── services/
│       ├── stock_service.go    # Yahoo Finance API integration
│       └── ai_service.go       # AI recommendation service
├── templates/
│   └── index.html              # Web interface
├── go.mod                      # Go module definition
└── README.md                   # This file
```

## Technologies Used

- **Backend**: Go with Gorilla Mux router
- **Frontend**: HTML5, CSS3, Chart.js for visualizations
- **APIs**:
  - Yahoo Finance API (free, no key required)
  - OpenAI API (optional, requires key)
- **Styling**: Modern CSS with gradients and animations

## Example Stocks to Try

- **AAPL** - Apple Inc.
- **TSLA** - Tesla Inc.
- **GOOGL** - Alphabet Inc.
- **MSFT** - Microsoft Corporation
- **AMZN** - Amazon.com Inc.
- **NVDA** - NVIDIA Corporation

## Environment Variables

| Variable         | Description                           | Required |
| ---------------- | ------------------------------------- | -------- |
| `OPENAI_API_KEY` | OpenAI API key for AI recommendations | No       |
| `PORT`           | Server port (default: 8080)           | No       |

## Features in Detail

### Stock Data

- Real-time price information from Yahoo Finance
- Market metrics including volume and market capitalization
- Daily open, high, low, and current prices
- Percentage and dollar change calculations

### AI Recommendations

With OpenAI API key:

- GPT-powered analysis considering market conditions
- Contextual recommendations based on current data
- Intelligent reasoning for investment decisions

Without OpenAI API key:

- Rule-based technical analysis
- Momentum and trend evaluation
- Price position analysis within daily range

### Visual Charts

- Interactive line charts showing price movement
- Color-coded for gains (green) and losses (red)
- Responsive design for all device sizes

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## Disclaimer

This application is for educational and informational purposes only. Stock market investments carry risk, and past performance does not guarantee future results. Always consult with a qualified financial advisor before making investment decisions.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
