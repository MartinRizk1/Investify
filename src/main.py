import sys
import os
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

from data.data_fetcher import DataFetcher
from analyzers.market_analyzer import MarketAnalyzer
from utils.helpers import print_trade_action, display_stock_graph, display_company_info, display_earnings_info
import matplotlib.pyplot as plt
import pandas as pd


def main():
    print("=" * 60)
    print("🎯 MONEY MAKER - AI Stock Analysis Tool")
    print("=" * 60)
    
    # Check for command line arguments
    if len(sys.argv) > 1:
        company = ' '.join(sys.argv[1:]).strip()
    else:
        company = input("\nEnter company name or ticker (e.g., Apple, AAPL, Tesla, TSLA): ").strip()
    
    if not company:
        print("❌ Please enter a valid company name or ticker symbol.")
        return
    
    print(f"\n🔍 Analyzing {company}...")
    
    # Initialize data fetcher
    fetcher = DataFetcher()
    
    # Get ticker symbol
    ticker = fetcher.get_ticker(company)
    print(f"📊 Resolved ticker: {ticker}")
    
    # Fetch data
    print("📈 Fetching stock data...")
    data = fetcher.fetch_data(company)
    
    print("🏢 Fetching company information...")
    company_info = fetcher.fetch_company_info(company)
    
    print("📰 Fetching latest news...")
    news = fetcher.fetch_news(company)
    
    if data.empty:
        print(f"❌ No stock data found for '{company}'. Please check the company name or ticker symbol.")
        return
    
    print("\n" + "=" * 60)
    print("📊 ANALYSIS RESULTS")
    print("=" * 60)
    
    # Display company information
    display_company_info(company_info, ticker)
    
    # Display earnings and financial information
    display_earnings_info(company_info, data)
    
    # Display stock graph
    display_stock_graph(data, ticker, company_info.get('shortName', ticker))
    
    # Perform AI analysis
    print("\n🤖 AI ANALYSIS")
    print("-" * 30)
    analyzer = MarketAnalyzer()
    analyzer.train(data)
    action = analyzer.predict(data)
    
    # Enhanced trade recommendation
    print_trade_action(ticker, action, data, company_info)
    
    # Display recent news
    if news:
        print("\n📰 RECENT NEWS")
        print("-" * 30)
        for i, article in enumerate(news[:5], 1):
            print(f"{i}. {article.get('title', 'No title')}")
            print(f"   📅 {article.get('providerPublishTime', 'Unknown date')}")
            print(f"   🔗 {article.get('link', 'No link')}\n")
    
    print("=" * 60)
    print("✅ Analysis complete!")


if __name__ == "__main__":
    main()
