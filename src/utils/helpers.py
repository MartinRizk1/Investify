import matplotlib.pyplot as plt
import matplotlib.dates as mdates
import pandas as pd
from datetime import datetime
import numpy as np

def print_trade_action(ticker: str, action: str, data: pd.DataFrame, company_info: dict):
    """Enhanced trade recommendation with detailed analysis."""
    current_price = data['Close'].iloc[-1]
    previous_price = data['Close'].iloc[-2] if len(data) > 1 else current_price
    
    # Convert to float to avoid pandas Series issues
    current_val = current_price.item() if hasattr(current_price, 'item') else float(current_price)
    previous_val = previous_price.item() if hasattr(previous_price, 'item') else float(previous_price)
    
    price_change = current_val - previous_val
    price_change_pct = (price_change / previous_val) * 100 if previous_val != 0 else 0
    
    # Calculate key metrics
    ma_20 = data['Close'].rolling(window=20).mean().iloc[-1] if len(data) >= 20 else current_price
    ma_50 = data['Close'].rolling(window=50).mean().iloc[-1] if len(data) >= 50 else current_price
    volume_avg = data['Volume'].rolling(window=20).mean().iloc[-1] if len(data) >= 20 else data['Volume'].iloc[-1]
    current_volume = data['Volume'].iloc[-1]
    
    # Convert to float
    ma_20_val = ma_20.item() if hasattr(ma_20, 'item') else float(ma_20)
    ma_50_val = ma_50.item() if hasattr(ma_50, 'item') else float(ma_50)
    volume_avg_val = volume_avg.item() if hasattr(volume_avg, 'item') else float(volume_avg)
    current_volume_val = current_volume.item() if hasattr(current_volume, 'item') else float(current_volume)
    
    # Volatility (standard deviation of returns)
    returns = data['Close'].pct_change().dropna()
    volatility = returns.std() * np.sqrt(252) * 100  # Annualized volatility in %
    volatility_val = volatility.item() if hasattr(volatility, 'item') else float(volatility)
    
    print(f"🎯 RECOMMENDATION: {action}")
    print(f"💰 Current Price: ${current_val:.2f}")
    print(f"📈 Price Change: ${price_change:.2f} ({price_change_pct:+.2f}%)")
    print(f"📊 20-day MA: ${ma_20_val:.2f}")
    print(f"📊 50-day MA: ${ma_50_val:.2f}")
    print(f"📦 Volume: {current_volume_val:,.0f} (Avg: {volume_avg_val:,.0f})")
    print(f"⚡ Volatility: {volatility_val:.1f}% (annualized)")
    
    # Analysis reasoning
    print(f"\n🧠 AI REASONING:")
    if action == "BUY":
        print("• Strong upward momentum detected")
        print("• Technical indicators suggest positive trend")
        print("• Consider this as a potential buying opportunity")
    elif action == "SELL":
        print("• Downward pressure identified")
        print("• Technical analysis suggests potential decline")
        print("• Consider reducing position or taking profits")
    else:  # HOLD
        print("• Mixed signals in the market data")
        print("• No clear directional bias detected")
        print("• Recommend maintaining current position")

def display_company_info(company_info: dict, ticker: str):
    """Display comprehensive company information."""
    if not company_info:
        print(f"❌ No company information available for {ticker}")
        return
    
    print(f"🏢 COMPANY: {company_info.get('shortName', ticker)}")
    print(f"🎫 Ticker: {ticker}")
    print(f"🏭 Sector: {company_info.get('sector', 'N/A')}")
    print(f"🔧 Industry: {company_info.get('industry', 'N/A')}")
    print(f"🌐 Website: {company_info.get('website', 'N/A')}")
    
    # Market cap formatting
    market_cap = company_info.get('marketCap')
    if market_cap:
        if market_cap >= 1e12:
            market_cap_str = f"${market_cap/1e12:.2f}T"
        elif market_cap >= 1e9:
            market_cap_str = f"${market_cap/1e9:.2f}B"
        elif market_cap >= 1e6:
            market_cap_str = f"${market_cap/1e6:.2f}M"
        else:
            market_cap_str = f"${market_cap:,.0f}"
        print(f"💎 Market Cap: {market_cap_str}")
    
    # Business summary
    summary = company_info.get('longBusinessSummary', '')
    if summary:
        print(f"\n📄 BUSINESS SUMMARY:")
        # Truncate long summaries
        if len(summary) > 300:
            summary = summary[:300] + "..."
        print(f"   {summary}")

