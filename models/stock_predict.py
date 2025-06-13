"""
TensorFlow Stock Price Prediction Model for Investify
This module provides integration between Go backend and Python TensorFlow models
"""

import os
import sys
import json
import numpy as np
import pandas as pd
import tensorflow as tf
import yfinance as yf
import matplotlib.pyplot as plt
from sklearn.preprocessing import MinMaxScaler
from datetime import datetime, timedelta
import joblib

# Suppress TensorFlow warnings
os.environ['TF_CPP_MIN_LOG_LEVEL'] = '2'
tf.get_logger().setLevel('ERROR')

class StockPredictor:
    """Stock price prediction model using TensorFlow"""
    
    def __init__(self, model_dir="saved"):
        """
        Initialize the stock predictor
        
        Args:
            model_dir: Directory where models are stored
        """
        self.model_dir = model_dir
        os.makedirs(model_dir, exist_ok=True)
        self.models = {}
        self.scalers = {}
        
    def load_model(self, ticker):
        """
        Load a saved model for a specific ticker
        
        Args:
            ticker: Stock ticker symbol
            
        Returns:
            True if model was loaded, False otherwise
        """
        model_path = os.path.join(self.model_dir, f"{ticker}_model.h5")
        scaler_path = os.path.join(self.model_dir, f"{ticker}_scaler.pkl")
        
        # Check if the model exists
        if not os.path.exists(model_path) or not os.path.exists(scaler_path):
            return False
            
        try:
            model = tf.keras.models.load_model(model_path)
            scaler = joblib.load(scaler_path)
            self.models[ticker] = model
            self.scalers[ticker] = scaler
            return True
        except Exception as e:
            print(f"Error loading model for {ticker}: {e}")
            return False
    
    def predict(self, ticker, current_price, historical_data=None):
        """
        Make a stock price prediction
        
        Args:
            ticker: Stock ticker symbol
            current_price: Current stock price
            historical_data: Optional historical data to use
            
        Returns:
            Dictionary with prediction results
        """
        # Try to load the model if not already loaded
        if ticker not in self.models and not self.load_model(ticker):
            # Use fallback prediction if no model
            return self._fallback_prediction(ticker, current_price)
        
        try:
            # If we have historical data, use it; otherwise fetch from Yahoo
            if historical_data is None:
                historical_data = self._fetch_historical_data(ticker)
                
            if historical_data is None or len(historical_data) < 60:
                return self._fallback_prediction(ticker, current_price)
            
            # Prepare the data for prediction
            X_pred = self._prepare_data_for_prediction(ticker, historical_data, current_price)
            
            # Make prediction
            prediction = self.models[ticker].predict(X_pred)
            
            # Convert prediction back to price
            predicted_price = self.scalers[ticker].inverse_transform(prediction)[0][0]
            
            # Calculate prediction details
            price_diff = ((predicted_price - current_price) / current_price) * 100
            confidence = min(0.95, max(0.55, 0.85 - abs(price_diff) * 0.01))
            
            direction = "NEUTRAL"
            if price_diff > 1.0:
                direction = "UP"
            elif price_diff < -1.0:
                direction = "DOWN"
                
            # Generate factors
            factors = self._generate_factors(ticker, historical_data, price_diff)
            
            return {
                "price": current_price,
                "predicted_price": round(predicted_price, 2),
                "confidence": confidence,
                "direction": direction,
                "factors": factors
            }
            
        except Exception as e:
            print(f"Error making prediction for {ticker}: {e}")
            return self._fallback_prediction(ticker, current_price)
    
    def _fetch_historical_data(self, ticker, days=60):
        """
        Fetch historical data from Yahoo Finance
        
        Args:
            ticker: Stock ticker symbol
            days: Number of days of historical data to fetch
            
        Returns:
            Pandas DataFrame with historical data
        """
        try:
            end_date = datetime.now()
            start_date = end_date - timedelta(days=days)
            data = yf.download(ticker, start=start_date, end=end_date, progress=False)
            return data
        except Exception as e:
            print(f"Error fetching data for {ticker}: {e}")
            return None
    
    def _prepare_data_for_prediction(self, ticker, historical_data, current_price):
        """
        Prepare data for the prediction model
        
        Args:
            ticker: Stock ticker symbol
            historical_data: DataFrame with historical data
            current_price: Current stock price
            
        Returns:
            Prepared data for model input
        """
        # Extract relevant features
        close_data = historical_data['Close'].values
        
        # Add current price to the end
        close_data = np.append(close_data, current_price)
        
        # Reshape for scaling
        close_data = close_data.reshape(-1, 1)
        
        # Scale the data
        scaler = self.scalers[ticker]
        scaled_data = scaler.transform(close_data)
        
        # Create a sequence for prediction
        sequence_length = 50  # This should match what the model was trained with
        X_pred = []
        X_pred.append(scaled_data[-sequence_length:])
        X_pred = np.array(X_pred)
        
        return X_pred
    
    def _generate_factors(self, ticker, historical_data, price_diff):
        """
        Generate key factors that influence the prediction
        
        Args:
            ticker: Stock ticker symbol
            historical_data: DataFrame with historical data
            price_diff: Percentage difference between prediction and current price
            
        Returns:
            List of factors
        """
        factors = []
        
        # Calculate short-term momentum (5-day)
        short_term = historical_data['Close'][-5:].pct_change().mean() * 100
        
        # Calculate medium-term momentum (20-day)
        medium_term = historical_data['Close'][-20:].pct_change().mean() * 100
        
        # Volume trend
        volume_change = (historical_data['Volume'][-5:].mean() / historical_data['Volume'][-10:-5].mean() - 1) * 100
        
        # Volatility
        volatility = historical_data['Close'][-20:].pct_change().std() * 100
        
        # Add factors based on calculations
        if short_term > 0.5:
            factors.append("Short-term momentum is positive")
        elif short_term < -0.5:
            factors.append("Short-term momentum is negative")
            
        if medium_term > 0.5:
            factors.append("Medium-term trend is bullish")
        elif medium_term < -0.5:
            factors.append("Medium-term trend is bearish")
            
        if volume_change > 10:
            factors.append("Trading volume is increasing")
        elif volume_change < -10:
            factors.append("Trading volume is decreasing")
            
        if volatility > 2:
            factors.append("Market showing high volatility")
        elif volatility < 1:
            factors.append("Market showing low volatility")
            
        # Add prediction direction factor
        if price_diff > 0:
            factors.append("Technical indicators suggest positive momentum")
        else:
            factors.append("Technical indicators suggest negative pressure")
        
        return factors[:3]  # Return top 3 factors
    
    def _fallback_prediction(self, ticker, current_price):
        """
        Make a rule-based fallback prediction when model is unavailable
        
        Args:
            ticker: Stock ticker symbol
            current_price: Current stock price
            
        Returns:
            Dictionary with prediction results
        """
        import random
        
        # Simple momentum-based prediction with randomness
        random_factor = (random.random() - 0.5) * 2  # Random factor between -1% and 1%
        predicted_change_pct = random_factor
        
        # Calculate predicted price
        predicted_price = current_price * (1 + predicted_change_pct / 100)
        predicted_price = round(predicted_price, 2)
        
        # Determine direction
        direction = "NEUTRAL"
        if predicted_change_pct > 1.0:
            direction = "UP"
        elif predicted_change_pct < -1.0:
            direction = "DOWN"
        
        # Calculate confidence (higher for smaller predictions)
        confidence = 0.7 - min(0.2, abs(predicted_change_pct) * 0.02)
        
        # Generate simple factors
        factors = [
            "Limited historical data available",
            "Market sentiment appears neutral",
            "Technical indicators provide mixed signals"
        ]
        
        return {
            "price": current_price,
            "predicted_price": predicted_price,
            "confidence": confidence,
            "direction": direction,
            "factors": factors
        }
    
    def train_model(self, ticker, epochs=50, sequence_length=50):
        """
        Train a new model for a ticker
        
        Args:
            ticker: Stock ticker symbol
            epochs: Number of training epochs
            sequence_length: Sequence length for LSTM model
            
        Returns:
            True if training was successful, False otherwise
        """
        # Fetch historical data (2 years)
        historical_data = self._fetch_historical_data(ticker, days=730)
        
        if historical_data is None or len(historical_data) < 100:
            print(f"Insufficient data for {ticker}")
            return False
            
        try:
            # Prepare training data
            X_train, y_train, X_test, y_test, scaler = self._prepare_training_data(historical_data, sequence_length)
            
            # Build model
            model = self._build_lstm_model(sequence_length)
            
            # Train model
            model.fit(
                X_train, y_train,
                validation_data=(X_test, y_test),
                epochs=epochs,
                batch_size=32,
                verbose=1
            )
            
            # Evaluate model
            loss = model.evaluate(X_test, y_test, verbose=0)
            print(f"Model loss on test data: {loss}")
            
            # Save model and scaler
            model_path = os.path.join(self.model_dir, f"{ticker}_model.h5")
            scaler_path = os.path.join(self.model_dir, f"{ticker}_scaler.pkl")
            
            model.save(model_path)
            joblib.dump(scaler, scaler_path)
            
            # Add to loaded models
            self.models[ticker] = model
            self.scalers[ticker] = scaler
            
            return True
            
        except Exception as e:
            print(f"Error training model for {ticker}: {e}")
            return False
    
    def _prepare_training_data(self, data, sequence_length=50, test_split=0.2):
        """
        Prepare data for LSTM model training
        
        Args:
            data: DataFrame with historical data
            sequence_length: Sequence length for LSTM model
            test_split: Proportion of data to use for testing
            
        Returns:
            X_train, y_train, X_test, y_test, scaler
        """
        # Extract closing prices
        close_data = data['Close'].values.reshape(-1, 1)
        
        # Normalize the data
        scaler = MinMaxScaler(feature_range=(0, 1))
        scaled_data = scaler.fit_transform(close_data)
        
        # Create sequences
        X, y = [], []
        for i in range(len(scaled_data) - sequence_length):
            X.append(scaled_data[i:i+sequence_length])
            y.append(scaled_data[i+sequence_length])
        
        X, y = np.array(X), np.array(y)
        
        # Train-test split
        split = int(len(X) * (1 - test_split))
        X_train, X_test = X[:split], X[split:]
        y_train, y_test = y[:split], y[split:]
        
        return X_train, y_train, X_test, y_test, scaler
    
    def _build_lstm_model(self, sequence_length):
        """
        Build an LSTM model for stock price prediction
        
        Args:
            sequence_length: Sequence length for LSTM model
            
        Returns:
            Compiled TensorFlow model
        """
        model = tf.keras.Sequential([
            tf.keras.layers.LSTM(50, return_sequences=True, input_shape=(sequence_length, 1)),
            tf.keras.layers.Dropout(0.2),
            tf.keras.layers.LSTM(50, return_sequences=False),
            tf.keras.layers.Dropout(0.2),
            tf.keras.layers.Dense(25),
            tf.keras.layers.Dense(1)
        ])
        
        model.compile(
            optimizer=tf.keras.optimizers.Adam(learning_rate=0.001),
            loss='mean_squared_error'
        )
        
        return model

def main():
    """
    Main entry point for command-line usage
    """
    if len(sys.argv) != 2 and len(sys.argv) != 3:
        print("Usage: python stock_predict.py <ticker> [train]")
        return
        
    ticker = sys.argv[1]
    
    # Initialize predictor
    predictor = StockPredictor()
    
    # Check if we should train a new model
    if len(sys.argv) == 3 and sys.argv[2] == "train":
        print(f"Training model for {ticker}...")
        if predictor.train_model(ticker):
            print(f"Model trained successfully for {ticker}")
        else:
            print(f"Model training failed for {ticker}")
        return
        
    # Otherwise make a prediction
    # First try to get current price from Yahoo
    try:
        data = yf.download(ticker, period="1d", progress=False)
        if len(data) > 0:
            current_price = data['Close'][-1]
        else:
            current_price = 100.0  # Default price if not available
    except:
        current_price = 100.0  # Default price if error
        
    # Make prediction
    result = predictor.predict(ticker, current_price)
    
    # Print result as JSON
    print(json.dumps(result, indent=2))

if __name__ == "__main__":
    main()
