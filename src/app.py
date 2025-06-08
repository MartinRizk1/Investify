import sys
import os
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

import re
from functools import wraps
from flask import Flask, render_template, request, jsonify
from flask_limiter import Limiter
from flask_limiter.util import get_remote_address
from flask_wtf.csrf import CSRFProtect
from dotenv import load_dotenv
from src.data.data_fetcher import DataFetcher
from src.analyzers.market_analyzer import MarketAnalyzer
from src.utils.logging_config import setup_security_logging

# Load environment variables
load_dotenv()

app = Flask(__name__)

# Security configuration
app.config['SECRET_KEY'] = os.getenv('SECRET_KEY', 'dev-key-change-in-production')
app.config['WTF_CSRF_TIME_LIMIT'] = int(os.getenv('CSRF_TIME_LIMIT', '3600'))

# Initialize security components
csrf = CSRFProtect(app)
limiter = Limiter(
    key_func=get_remote_address,
    app=app,
    default_limits=[os.getenv('RATE_LIMIT_DEFAULT', '100 per hour')]
)

# Setup security logging
logger = setup_security_logging()

# Security headers middleware
@app.after_request
def add_security_headers(response):
    """Add security headers to all responses"""
    response.headers['X-Content-Type-Options'] = 'nosniff'
    response.headers['X-Frame-Options'] = 'DENY'
    response.headers['X-XSS-Protection'] = '1; mode=block'
    response.headers['Strict-Transport-Security'] = 'max-age=31536000; includeSubDomains'
    response.headers['Content-Security-Policy'] = "default-src 'self'; script-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; img-src 'self' data:"
    return response

# Input validation decorator
def validate_input(validation_func):
    """Decorator for input validation"""
    def decorator(f):
        @wraps(f)
        def decorated_function(*args, **kwargs):
            try:
                validation_func()
                return f(*args, **kwargs)
            except ValueError as e:
                logger.warning(f"Input validation failed: {str(e)} - IP: {get_remote_address()}")
                return jsonify({'error': 'Invalid input provided'}), 400
        return decorated_function
    return decorator

def validate_ticker_input():
    """Validate ticker symbol input for POST requests"""
    if request.method == 'POST':
        company = request.form.get('company', '').strip()
        if not company:
            raise ValueError("Company/ticker is required")
        # Allow alphanumeric characters, dots, hyphens, and spaces (max 50 chars)
        if not re.match(r'^[a-zA-Z0-9.\s-]{1,50}$', company):
            raise ValueError("Invalid company/ticker format")
    return True

def validate_api_ticker_input():
    """Validate ticker symbol input for API endpoints"""
    company = request.view_args.get('ticker', '').strip() if request.view_args else ''
    if not company:
        raise ValueError("Company/ticker is required")
    # Allow alphanumeric characters, dots, hyphens, and spaces (max 50 chars)
    if not re.match(r'^[a-zA-Z0-9.\s-]{1,50}$', company):
        raise ValueError("Invalid company/ticker format")
    return company

def validate_period_input():
    """Validate period parameter"""
    period = request.view_args.get('period', '').strip()
    valid_periods = ['1M', '3M', '6M', '1Y', '5Y']
    
    if period not in valid_periods:
        raise ValueError(f"Invalid period. Must be one of: {', '.join(valid_periods)}")
    
    return period

