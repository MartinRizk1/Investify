import unittest
import pandas as pd
from src.analyzers.market_analyzer import MarketAnalyzer

class TestMarketAnalyzer(unittest.TestCase):
    def test_train_and_predict(self):
        # Create dummy data
        data = pd.DataFrame({
            'Open': [1, 2, 3, 4, 5],
            'High': [2, 3, 4, 5, 6],
            'Low': [0, 1, 2, 3, 4],
            'Close': [1.5, 2.5, 3.5, 4.5, 5.5],
            'Volume': [100, 110, 120, 130, 140]
        })
        analyzer = MarketAnalyzer()
        analyzer.train(data)
        action = analyzer.predict(data)
        self.assertIn(action, ['BUY', 'SELL', 'HOLD'])

if __name__ == '__main__':
    unittest.main()
