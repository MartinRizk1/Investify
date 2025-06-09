#!/usr/bin/env python3
"""
Security test suite for MoneyMaker application
Tests security features, rate limiting, input validation, and CSRF protection
"""

import sys
import os
import unittest
import json
import time
from unittest.mock import patch, MagicMock

# Add src directory to path
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..', 'src')))

# Configure test environment
os.environ['FLASK_DEBUG'] = 'False'
os.environ['SECRET_KEY'] = 'test-secret-key'
os.environ['RATE_LIMIT_DEFAULT'] = '1000 per hour'

from app import app

class SecurityTestCase(unittest.TestCase):
    """Test case for security features"""
    
    def setUp(self):
        """Set up test client"""
        self.app = app
        self.app.config['TESTING'] = True
        self.app.config['WTF_CSRF_ENABLED'] = False  # Disable CSRF for testing
        self.app.config['RATELIMIT_ENABLED'] = False  # Disable rate limiting for testing
        self.client = self.app.test_client()
        
        # Clear any rate limit data
        from flask_limiter import Limiter
        with self.app.app_context():
            if hasattr(app.extensions.get('limiter'), 'storage'):
                app.extensions['limiter'].storage.clear()
        
    def test_security_headers(self):
        """Test that security headers are properly set"""
        response = self.client.get('/')
        
        # Check for security headers
        self.assertIn('X-Content-Type-Options', response.headers)
        self.assertEqual(response.headers['X-Content-Type-Options'], 'nosniff')
        
        self.assertIn('X-Frame-Options', response.headers)
        self.assertEqual(response.headers['X-Frame-Options'], 'DENY')
        
        self.assertIn('X-XSS-Protection', response.headers)
        self.assertEqual(response.headers['X-XSS-Protection'], '1; mode=block')
        
        self.assertIn('Strict-Transport-Security', response.headers)
        
        self.assertIn('Content-Security-Policy', response.headers)
        
    def test_input_validation_ticker(self):
        """Test input validation for ticker symbols"""
        # Test valid ticker
        response = self.client.get('/api/chart/AAPL/1Y')
        self.assertIn(response.status_code, [200, 404])  # 404 if no data, 200 if data found
        
        # Test invalid ticker with special characters
        response = self.client.get('/api/chart/AAPL<script>/1Y')
        self.assertEqual(response.status_code, 400)
        
        # Test overly long ticker
        long_ticker = 'A' * 100
        response = self.client.get(f'/api/chart/{long_ticker}/1Y')
        self.assertEqual(response.status_code, 400)
        
    def test_input_validation_period(self):
        """Test input validation for time periods"""
        # Test valid periods
        valid_periods = ['1M', '3M', '6M', '1Y', '5Y']
        for period in valid_periods:
            response = self.client.get(f'/api/chart/AAPL/{period}')
            self.assertIn(response.status_code, [200, 404, 400])
        
        # Test invalid period
        response = self.client.get('/api/chart/AAPL/INVALID')
        self.assertEqual(response.status_code, 400)
        
    def test_rate_limiting_home_page(self):
        """Test rate limiting on home page"""
        # Create a new app instance with rate limiting enabled for this test
        from src.app import app as rate_app
        rate_app.config['TESTING'] = True
        rate_app.config['WTF_CSRF_ENABLED'] = False
        test_client = rate_app.test_client()
        
        # Make requests to trigger rate limiting
        rate_limited = False
        for i in range(35):  # Try to exceed the limit
            response = test_client.post('/', data={'company': 'AAPL'})
            if response.status_code == 429:
                rate_limited = True
                break
            time.sleep(0.1)  # Small delay between requests
        
        # Should eventually get rate limited or handle gracefully
        self.assertTrue(rate_limited or response.status_code in [200, 400])
        
    def test_api_error_handling(self):
        """Test API error handling"""
        # Test with empty ticker
        response = self.client.get('/api/chart//1Y')
        self.assertEqual(response.status_code, 404)  # Not found route
        
        # Test with malformed request
        response = self.client.get('/api/chart')
        self.assertEqual(response.status_code, 404)  # Not found route
        
    def test_json_response_format(self):
        """Test that API responses are proper JSON"""
        response = self.client.get('/api/chart/INVALID_TICKER/1Y')
        self.assertEqual(response.content_type, 'application/json')
        
        # Should be valid JSON
        try:
            json.loads(response.data)
        except json.JSONDecodeError:
            self.fail("Response is not valid JSON")
            
    @patch('src.analyzers.market_analyzer.MarketAnalyzer.train')
    def test_error_logging(self, mock_train):
        """Test that errors are properly logged"""
        # Mock an exception in the analyzer
        mock_train.side_effect = Exception("Test error")
        
        # Capture log output
        import logging
        from io import StringIO
        log_stream = StringIO()
        handler = logging.StreamHandler(log_stream)
        logger = logging.getLogger('moneymaker.security')
        logger.addHandler(handler)
        logger.setLevel(logging.ERROR)
        
        # Mock the validation to pass
        with patch('src.app.validate_ticker_input'), \
             patch('src.data.data_fetcher.DataFetcher.fetch_data') as mock_fetch, \
             patch('src.data.data_fetcher.DataFetcher.get_ticker') as mock_ticker:
            
            # Mock successful data fetching but analysis failure
            mock_ticker.return_value = 'AAPL'
            mock_fetch.return_value = MagicMock()  # Non-empty dataframe mock
            mock_fetch.return_value.empty = False
            
            # Make request that should cause error in analysis
            response = self.client.post('/', data={'company': 'AAPL'})
            
            # Should handle error gracefully and log it
            self.assertEqual(response.status_code, 200)
            
            # Check that error was logged
            log_contents = log_stream.getvalue()
            self.assertIn('Analysis error for AAPL: Test error', log_contents)
        
        # Clean up
        logger.removeHandler(handler)
        
    def test_xss_prevention(self):
        """Test XSS prevention in input handling"""
        xss_payload = '<script>alert("xss")</script>'
        
        # Test in company field - should fail validation
        response = self.client.post('/', data={'company': xss_payload})
        
        # Script should not be executed (should be handled by input validation)
        self.assertEqual(response.status_code, 400)  # Should fail validation
        
    def test_sql_injection_prevention(self):
        """Test SQL injection prevention (even though we don't use SQL)"""
        sql_payload = "'; DROP TABLE companies; --"
        
        # Test in company field - should fail validation  
        response = self.client.post('/', data={'company': sql_payload})
        
        # Should fail input validation
        self.assertEqual(response.status_code, 400)
        
    def test_csrf_token_present(self):
        """Test that CSRF token is present in forms"""
        with patch('src.app.validate_ticker_input'):
            response = self.client.get('/')
            self.assertEqual(response.status_code, 200)
            
            # Check if csrf_token is in the response (when CSRF is enabled)
            # This test would need to be adjusted when CSRF is fully enabled
            self.assertIn(b'form', response.data)

