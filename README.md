# 🚀 Investify - AI-Powered Stock Analysis Platform

A sophisticated Go-based web application that provides real-time stock information, interactive price charts, and dual AI-powered investment recommendations using Yahoo Finance API data. The platform combines OpenAI's language model for qualitative analysis and a TensorFlow model for quantitative predictions.

## Features

- **Real-time Stock Data**: Fetches live stock prices, market metrics, and company information from Yahoo Finance
- **Dual AI Analysis**:
  - OpenAI GPT for qualitative market insights and recommendations
  - TensorFlow LSTM models for accurate price predictions and trend analysis
- **Interactive Charts**: Visual price movement representation with day's range and historical data
- **Beautiful UI**: Modern, responsive design with gradient themes and intuitive layout
- **Smart Search**: Company name mapping to correct ticker symbols
- **Technical Indicators**: Price position analysis, volatility metrics, and momentum tracking
- **AI Confidence Meter**: Visual representation of prediction confidence with key factors
- **Go-Python Bridge**: Seamless integration between Go backend and Python TensorFlow models

## Prerequisites

- Go 1.20 or higher
- Python 3.7+ (for TensorFlow models)
- (Optional) OpenAI API key for advanced AI recommendations

## Installation & Setup

1. Clone the repository
```bash
git clone https://github.com/yourusername/investify.git
cd investify
```

2. Install Go dependencies
```bash
make deps
```

3. Install Python dependencies for TensorFlow models
```bash
make python-setup
```

4. Set OpenAI API key (optional, for enhanced AI recommendations)
```bash
export OPENAI_API_KEY=your_api_key_here
```

5. Build the application
```bash
make build
```

6. Run the application
```bash
make run
```

7. Access the application in your browser
```
http://localhost:8080
```

## Advanced: Training TensorFlow Models

1. Train models for specific stocks:
```bash
make train-models
```

2. Or train a model for a specific ticker:
```bash
cd models && python3 train_stock_model.py TICKER
```

3. The trained models will be stored in `models/saved/` directory and automatically used by the application.

## Running Tests

```bash
make test
```

For coverage report:
```bash
make test-coverage
```

## Development

For hot-reloading during development:
```bash
# Install air first if you don't have it
go install github.com/cosmtrek/air@latest

# Run with hot reloading
air
```

## Error Handling

The application implements robust error handling for:
- Network failures with Yahoo Finance API
- Rate limiting scenarios with exponential backoff
- OpenAI API connectivity issues with fallback recommendations
- TensorFlow model prediction failures with rule-based alternatives
- Invalid stock data with auto-repair mechanisms

## Usage

1. Enter a company name or ticker symbol in the search field
2. View comprehensive stock analysis including:
   - Current price and day's movement
   - Key market metrics
   - AI-powered price prediction
   - Technical analysis
   - Investment recommendation
   - Historical price chart
   - Confidence meter with key factors

## API Keys Required

For full functionality, set your OpenAI API key:
```bash
export OPENAI_API_KEY="your-key-here"
```

The application will still work without the key, using rule-based recommendations instead of AI-powered ones.

## Project Structure

```
investify/
├── cmd/
│   └── main.go           # Application entry point
├── internal/
│   ├── handlers/         # HTTP request handlers
│   └── services/         # Business logic and data services
│       ├── ai_service.go # OpenAI API integration
│       ├── py_bridge.go  # Go-Python bridge
│       ├── stock_service.go # Yahoo Finance API integration
│       └── tensorflow_model.go # TensorFlow model integration
├── models/
│   ├── predict.py        # Python script for predictions
│   ├── stock_predict.py  # Advanced prediction model
│   ├── requirements.txt  # Python dependencies
│   └── train_stock_model.py # Model training script
└── templates/            # HTML templates and static assets
```

## Technologies Used

- **Backend**: Go, Gorilla Mux
- **APIs**: Yahoo Finance, OpenAI
- **Machine Learning**: TensorFlow (Python)
- **Frontend**: HTML, CSS, JavaScript, Chart.js

## Example Stocks to Try

- AAPL (Apple)
- GOOGL (Alphabet/Google)
- MSFT (Microsoft)
- AMZN (Amazon)
- TSLA (Tesla)

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| OPENAI_API_KEY | OpenAI API key | No (optional) |
| PORT | Server port | No (defaults to 8080) |

## Features in Detail

### Real-time Stock Data
- Current price and day's change
- Trading volume and market cap
- Day's range with position indicator
- 52-week high and low

### AI Analysis
- Predicted price based on machine learning models
- Confidence rating with visual meter
- Trend direction with visual indicators
- Key factors influencing the prediction

### Visual Charts
- Interactive price chart
- Historical data visualization
- Prediction markers
- Technical indicators

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Disclaimer

This application is for educational purposes only. The AI recommendations and price predictions should not be considered financial advice. Always do your own research before making investment decisions.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
