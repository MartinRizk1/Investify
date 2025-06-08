import unittest
from src.data.data_fetcher import DataFetcher

class TestDataFetcher(unittest.TestCase):
    def test_get_ticker(self):
        fetcher = DataFetcher()
        self.assertEqual(fetcher.get_ticker('AAPL'), 'AAPL')

    def test_fetch_data(self):
        fetcher = DataFetcher()
        data = fetcher.fetch_data('AAPL', period='5d')
        self.assertFalse(data.empty)

if __name__ == '__main__':
    unittest.main()
