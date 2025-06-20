"""
Advanced Stock Price Analysis for Investify
Uses yfinance for data retrieval and technical indicators for trend analysis
Includes RSI, MACD, and Bollinger Bands for more accurate predictions
"""

import os
import sys
import json
import random
import numpy as np
import pandas as pd
import yfinance as yf
from datetime import datetime, timedelta

class SimpleStockAnalyzer:
    """A stock analyzer with technical indicators for accurate predictions"""
    
    def __init__(self):
        """Initialize the stock analyzer"""
        pass
    
    def calculate_rsi(self, data, window=14):
        """Calculate Relative Strength Index
        
        Args:
            data: Pandas Series of prices
            window: RSI window period
            
        Returns:
            Pandas Series with RSI values
        """
        delta = data.diff()
        gain = delta.where(delta > 0, 0)
        loss = -delta.where(delta < 0, 0)
        
        avg_gain = gain.rolling(window=window).mean()
        avg_loss = loss.rolling(window=window).mean()
        
        rs = avg_gain / avg_loss
        rsi = 100 - (100 / (1 + rs))
        
        return rsi
    
    def calculate_macd(self, data, fast=12, slow=26, signal=9):
        """Calculate Moving Average Convergence Divergence
        
        Args:
            data: Pandas Series of prices
            fast: Fast EMA period
            slow: Slow EMA period
            signal: Signal line period
            
        Returns:
            DataFrame with MACD, Signal and Histogram
        """
        ema_fast = data.ewm(span=fast, adjust=False).mean()
        ema_slow = data.ewm(span=slow, adjust=False).mean()
        
        macd_line = ema_fast - ema_slow
        signal_line = macd_line.ewm(span=signal, adjust=False).mean()
        histogram = macd_line - signal_line
        
        return pd.DataFrame({
            'macd': macd_line,
            'signal': signal_line,
            'histogram': histogram
        })
    
    def calculate_bollinger_bands(self, data, window=20, num_std=2):
        """Calculate Bollinger Bands
        
        Args:
            data: Pandas Series of prices
            window: Moving average window
            num_std: Number of standard deviations
            
        Returns:
            DataFrame with Middle, Upper and Lower bands
        """
        middle_band = data.rolling(window=window).mean()
        std_dev = data.rolling(window=window).std()
        
        upper_band = middle_band + (std_dev * num_std)
        lower_band = middle_band - (std_dev * num_std)
        
        return pd.DataFrame({
            'middle': middle_band,
            'upper': upper_band,
            'lower': lower_band
        })
        
    def analyze(self, ticker, period="1mo"):
        """
        Analyze a stock ticker and provide predictions with technical indicators
        
        Args:
            ticker: Stock ticker symbol
            period: Data period (1d, 5d, 1mo, 3mo, 6mo, 1y, 2y, 5y, 10y, ytd, max)
            
        Returns:
            dict with prediction data and technical indicators
        """
        try:
            # Get stock data with extended history
            stock = yf.Ticker(ticker)
            hist = stock.history(period=period if period else "6mo")
            
            if len(hist) == 0:
                return {
                    "error": f"No data available for ticker {ticker}"
                }
            
            # Calculate basic statistics
            current_price = hist['Close'].iloc[-1]
            avg_price = hist['Close'].mean()
            std_price = hist['Close'].std()
            min_price = hist['Close'].min()
            max_price = hist['Close'].max()
            
            # Calculate technical indicators
            rsi = self.calculate_rsi(hist['Close'])
            macd_data = self.calculate_macd(hist['Close'])
            bb_data = self.calculate_bollinger_bands(hist['Close'])
            
            # Get latest values
            latest_rsi = rsi.iloc[-1] if not pd.isna(rsi.iloc[-1]) else 50
            latest_macd = macd_data['histogram'].iloc[-1] if not pd.isna(macd_data['histogram'].iloc[-1]) else 0
            latest_bb_pos = (current_price - bb_data['lower'].iloc[-1]) / (bb_data['upper'].iloc[-1] - bb_data['lower'].iloc[-1]) if not pd.isna(bb_data['lower'].iloc[-1]) and not pd.isna(bb_data['upper'].iloc[-1]) and (bb_data['upper'].iloc[-1] - bb_data['lower'].iloc[-1]) != 0 else 0.5
            
            # Determine trend direction based on TECHNICAL indicators
            # RSI: Above 70 = overbought (possible DOWN), Below 30 = oversold (possible UP)
            # MACD: Positive histogram = UP, Negative histogram = DOWN
            # BB Position: Close to upper band might be overbought, close to lower band might be oversold
            
            # Weight the signals
            signals = {
                "rsi": 0.3 * (-1 if latest_rsi > 70 else 1 if latest_rsi < 30 else 0),
                "macd": 0.4 * (1 if latest_macd > 0 else -1 if latest_macd < 0 else 0),
                "bb": 0.3 * (-1 if latest_bb_pos > 0.8 else 1 if latest_bb_pos < 0.2 else 0)
            }
            
            # Calculate weighted signal
            weighted_signal = sum(signals.values())
            
            # Determine trend direction
            if weighted_signal > 0.1:
                direction = "UP"
                confidence = min(90, 50 + abs(weighted_signal) * 50)
            elif weighted_signal < -0.1:
                direction = "DOWN"
                confidence = min(90, 50 + abs(weighted_signal) * 50)
            else:
                direction = "NEUTRAL"
                confidence = 50 + abs(weighted_signal) * 50
            
            # Calculate predicted price
            # Use a more sophisticated prediction based on technical indicators
            predicted_change_pct = weighted_signal * 2  # Scale the signal
            predicted_price = current_price * (1 + predicted_change_pct / 100)
            predicted_price = round(predicted_price, 2)
            
            # Key factors based on technical indicators
            factors = []
            
            # RSI factor
            if latest_rsi > 70:
                factors.append(f"RSI indicates overbought at {latest_rsi:.1f}")
            elif latest_rsi < 30:
                factors.append(f"RSI indicates oversold at {latest_rsi:.1f}")
            else:
                factors.append(f"RSI is neutral at {latest_rsi:.1f}")
            
            # MACD factor
            if latest_macd > 0.1:
                factors.append("MACD shows strong bullish momentum")
            elif latest_macd < -0.1:
                factors.append("MACD shows strong bearish momentum")
            else:
                factors.append("MACD indicates sideways momentum")
            
            # Bollinger Bands factor
            if latest_bb_pos > 0.85:
                factors.append("Price near upper Bollinger Band (potential resistance)")
            elif latest_bb_pos < 0.15:
                factors.append("Price near lower Bollinger Band (potential support)")
            else:
                factors.append("Price within normal Bollinger Band range")
                
            # Volume trend
            avg_volume = hist['Volume'].mean()
            recent_volume = hist['Volume'].iloc[-5:].mean()
            if recent_volume > avg_volume * 1.2:
                factors.append("Above average trading volume")
            elif recent_volume < avg_volume * 0.8:
                factors.append("Below average trading volume")
            
            # Create technical indicator data for charting
            technical_data = {
                "rsi": rsi.iloc[-20:].tolist(),
                "macd": macd_data['macd'].iloc[-20:].tolist(),
                "macd_signal": macd_data['signal'].iloc[-20:].tolist(),
                "macd_histogram": macd_data['histogram'].iloc[-20:].tolist(),
                "bollinger_middle": bb_data['middle'].iloc[-20:].tolist(),
                "bollinger_upper": bb_data['upper'].iloc[-20:].tolist(),
                "bollinger_lower": bb_data['lower'].iloc[-20:].tolist(),
                "dates": [str(d.date()) for d in hist.index[-20:].tolist()]
            }
            
            return {
                "predicted_price": predicted_price,
                "confidence": round(confidence, 1),
                "direction": direction,
                "factors": factors[:4],  # Limit to 4 factors
                "current_price": float(current_price),
                "avg_price": float(avg_price),
                "std_price": float(std_price),
                "min_price": float(min_price),
                "max_price": float(max_price),
                "technical": technical_data,
                "latest_rsi": float(latest_rsi),
                "latest_macd": float(latest_macd),
                "latest_bb_position": float(latest_bb_pos)
            }
        
        except Exception as e:
            return {
                "error": f"Analysis failed: {str(e)}"
            }

def main():
    """Main function to handle command-line usage"""
    if len(sys.argv) != 2:
        print(json.dumps({"error": "Please provide a ticker symbol"}))
        sys.exit(1)
        
    ticker = sys.argv[1]
    analyzer = SimpleStockAnalyzer()
    result = analyzer.analyze(ticker)
    print(json.dumps(result))

if __name__ == "__main__":
    main()
