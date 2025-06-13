# Investify Stock Prediction Model

This directory contains the machine learning models and training scripts for the Investify application's stock prediction feature.

## Setup

1. Ensure you have Python 3.7+ installed
2. Install dependencies:
   ```
   pip install -r requirements.txt
   ```

## Training Models

To train a model for a specific stock ticker:

```bash
python train_stock_model.py AAPL
```

To train models for a set of popular stocks:

```bash
python train_stock_model.py
```

## Model Architecture

The model uses a stacked LSTM (Long Short-Term Memory) neural network architecture:
- Two LSTM layers with 50 units each
- Dropout layers (0.2) for regularization
- Dense output layers

The model is trained on sequences of 60 days of price and volume data to predict the next day's closing price.

## Model Files

After training, the following files will be generated in the `models/saved` directory:
- `{TICKER}_model.h5`: The trained TensorFlow model
- `{TICKER}_scaler.pkl`: The scaler used to normalize the data (needed for predictions)

Training history plots are saved in the `models/plots` directory.

## Integration with Go Application

The Go application uses the trained models through a TensorFlow Serving API or by directly loading the models in a Python subprocess.

The `tensorflow_model.go` file in the Investify application handles the communication between the Go application and the TensorFlow models.