def display_earnings_info(company_info: dict, data: pd.DataFrame):
    """Display earnings and financial metrics."""
    print(f"\n💰 FINANCIAL METRICS")
    print("-" * 30)
    
    if company_info:
        # Key financial metrics
        metrics = [
            ('P/E Ratio', company_info.get('trailingPE')),
            ('Forward P/E', company_info.get('forwardPE')),
            ('EPS (TTM)', company_info.get('trailingEps')),
            ('Revenue (TTM)', company_info.get('totalRevenue')),
            ('Profit Margin', company_info.get('profitMargins')),
            ('Return on Equity', company_info.get('returnOnEquity')),
            ('Debt to Equity', company_info.get('debtToEquity')),
            ('Dividend Yield', company_info.get('dividendYield')),
        ]
        
        for name, value in metrics:
            if value is not None:
                if name in ['Revenue (TTM)'] and isinstance(value, (int, float)):
                    if value >= 1e9:
                        formatted_value = f"${value/1e9:.2f}B"
                    elif value >= 1e6:
                        formatted_value = f"${value/1e6:.2f}M"
                    else:
                        formatted_value = f"${value:,.0f}"
                elif name in ['Profit Margin', 'Return on Equity', 'Dividend Yield'] and isinstance(value, (int, float)):
                    # Handle percentage values - yfinance returns these as decimals (e.g., 0.05 = 5%)
                    if name == 'Dividend Yield':
                        # Dividend yield: if value > 0.2 (20%), it's likely already a percentage
                        formatted_value = f"{value:.2f}%" if value > 0.2 else f"{value*100:.2f}%"
                    else:
                        # Other percentages: typically decimals that need to be converted
                        formatted_value = f"{value*100:.2f}%" if value <= 1 else f"{value:.2f}%"
                elif isinstance(value, (int, float)):
                    formatted_value = f"{value:.2f}"
                else:
                    formatted_value = str(value)
                print(f"📊 {name}: {formatted_value}")
    
    # Technical analysis from price data
    if not data.empty:
        print(f"\n📈 TECHNICAL ANALYSIS")
        print("-" * 30)
        
        # 52-week high/low  
        high_52w = data['High'].max()
        low_52w = data['Low'].min()
        current_price = data['Close'].iloc[-1]
        
        # Convert to float to avoid Series formatting issues
        high_val = high_52w.item() if hasattr(high_52w, 'item') else float(high_52w)
        low_val = low_52w.item() if hasattr(low_52w, 'item') else float(low_52w)
        current_val = current_price.item() if hasattr(current_price, 'item') else float(current_price)
        
        print(f"📊 52-Week High: ${high_val:.2f}")
        print(f"📊 52-Week Low: ${low_val:.2f}")
        print(f"📊 Current vs 52W High: {((current_val/high_val - 1) * 100):+.1f}%")
        print(f"📊 Current vs 52W Low: {((current_val/low_val - 1) * 100):+.1f}%")

