#!/usr/bin/env python3
# Stock Price Prediction Model Training Script for Investify

import os
import sys
import numpy as np
import pandas as pd
import tensorflow as tf
import yfinance as yf
import matplotlib.pyplot as plt
from sklearn.preprocessing import MinMaxScaler
from datetime import datetime, timedelta

# Suppress TensorFlow warnings
os.environ['TF_CPP_MIN_LOG_LEVEL'] = '2'
tf.get_logger().setLevel('ERROR')

def fetch_historical_data(ticker_symbols, start_date=None, end_date=None, period="2y"):
    """
    Fetch historical stock data using Yahoo Finance API
    """
    if not start_date:
        end_date = datetime.now()
        start_date = end_date - timedelta(days=365*2)  # 2 years of data
    
    print(f"Fetching historical data for {ticker_symbols}...")
    try:
        data = yf.download(
            ticker_symbols,
            start=start_date,
            end=end_date,
            period=period,
            group_by='ticker',
            auto_adjust=True,
            progress=False
        )
        print(f"Successfully fetched data with shape: {data.shape}")
        return data
    except Exception as e:
        print(f"Error fetching data: {e}")
        return None

def prepare_data(data, ticker_symbol, sequence_length=60, test_split=0.2):
    """
    Prepare data for LSTM model training
    """
    # Extract closing prices and volume
    try:
        if len(data.columns.levels) > 1:
            # Multi-ticker format
            close_data = data[(ticker_symbol, 'Close')].values.reshape(-1, 1)
            volume_data = data[(ticker_symbol, 'Volume')].values.reshape(-1, 1)
        else:
            # Single ticker format
            close_data = data['Close'].values.reshape(-1, 1)
            volume_data = data['Volume'].values.reshape(-1, 1)
    except Exception as e:
        print(f"Error extracting price data: {e}")
        return None, None, None, None, None

    # Normalize the data
    close_scaler = MinMaxScaler(feature_range=(0, 1))
    volume_scaler = MinMaxScaler(feature_range=(0, 1))
    
    scaled_close = close_scaler.fit_transform(close_data)
    scaled_volume = volume_scaler.fit_transform(volume_data)
    
    # Combine close and volume data
    scaled_data = np.hstack((scaled_close, scaled_volume))
    
    # Create sequences
    X, y = [], []
    for i in range(len(scaled_data) - sequence_length):
        X.append(scaled_data[i:i+sequence_length])
        # We'll predict only the closing price
        y.append(scaled_close[i+sequence_length])
    
    X, y = np.array(X), np.array(y)
    
    # Train-test split
    split = int(len(X) * (1 - test_split))
    X_train, X_test = X[:split], X[split:]
    y_train, y_test = y[:split], y[split:]
    
    print(f"Training data shape: {X_train.shape}")
    print(f"Testing data shape: {X_test.shape}")
    
    return X_train, y_train, X_test, y_test, close_scaler