@app.route('/', methods=['GET', 'POST'])
@limiter.limit("30 per minute")
def index():
    info = None
    action = None
    error = None
    
    if request.method == 'POST':
        # Apply validation only for POST requests
        try:
            validate_ticker_input()
        except ValueError as e:
            logger.warning(f"Input validation failed: {str(e)} - IP: {get_remote_address()}")
            return jsonify({'error': 'Invalid input provided'}), 400
        
        company = request.form.get('company', '').strip()
        
        # Log the request for security monitoring
        logger.info(f"Stock analysis request for: {company[:20]}... - IP: {get_remote_address()}")
        
        if not company:
            error = 'Please enter a company name or ticker symbol.'
            return render_template('index.html', info=info, action=action, error=error)
            
        fetcher = DataFetcher()
        
        # Get the resolved ticker using enhanced search
        ticker = fetcher.get_ticker(company)
        
        if not ticker:
            error = f'Could not find ticker for "{company}". Please try:\n• Using the exact ticker symbol (e.g., "AAPL" for Apple)\n• A different spelling or the full company name\n• A well-known publicly traded company'
            logger.warning(f"Ticker not found for company: {company[:20]}... - IP: {get_remote_address()}")
            return render_template('index.html', info=info, action=action, error=error)
        
        # Fetch all data using the enhanced methods
        data = fetcher.fetch_data(company)
        company_info = fetcher.fetch_company_info(company)
        news = fetcher.fetch_news(company)
        
        # Check if we have sufficient data
        if data.empty and not company_info:
            error = f'No stock data found for "{company}" (ticker: {ticker}). This may be:\n• A private company\n• A delisted stock\n• An invalid ticker symbol\nPlease try a different publicly traded company.'
            logger.warning(f"No data found for: {company} ({ticker}) - IP: {get_remote_address()}")
        elif data.empty:
            error = f'Found company information for "{company}" ({ticker}) but no stock price data available. This stock may be suspended from trading.'
            logger.warning(f"No stock data for: {company} ({ticker}) - IP: {get_remote_address()}")
        else:
            try:
                analyzer = MarketAnalyzer()
                analyzer.train(data)
                action = analyzer.predict(data)
                
                logger.info(f"Successful analysis for: {company} ({ticker}) - Action: {action} - IP: {get_remote_address()}")
                
                # Prepare chart data (last 60 days or all available data)
                chart_days = min(60, len(data))
                chart_dates = [d.strftime('%Y-%m-%d') for d in data.index[-chart_days:]]
                chart_closes = []
                chart_volumes = []
                
                # Safe conversion for chart data
                for close_val in data['Close'].iloc[-chart_days:]:
                    safe_close = close_val.item() if hasattr(close_val, 'item') else float(close_val)
                    chart_closes.append(safe_close)
                
                for volume_val in data['Volume'].iloc[-chart_days:]:
                    safe_volume = volume_val.item() if hasattr(volume_val, 'item') else float(volume_val)
                    chart_volumes.append(safe_volume)
                
                chart_data = {
                    'dates': chart_dates,
                    'closes': chart_closes,
                    'volumes': chart_volumes
                }
                
                # Calculate additional metrics
                latest_data = data.iloc[-1]
                previous_data = data.iloc[-2] if len(data) > 1 else data.iloc[-1]
                
                price_change = latest_data['Close'] - previous_data['Close'] if len(data) > 1 else 0
                price_change_pct = (price_change / previous_data['Close'] * 100) if len(data) > 1 and previous_data['Close'] != 0 else 0
                
                # Calculate 52-week high/low
                fifty_two_week_high = data['High'].max()
                fifty_two_week_low = data['Low'].min()
                
                # Calculate daily high/low changes
                high_change = latest_data['High'] - previous_data['High'] if len(data) > 1 else 0
                low_change = latest_data['Low'] - previous_data['Low'] if len(data) > 1 else 0
                
                # Format market cap
                market_cap = company_info.get('marketCap')
                if market_cap:
                    if market_cap >= 1e12:
                        market_cap_formatted = f"${market_cap/1e12:.2f}T"
                    elif market_cap >= 1e9:
                        market_cap_formatted = f"${market_cap/1e9:.2f}B"
                    elif market_cap >= 1e6:
                        market_cap_formatted = f"${market_cap/1e6:.2f}M"
                    else:
                        market_cap_formatted = f"${market_cap:,.0f}"
                else:
                    market_cap_formatted = "N/A"
                
                # Safe conversion for pandas Series values
                last_close = latest_data['Close']
                last_volume = latest_data['Volume']
                last_open = latest_data['Open']
                last_high = latest_data['High']
                last_low = latest_data['Low']
                
                # Convert to float using the same safe method as in helpers.py
                last_close_val = last_close.item() if hasattr(last_close, 'item') else float(last_close)
                last_volume_val = last_volume.item() if hasattr(last_volume, 'item') else int(last_volume)
                last_open_val = last_open.item() if hasattr(last_open, 'item') else float(last_open)
                last_high_val = last_high.item() if hasattr(last_high, 'item') else float(last_high)
                last_low_val = last_low.item() if hasattr(last_low, 'item') else float(last_low)
                price_change_val = price_change.item() if hasattr(price_change, 'item') else float(price_change)
                high_change_val = high_change.item() if hasattr(high_change, 'item') else float(high_change)
                low_change_val = low_change.item() if hasattr(low_change, 'item') else float(low_change)
                fifty_two_week_high_val = fifty_two_week_high.item() if hasattr(fifty_two_week_high, 'item') else float(fifty_two_week_high)
                fifty_two_week_low_val = fifty_two_week_low.item() if hasattr(fifty_two_week_low, 'item') else float(fifty_two_week_low)
                
                # Get financial metrics with safe conversions
                pe_ratio = company_info.get('trailingPE', 'N/A')
                if pe_ratio != 'N/A' and pe_ratio:
                    pe_ratio = f"{pe_ratio:.2f}"
                
                forward_pe = company_info.get('forwardPE', 'N/A')  
                if forward_pe != 'N/A' and forward_pe:
                    forward_pe = f"{forward_pe:.2f}"
                
                dividend_yield = company_info.get('dividendYield')
                if dividend_yield:
                    dividend_yield = f"{dividend_yield * 100:.2f}%"
                else:
                    dividend_yield = "N/A"
                
                info = {
                    'company_name': company_info.get('shortName') or company_info.get('longName') or company.upper(),
                    'ticker': ticker,
                    'last_close': f"${last_close_val:.2f}",
                    'open': f"${last_open_val:.2f}",
                    'high': f"${last_high_val:.2f}",
                    'low': f"${last_low_val:.2f}",
                    'fifty_two_week_high': f"${fifty_two_week_high_val:.2f}",
                    'fifty_two_week_low': f"${fifty_two_week_low_val:.2f}",
                    'price_change': f"${price_change_val:+.2f}",
                    'price_change_pct': f"{price_change_pct:+.2f}%",
                    'high_change': f"${high_change_val:+.2f}" if abs(high_change_val) > 0.01 else None,
                    'low_change': f"${low_change_val:+.2f}" if abs(low_change_val) > 0.01 else None,
                    'last_volume': f"{last_volume_val:,.0f}",
                    'last_date': data.index[-1].strftime('%Y-%m-%d'),
                    'chart_data': chart_data,
                    'sector': company_info.get('sector', 'N/A'),
                    'industry': company_info.get('industry', 'N/A'),
                    'website': company_info.get('website'),
                    'summary': company_info.get('longBusinessSummary', 'No summary available.')[:500] + '...' if company_info.get('longBusinessSummary') and len(company_info.get('longBusinessSummary', '')) > 500 else company_info.get('longBusinessSummary', 'No summary available.'),
                    'market_cap': market_cap_formatted,
                    'pe_ratio': pe_ratio,
                    'forward_pe': forward_pe,
                    'dividend_yield': dividend_yield,
                    'eps': company_info.get('trailingEps', 'N/A'),
                    'revenue': company_info.get('totalRevenue'),
                    'news': news[:5] if news else []
                }
                
                # Format revenue
                if info['revenue']:
                    revenue = info['revenue']
                    if revenue >= 1e12:
                        info['revenue'] = f"${revenue/1e12:.2f}T"
                    elif revenue >= 1e9:
                        info['revenue'] = f"${revenue/1e9:.2f}B"
                    elif revenue >= 1e6:
                        info['revenue'] = f"${revenue/1e6:.2f}M"
                    else:
                        info['revenue'] = f"${revenue:,.0f}"
                else:
                    info['revenue'] = "N/A"
                    
            except Exception as e:
                error = f'Error analyzing data for "{company}" ({ticker}): {str(e)}'
                logger.error(f"Analysis error for {company}: {str(e)} - IP: {get_remote_address()}")
                
    return render_template('index.html', info=info, action=action, error=error)

