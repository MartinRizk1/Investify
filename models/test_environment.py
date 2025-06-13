#!/usr/bin/env python3
# Test script to verify the Python environment and TensorFlow model training

import os
import numpy as np
import tensorflow as tf
import joblib
from sklearn.preprocessing import MinMaxScaler
import sys

def test_environment():
    """Test that the Python environment is set up correctly"""
    print("Testing Python environment for Investify...")
    
    # Check TensorFlow
    print(f"TensorFlow version: {tf.__version__}")
    
    # Test GPU availability (optional)
    gpu_available = len(tf.config.list_physical_devices('GPU')) > 0
    print(f"GPU available: {gpu_available}")
    
    # Create a simple model
    print("Creating and training a simple model...")
    
    # Generate some test data
    x_train = np.array(np.random.random((100, 1)))
    y_train = 3 * x_train + 2 + np.random.random((100, 1)) * 0.1
    
    # Build a simple model
    model = tf.keras.Sequential([
        tf.keras.layers.Dense(10, activation='relu', input_shape=(1,)),
        tf.keras.layers.Dense(1)
    ])
    
    # Compile and train
    model.compile(optimizer='adam', loss='mse')
    model.fit(x_train, y_train, epochs=5, verbose=1)
    
    # Test prediction
    test_input = np.array([[0.5]])
    prediction = model.predict(test_input)
    expected = 3 * 0.5 + 2  # Based on our data generation
    
    print(f"Test prediction: {prediction[0][0]:.4f}, Expected: {expected:.4f}")
    
    # Save and load the model to test the full pipeline
    os.makedirs('saved', exist_ok=True)
    
    # Save model
    model.save('saved/TEST_model.h5')
    
    # Create and save a scaler
    scaler = MinMaxScaler()
    scaler.fit(np.array([[0], [1]]))
    joblib.dump(scaler, 'saved/TEST_scaler.pkl')
    
    # Load model
    loaded_model = tf.keras.models.load_model('saved/TEST_model.h5')
    loaded_scaler = joblib.load('saved/TEST_scaler.pkl')
    
    # Test the loaded model
    loaded_prediction = loaded_model.predict(test_input)
    
    print(f"Loaded model prediction: {loaded_prediction[0][0]:.4f}")
    
    print("\nEnvironment test completed successfully!")
    return True

if __name__ == "__main__":
    try:
        success = test_environment()
        sys.exit(0 if success else 1)
    except Exception as e:
        print(f"Error during testing: {e}")
        sys.exit(1)