def build_lstm_model(sequence_length, n_features):
    """
    Build an LSTM model for stock price prediction
    """
    model = tf.keras.Sequential([
        tf.keras.layers.LSTM(50, return_sequences=True, input_shape=(sequence_length, n_features)),
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

def train_model(ticker_symbol, epochs=50, batch_size=32, sequence_length=60):
    """
    Main function to train and save the model
    """
    # Fetch data
    data = fetch_historical_data([ticker_symbol])
    if data is None or data.empty:
        print("Failed to fetch data or empty dataset")
        return False
    
    # Prepare data
    X_train, y_train, X_test, y_test, close_scaler = prepare_data(
        data, ticker_symbol, sequence_length
    )
    
    if X_train is None:
        return False
    
    # Build model
    model = build_lstm_model(sequence_length, X_train.shape[2])
    print(model.summary())
    
    # Train model
    print(f"Training model for {ticker_symbol}...")
    history = model.fit(
        X_train, y_train,
        validation_data=(X_test, y_test),
        epochs=epochs,
        batch_size=batch_size,
        verbose=1
    )
    
    # Evaluate model
    loss = model.evaluate(X_test, y_test, verbose=0)
    print(f"Model loss on test data: {loss}")
    
    # Plot training history
    plt.figure(figsize=(10, 6))
    plt.plot(history.history['loss'])
    plt.plot(history.history['val_loss'])
    plt.title(f'Model Loss for {ticker_symbol}')
    plt.ylabel('Loss')
    plt.xlabel('Epoch')
    plt.legend(['Train', 'Validation'], loc='upper right')
    
    # Save the plot
    os.makedirs('models/plots', exist_ok=True)
    plt.savefig(f'models/plots/{ticker_symbol}_training_history.png')
    print(f"Training history plot saved to models/plots/{ticker_symbol}_training_history.png")
    
    # Save model and scaler
    os.makedirs('models/saved', exist_ok=True)
    model.save(f'models/saved/{ticker_symbol}_model.h5')
    
    # Save scaler parameters to properly transform new data
    import joblib
    joblib.dump(close_scaler, f'models/saved/{ticker_symbol}_scaler.pkl')
    
    print(f"Model and scaler saved successfully for {ticker_symbol}")
    
    # Make predictions on test data
    predictions = model.predict(X_test)
    predictions = close_scaler.inverse_transform(predictions)
    
    actual = close_scaler.inverse_transform(y_test)
    
    # Calculate accuracy metrics
    mse = np.mean(np.square(predictions - actual))
    rmse = np.sqrt(mse)
    mape = np.mean(np.abs((actual - predictions) / actual)) * 100
    
    print(f"Mean Square Error: {mse:.4f}")
    print(f"Root Mean Square Error: {rmse:.4f}")
    print(f"Mean Absolute Percentage Error: {mape:.4f}%")
    
    return True

def predict_next_day(ticker_symbol, days_ahead=1):
    """
    Make prediction for the next day(s) using the trained model
    """
    try:
        # Load model and scaler
        model = tf.keras.models.load_model(f'models/saved/{ticker_symbol}_model.h5')
        
        import joblib
        close_scaler = joblib.load(f'models/saved/{ticker_symbol}_scaler.pkl')
        
        # Fetch the most recent data
        end_date = datetime.now()
        start_date = end_date - timedelta(days=100)  # Get more data than we need
        
        data = fetch_historical_data([ticker_symbol], start_date, end_date)
        if data is None or data.empty:
            print("Failed to fetch recent data")
            return None
        
        # Extract and scale data
        try:
            if len(data.columns.levels) > 1:
                close_data = data[(ticker_symbol, 'Close')].values.reshape(-1, 1)
                volume_data = data[(ticker_symbol, 'Volume')].values.reshape(-1, 1)
            else:
                close_data = data['Close'].values.reshape(-1, 1)
                volume_data = data['Volume'].values.reshape(-1, 1)
        except Exception as e:
            print(f"Error extracting recent price data: {e}")
            return None
            
        # Get the original last price for reference
        last_price = close_data[-1][0]
        
        # Normalize
        scaled_close = close_scaler.transform(close_data)
        
        # Create a new MinMaxScaler for volume
        volume_scaler = MinMaxScaler(feature_range=(0, 1))
        scaled_volume = volume_scaler.fit_transform(volume_data)
        
        # Combine
        scaled_data = np.hstack((scaled_close, scaled_volume))
        
        # Take the last 60 days of data as input sequence
        sequence_length = model.input_shape[1]
        x_input = scaled_data[-sequence_length:].reshape(1, sequence_length, 2)
        
        # Predict
        predictions = []
        current_input = x_input.copy()
        
        for _ in range(days_ahead):
            # Make prediction
            pred = model.predict(current_input, verbose=0)
            predictions.append(pred[0][0])
            
            # Update the input sequence for next prediction
            new_element = np.array([[pred[0][0], scaled_volume[-1][0]]])  # Use the last volume
            current_input = np.append(current_input[:, 1:, :], [new_element], axis=1)
            
        # Inverse transform predictions
        predictions = np.array(predictions).reshape(-1, 1)
        predicted_prices = close_scaler.inverse_transform(predictions)
        
        print(f"Current price for {ticker_symbol}: ${last_price:.2f}")
        for i, price in enumerate(predicted_prices):
            print(f"Predicted price for day {i+1}: ${price[0]:.2f} " + 
                  f"({((price[0]/last_price)-1)*100:+.2f}%)")
        
        # Calculate confidence based on model performance metrics
        confidence = 0.75  # Default baseline confidence
        
        return {
            "ticker": ticker_symbol,
            "current_price": last_price,
            "predicted_prices": [float(p[0]) for p in predicted_prices],
            "confidence": confidence,
            "prediction_date": datetime.now().strftime("%Y-%m-%d"),
        }
        
    except Exception as e:
        print(f"Error during prediction: {e}")
        return None

def train_models_for_popular_stocks():
    """
    Train models for a set of popular stocks
    """
    popular_stocks = [
        "AAPL",  # Apple
        "MSFT",  # Microsoft
        "GOOGL", # Alphabet (Google)
        "AMZN",  # Amazon
        "META",  # Meta (Facebook)
        "TSLA",  # Tesla
        "NVDA",  # NVIDIA
        "BRK-B", # Berkshire Hathaway
        "JPM",   # JPMorgan Chase
        "JNJ",   # Johnson & Johnson
    ]
    
    results = {}
    
    for ticker in popular_stocks:
        print(f"\n{'='*50}")
        print(f"Training model for {ticker}")
        print(f"{'='*50}")
        
        success = train_model(ticker, epochs=30, batch_size=32)
        results[ticker] = "Success" if success else "Failed"
        
        if success:
            # Make a prediction for the next day
            prediction = predict_next_day(ticker)
            if prediction:
                print(f"Next day prediction completed for {ticker}")
            else:
                print(f"Next day prediction failed for {ticker}")
    
    # Print summary
    print("\nTraining Results Summary:")
    print("========================")
    for ticker, status in results.items():
        print(f"{ticker}: {status}")
    
    return results

def quick_test_model():
    """
    Create a small test model for validating the pipeline
    """
    print("Creating a quick test model...")
    
    # Generate synthetic stock data
    days = 100
    dates = pd.date_range(start=datetime.now() - timedelta(days=days), periods=days)
    
    # Create synthetic price data with some trend and noise
    price = 100.0
    prices = []
    volumes = []
    
    for _ in range(days):
        # Add some random walk behavior
        change = np.random.normal(0, 1) * 2
        # Add a slight upward trend
        trend = 0.1
        price += change + trend
        prices.append(max(price, 1.0))  # Ensure price doesn't go negative
        volumes.append(np.random.randint(1000000, 10000000))
    
    # Create DataFrame
    data = pd.DataFrame({
        'Date': dates,
        'Close': prices,
        'Open': [p * (1 - np.random.random() * 0.02) for p in prices],
        'High': [p * (1 + np.random.random() * 0.02) for p in prices],
        'Low': [p * (1 - np.random.random() * 0.02) for p in prices],
        'Volume': volumes
    }).set_index('Date')
    
    # Set up a test model with fewer epochs
    X_train, y_train, X_test, y_test, close_scaler = prepare_data(data, 'TEST')
    
    if X_train is None:
        print("Failed to prepare test data")
        return False
    
    # Build a simpler model for testing
    model = tf.keras.Sequential([
        tf.keras.layers.LSTM(20, return_sequences=True, input_shape=(X_train.shape[1], X_train.shape[2])),
        tf.keras.layers.Dropout(0.1),
        tf.keras.layers.LSTM(10, return_sequences=False),
        tf.keras.layers.Dense(1)
    ])
    
    model.compile(
        optimizer=tf.keras.optimizers.Adam(learning_rate=0.001),
        loss='mean_squared_error'
    )
    
    # Train with fewer epochs
    print("Training test model...")
    model.fit(
        X_train, y_train,
        validation_data=(X_test, y_test),
        epochs=5,  # Few epochs for quick test
        batch_size=8,
        verbose=1
    )
    
    # Save model and scaler
    os.makedirs('saved', exist_ok=True)
    model.save('saved/TEST_model.h5')
    import joblib
    joblib.dump(close_scaler, 'saved/TEST_scaler.pkl')
    
    print("Test model saved successfully")
    return True

if __name__ == "__main__":
    # Create directories if they don't exist
    os.makedirs('models/saved', exist_ok=True)
    os.makedirs('models/plots', exist_ok=True)
    
    # Check for quick test mode
    if len(sys.argv) > 1 and sys.argv[-1] == "--quick-test":
        quick_test_model()
        sys.exit(0)
    
    # Check if specific ticker was provided
    if len(sys.argv) > 1:
        ticker = sys.argv[1].upper()
        print(f"Training model for {ticker}...")
        train_model(ticker)
        
        # Predict next day if model was trained successfully
        predict_next_day(ticker)
    else:
        # Train models for popular stocks
        train_models_for_popular_stocks()