@app.route('/api/chart/<ticker>/<period>')
@limiter.limit("60 per minute")
@validate_input(lambda: (validate_api_ticker_input(), validate_period_input()))
def get_chart_data(ticker, period):
    """API endpoint to fetch chart data for different time periods"""
    try:
        # Log API request for monitoring
        logger.info(f"Chart API request - Ticker: {ticker}, Period: {period} - IP: {get_remote_address()}")
        
        fetcher = DataFetcher()
        
        # Validate and get actual ticker
        actual_ticker = fetcher.get_ticker(ticker)
        if not actual_ticker:
            logger.warning(f"Invalid ticker in API request: {ticker} - IP: {get_remote_address()}")
            return jsonify({'error': f'Invalid ticker: {ticker}'}), 400
        
        # Map frontend periods to yfinance periods
        period_mapping = {
            '1M': '1mo',
            '3M': '3mo', 
            '6M': '6mo',
            '1Y': '1y',
            '5Y': '5y'
        }
        
        yf_period = period_mapping.get(period, '1y')
        data = fetcher.fetch_data(actual_ticker, period=yf_period)
        
        if data.empty:
            logger.warning(f"No chart data available for {ticker} ({period}) - IP: {get_remote_address()}")
            return jsonify({'error': 'No data available'}), 404
            
        # Prepare chart data
        chart_dates = [d.strftime('%Y-%m-%d') for d in data.index]
        chart_closes = []
        chart_volumes = []
        
        # Safe conversion for chart data
        for close_val in data['Close']:
            safe_close = close_val.item() if hasattr(close_val, 'item') else float(close_val)
            chart_closes.append(safe_close)
        
        for volume_val in data['Volume']:
            safe_volume = volume_val.item() if hasattr(volume_val, 'item') else float(volume_val)
            chart_volumes.append(safe_volume)
        
        chart_data = {
            'dates': chart_dates,
            'closes': chart_closes,
            'volumes': chart_volumes,
            'period': period
        }
        
        return jsonify(chart_data)
        
    except Exception as e:
        logger.error(f"Chart API error for {ticker}: {str(e)} - IP: {get_remote_address()}")
        return jsonify({'error': str(e)}), 500

if __name__ == '__main__':
    # Load environment variables
    debug_mode = os.getenv('FLASK_DEBUG', 'False').lower() == 'true'
    port = int(os.getenv('FLASK_PORT', '5001'))
    host = os.getenv('FLASK_HOST', '127.0.0.1')
    
    # Log startup
    logger.info(f"Starting MoneyMaker application - Debug: {debug_mode}, Host: {host}, Port: {port}")
    
    # For production, set debug=False to prevent exposing sensitive information
    app.run(debug=debug_mode, port=port, host=host)
