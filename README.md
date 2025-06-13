# Investify - Advanced Stock Analysis Application

Investify is a Go-based stock analysis application that provides real-time stock data, AI-powered recommendations, and TensorFlow-based price predictions.

## Features

- **Real-time Stock Data**: Fetch up-to-date stock information from the Yahoo Finance API
- **AI Analysis**: Get AI-powered stock recommendations and insights
- **TensorFlow Predictions**: Use machine learning models to predict future stock prices
- **Beautiful UI**: Modern, responsive interface with charts and animations
- **Go-Python Bridge**: Seamless integration between Go backend and Python TensorFlow models
- **Security-First Design**: Input validation, CSP security headers, and secure API handling

## Architecture

The application has the following components:
- **Go Backend**: Handles HTTP requests, stock data fetching, and business logic
- **Python ML Layer**: TensorFlow-based prediction models for stock price analysis
- **HTML/JS Frontend**: Responsive UI with charts and animations
- **Go-Python Bridge**: Communication layer between Go and Python

## Requirements

### Basic Requirements
- Go 1.16+
- Web browser

### Optional (for ML features)
- Python 3.8+
- TensorFlow 2.6+
- Other Python dependencies listed in `models/requirements.txt`

## Installation

### 1. Clone the repository
```bash
git clone https://github.com/martinrizk/investify.git
cd investify
```

### 2. Install Go dependencies
```bash
go mod download
```

### 3. Set up Python environment (optional)
```bash
make python-setup
```

### 4. Build the application
```bash
make build
```

### 5. Run the application
```bash
./investify
```

## Development

### Building the application
```bash
make build
```

### Running tests
```bash
make test
```

### Running with hot reload
```bash
make run
```

## Security

### Environment Variables
Investify uses environment variables for all sensitive information:

1. Copy the example environment file
```bash
cp .env.example .env
```

2. Edit the `.env` file to add your API keys:
```bash
# OpenAI API Key for AI-powered recommendations
OPENAI_API_KEY=your_openai_key_here
```

### Security Features
- Input validation for all user inputs
- Content Security Policy (CSP) headers
- Protected API access points
- No sensitive data leakage in responses or logs
- Built-in rate limiting for external API calls

## TensorFlow Integration

The application can work in two modes:

1. **Full Mode**: Using Python TensorFlow models for predictions
2. **Fallback Mode**: Using Go-based prediction simulations when Python/TensorFlow is not available

### Training TensorFlow Models

```bash
# Train a quick test model
make train-test-model

# Train models for popular stocks (takes longer)
make train-models
```

## API Documentation

### GET /health
Returns the health status of the application

### GET /
Serves the main application UI

### POST /
Accepts stock search requests and returns analysis

## Security

### API Keys
This application uses the following API keys that should be set in your environment:

- `OPENAI_API_KEY`: Required for AI-powered stock recommendations (optional, will use rule-based recommendations if not provided)
- `PORT`: Optional, defaults to 8080

Create a `.env` file in the root directory by copying the `.env.example` file and filling in your values:

```bash
cp .env.example .env
# Edit .env with your API keys
```

> ⚠️ **Warning**: Never commit your actual API keys to GitHub. The `.env` file is included in `.gitignore` to prevent accidental commits.

## GitHub Repository Security

When pushing to GitHub, ensure:

1. No sensitive environment variables are committed:
   ```bash
   # Check for any potentially sensitive files before commit
   git diff --cached | grep -i "key\|secret\|password\|credential"
   ```

2. The `.env.example` file is kept up to date but does not contain real credentials

3. Your `.gitignore` file is excluding all sensitive files:
   ```
   # Sensitive files
   .env
   .env.local
   *.pem
   *.key
   ```

4. Regularly scan your repository for sensitive information using tools like:
   - GitHub's secret scanning feature
   - Git-secrets (https://github.com/awslabs/git-secrets)

5. Consider enabling branch protection rules in the GitHub repository settings

## License

MIT License