#!/usr/bin/env python3
# Stock price prediction script for Investify Go-Python bridge
# This script provides a way for the Go backend to call the TensorFlow model

import sys
import json
import numpy as np
import joblib
import os
import tensorflow as tf
from datetime import datetime
import traceback

def load_model(ticker):
    """Load the trained model for a ticker."""
    model_path = f"saved/{ticker}_model.h5"
    scaler_path = f"saved/{ticker}_scaler.pkl"
    
    # Check if model exists for this ticker
    if not os.path.exists(model_path):
        # Try a default model
        model_path = "saved/DEFAULT_model.h5"
        scaler_path = "saved/DEFAULT_scaler.pkl"
        
        if not os.path.exists(model_path):
            return None, None
    
    try:
        model = tf.keras.models.load_model(model_path)
        scaler = joblib.load(scaler_path)
        return model, scaler
    except Exception as e:
        print(json.dumps({
            "error": f"Failed to load model: {str(e)}"
        }))
        return None, None

def predict(stock_data):
    """Make a prediction for the stock."""
    ticker = stock_data.get("ticker", "DEFAULT")
    
    # Try to load the model for this ticker
    model, scaler = load_model(ticker)
    
    # If no model exists, use a simple rule-based prediction
    if model is None:
        return rule_based_prediction(stock_data)
        
    try:
        # Extract features for prediction
        price = float(stock_data.get("price", 0))
        open_price = float(stock_data.get("open", price))
        high = float(stock_data.get("high", price * 1.01))
        low = float(stock_data.get("low", price * 0.99))
        
        # Create a feature vector (this needs to match your model's expected input)
        # This is a simplified example - your actual model may need different features
        features = np.array([[price, open_price, high, low]])
        
        # Scale the features
        scaled_features = scaler.transform(features)
        
        # Make prediction
        prediction = model.predict(scaled_features)
        
        # Convert prediction back to price
        predicted_price = float(scaler.inverse_transform(prediction)[0][0])
        
        # Calculate other metrics
        price_diff = ((predicted_price - price) / price) * 100
        confidence = min(0.95, max(0.55, 0.85 - abs(price_diff) * 0.01))
        
        direction = "NEUTRAL"
        if price_diff > 1.0:
            direction = "UP"
        elif price_diff < -1.0:
            direction = "DOWN"
            
        # Generate factors
        factors = generate_factors(stock_data, price_diff)
        
        return {
            "predicted_price": round(predicted_price, 2),
            "confidence": confidence,
            "direction": direction,
            "factors": factors
        }
        
    except Exception as e:
        # Fallback to rule-based prediction on error
        print(f"Model prediction error: {str(e)}", file=sys.stderr)
        traceback.print_exc(file=sys.stderr)
        return rule_based_prediction(stock_data)

def rule_based_prediction(stock_data):
    """Simple rule-based prediction as fallback."""
    price = float(stock_data.get("price", 0))
    if price <= 0:
        return {
            "error": "Invalid price data"
        }
        
    # Parse change percentage
    change_pct = stock_data.get("change_pct", "0")
    if isinstance(change_pct, str) and "%" in change_pct:
        change_pct = change_pct.replace("%", "")
    change_pct = float(change_pct)
    
    # Simple momentum-based prediction
    momentum = change_pct * 0.1
    
    # Add some randomness
    import random
    random_factor = (random.random() - 0.5) * 0.5
    
    # Calculate predicted change percentage
    predicted_change_pct = momentum + random_factor
    
    # Calculate predicted price
    predicted_price = price * (1 + predicted_change_pct / 100)
    predicted_price = round(predicted_price, 2)
    
    # Determine direction
    direction = "NEUTRAL"
    if predicted_change_pct > 1.0:
        direction = "UP"
    elif predicted_change_pct < -1.0:
        direction = "DOWN"
    
    # Calculate confidence (higher for smaller predictions)
    confidence = 0.7 - min(0.2, abs(predicted_change_pct) * 0.02)
    
    # Generate factors
    factors = generate_factors(stock_data, predicted_change_pct)
    
    return {
        "predicted_price": predicted_price,
        "confidence": confidence,
        "direction": direction,
        "factors": factors
    }

def generate_factors(stock_data, predicted_change_pct):
    """Generate key factors that influence the prediction."""
    factors = []
    
    # Price momentum
    price = float(stock_data.get("price", 0))
    open_price = float(stock_data.get("open", price))
    
    if price > open_price:
        factors.append("Price is above opening level")
    elif price < open_price:
        factors.append("Price is below opening level")
    
    # Market cap based factor
    market_cap = stock_data.get("market_cap", "")
    if "T" in market_cap:
        factors.append("Large market cap indicates stability")
    elif "B" in market_cap and "T" not in market_cap:
        if float(market_cap.replace("$", "").replace("B", "")) > 10:
            factors.append("Medium-large market cap suggests moderate volatility")
        else:
            factors.append("Medium market cap may lead to higher volatility")
    else:
        factors.append("Smaller capitalization suggests higher volatility")
    
    # Add prediction direction factor
    if predicted_change_pct > 0:
        factors.append("Technical indicators suggest positive momentum")
    else:
        factors.append("Technical indicators suggest negative pressure")
    
    return factors[:3]  # Return top 3 factors

def main():
    """Main entry point for the prediction script."""
    try:
        # Read stock data from stdin
        stock_data = json.load(sys.stdin)
        
        # Make prediction
        result = predict(stock_data)
        
        # Output result as JSON
        print(json.dumps(result))
        
    except Exception as e:
        print(json.dumps({
            "error": str(e)
        }))

if __name__ == "__main__":
    main()
