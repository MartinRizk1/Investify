# Money Maker - Stock Analysis Web App

A modern, functional stock analysis web application that provides stock data and analysis using public Yahoo Finance APIs.

## Features

- **Company Search**: Enter company names (like "Apple", "Tesla") or ticker symbols (like "AAPL", "TSLA").
- **Real-time Data**: Fetches live stock data from Yahoo Finance.
- **Interactive Charts**: Displays stock price history with Chart.js.
- **Company Information**: Shows sector, industry, market cap, website, and business summary.
- **Modern UI**: Glossy black design with gradient accents and responsive layout.

## Installation

1. Clone or download the project:
   ```bash
   git clone https://github.com/yourusername/MoneyMaker.git
   cd MoneyMaker
   ```
2. Install Go (if not already installed):
   - Follow the instructions at [https://golang.org/doc/install](https://golang.org/doc/install).
3. Build the application:
   ```bash
   go build -o moneymaker main.go
   ```

## Usage

1. Run the application:
   ```bash
   ./moneymaker
   ```
2. Open your browser and go to `http://localhost:8080`.
3. Enter a company name or ticker symbol in the search bar.
4. View the analysis results, charts, and company information.

## Supported Companies

The app works with any publicly traded company, fetching data directly from Yahoo Finance.

## Technical Details

- **Backend**: Go web server using the `net/http` package.
- **Data Source**: Yahoo Finance API via public endpoints.
- **Frontend**: HTML/CSS/JavaScript with Chart.js.
- **Styling**: Modern gradient-based design with Orbitron font.

## Architecture

```
MoneyMaker/
├── main.go               # Main Go application
├── templates/
│   ├── index.html       # Frontend UI
│   └── clean.html       # Additional template
└── static/              # Static assets (CSS, JS, images)
```

## License

This project is for educational and personal use.
