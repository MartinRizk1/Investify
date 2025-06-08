import pandas as pd
from sklearn.linear_model import LogisticRegression
import numpy as np

class MarketAnalyzer:
    """Analyzes market data and predicts trade actions."""
    def __init__(self):
        self.model = LogisticRegression()

    def train(self, data: pd.DataFrame):
        # Placeholder: Use price change as target
        data = data.dropna()
        X = data[['Open', 'High', 'Low', 'Close', 'Volume']].values[:-1]
        y = (data['Close'].shift(-1) > data['Close']).astype(int)[:-1]
        self.model.fit(X, y)

    def predict(self, data: pd.DataFrame) -> str:
        X = data[['Open', 'High', 'Low', 'Close', 'Volume']].values[-1:]
        prob = self.model.predict_proba(X)[0][1]
        if prob > 0.6:
            return 'BUY'
        elif prob < 0.4:
            return 'SELL'
        else:
            return 'HOLD'