def display_stock_graph(data: pd.DataFrame, ticker: str, company_name: str):
    """Display stock price graph or text summary."""
    if data.empty:
        print("❌ No data available for graph generation")
        return
    
    print(f"\n📈 STOCK CHART FOR {ticker}")
    print("-" * 40)
    
    # Always show text-based summary (more reliable than matplotlib in all environments)
    recent_data = data.tail(10)
    print(f"\n📊 Last 10 trading days for {ticker}:")
    print("-" * 60)
    print("Date          Close     Change    Volume")
    print("-" * 60)
    
    for i, (date, row) in enumerate(recent_data.iterrows()):
        close_price = row['Close'].item() if hasattr(row['Close'], 'item') else float(row['Close'])
        open_price = row['Open'].item() if hasattr(row['Open'], 'item') else float(row['Open'])
        volume = row['Volume'].item() if hasattr(row['Volume'], 'item') else int(row['Volume'])
        
        change = close_price - open_price
        change_pct = (change / open_price) * 100 if open_price != 0 else 0
        
        # Format volume
        if volume >= 1e9:
            vol_str = f"{volume/1e9:.1f}B"
        elif volume >= 1e6:
            vol_str = f"{volume/1e6:.1f}M"
        elif volume >= 1e3:
            vol_str = f"{volume/1e3:.1f}K"
        else:
            vol_str = f"{volume}"
            
        print(f"{date.strftime('%Y-%m-%d')}  ${close_price:7.2f}  {change_pct:+6.1f}%  {vol_str:>8}")
    
    print("-" * 60)
    
    # Price trend analysis
    current_price = data['Close'].iloc[-1]
    week_ago_price = data['Close'].iloc[-5] if len(data) >= 5 else current_price
    month_ago_price = data['Close'].iloc[-20] if len(data) >= 20 else current_price
    
    # Convert to float
    current_val = current_price.item() if hasattr(current_price, 'item') else float(current_price)
    week_ago_val = week_ago_price.item() if hasattr(week_ago_price, 'item') else float(week_ago_price)
    month_ago_val = month_ago_price.item() if hasattr(month_ago_price, 'item') else float(month_ago_price)
    
    week_change = ((current_val / week_ago_val) - 1) * 100 if week_ago_val != 0 else 0
    month_change = ((current_val / month_ago_val) - 1) * 100 if month_ago_val != 0 else 0
    
    print(f"\n📈 PRICE TRENDS:")
    print(f"   5-day change:  {week_change:+.1f}%")
    print(f"   20-day change: {month_change:+.1f}%")
    
    # Optional: Try to create matplotlib chart (with better error handling)
    try:
        create_matplotlib_chart(data, ticker, company_name)
    except Exception as e:
        print(f"\n📊 Chart visualization not available in this environment")
        print(f"   (Technical note: {str(e)[:50]}...)")

def create_matplotlib_chart(data: pd.DataFrame, ticker: str, company_name: str):
    """Create matplotlib chart with proper error handling."""
    import matplotlib
    matplotlib.use('Agg')  # Use non-interactive backend
    
    plt.style.use('dark_background')
    fig, ax = plt.subplots(figsize=(12, 6))
    
    # Convert to numpy arrays to avoid pandas Series issues
    dates = data.index.to_numpy()
    closes = data['Close'].to_numpy()
    
    # Plot stock price
    ax.plot(dates, closes, linewidth=2, color='#00ff99', label='Close Price')
    
    # Add moving averages if enough data
    if len(data) >= 20:
        ma20 = data['Close'].rolling(window=20).mean().to_numpy()
        ax.plot(dates, ma20, linewidth=1, color='#ffaa00', label='20-day MA', alpha=0.8)
    
    # Customize chart
    ax.set_title(f'{company_name} ({ticker}) - Stock Price', fontsize=14, color='white')
    ax.set_ylabel('Price ($)', fontsize=12, color='white')
    ax.legend()
    ax.grid(True, alpha=0.3)
    ax.tick_params(colors='white')
    
    plt.tight_layout()
    plt.savefig(f'/tmp/{ticker}_chart.png', dpi=100, bbox_inches='tight')
    plt.close()
    
    print(f"✅ Chart saved as /tmp/{ticker}_chart.png")
