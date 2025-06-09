# Money Maker - Stock Analysis Web App

A modern, functional stock analysis web application that provides buy/sell/hold recommendations using machine learning.

## Features

- **Company Search**: Enter company names (like "Apple", "Tesla") or ticker symbols (like "AAPL", "TSLA")
- **Smart Ticker Resolution**: Automatically converts company names to correct ticker symbols
- **Real-time Data**: Fetches live stock data from Yahoo Finance
- **ML Predictions**: Uses logistic regression to provide buy/sell/hold recommendations
- **Interactive Charts**: Displays stock price history with Chart.js
- **Company Information**: Shows sector, industry, market cap, website, and business summary
- **Latest News**: Displays recent news articles related to the company
- **Modern UI**: Glossy black design with gradient accents and responsive layout

## Installation

1. Clone or download the project
2. Create a virtual environment:
   ```bash
   python3 -m venv venv
   source venv/bin/activate  # On Windows: venv\Scripts\activate
   ```
3. Install dependencies:
   ```bash
   pip install -r requirements.txt
   ```

## Usage

1. Start the application:
   ```bash
   cd src
   python app.py
   ```
2. Open your browser and go to `http://localhost:5001`
3. Enter a company name or ticker symbol in the search bar
4. View the analysis results, charts, and recommendations

## Supported Companies

The app works with any publicly traded company, but includes built-in mappings for popular companies:

- Apple, Microsoft, Google/Alphabet, Amazon, Tesla, Meta/Facebook
- Netflix, NVIDIA, Salesforce, Adobe, Intel, Oracle, IBM
- Coca-Cola, Pepsi, Johnson & Johnson, Visa, Mastercard
- And many more...

## Technical Details

- **Backend**: Flask web framework
- **Data Source**: Yahoo Finance API via yfinance library
- **Machine Learning**: Scikit-learn logistic regression
- **Frontend**: HTML/CSS/JavaScript with Chart.js
- **Styling**: Modern gradient-based design with Orbitron font

## Architecture

```
src/
├── app.py                 # Main Flask application
├── data/
│   └── data_fetcher.py   # Yahoo Finance data retrieval
├── analyzers/
│   └── market_analyzer.py # ML prediction model
└── templates/
    └── index.html        # Frontend UI
```

## API

The application uses the yfinance library to fetch data from Yahoo Finance, including:

- Historical stock prices
- Company information (sector, industry, market cap)
- Recent news articles
- Trading volume and other metrics

## License

This project is for educational and personal use.