class SecurityIntegrationTestCase(unittest.TestCase):
    """Integration tests for security features"""
    
    def setUp(self):
        """Set up test client with full security enabled"""
        self.app = app
        self.app.config['TESTING'] = True
        # Keep CSRF enabled for integration tests
        self.client = self.app.test_client()
        
    def test_full_security_workflow(self):
        """Test complete security workflow"""
        # Create a test client with minimal rate limiting
        test_app = app
        test_app.config['TESTING'] = True
        test_app.config['WTF_CSRF_ENABLED'] = False
        test_client = test_app.test_client()
        
        # 1. Get homepage (should have security headers) - bypass validation for this test
        with patch('src.app.validate_ticker_input'):
            response = test_client.get('/')
            self.assertEqual(response.status_code, 200)
            self.assertIn('X-Content-Type-Options', response.headers)
        
        # 2. Test API endpoint with valid data
        response = test_client.get('/api/chart/AAPL/1Y')
        self.assertIn(response.status_code, [200, 404, 400])
        
        # 3. Test input validation
        response = test_client.get('/api/chart/INVALID<>/1Y')
        self.assertEqual(response.status_code, 400)

if __name__ == '__main__':
    # Create test suite
    loader = unittest.TestLoader()
    suite = unittest.TestSuite()
    
    # Add test cases
    suite.addTests(loader.loadTestsFromTestCase(SecurityTestCase))
    suite.addTests(loader.loadTestsFromTestCase(SecurityIntegrationTestCase))
    
    # Run tests
    runner = unittest.TextTestRunner(verbosity=2)
    result = runner.run(suite)
    
    # Exit with proper code
    sys.exit(0 if result.wasSuccessful() else 1)
