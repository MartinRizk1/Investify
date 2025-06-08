# Test configuration for security tests
import os

class TestConfig:
    """Configuration for testing"""
    TESTING = True
    WTF_CSRF_ENABLED = False  # Disable CSRF for testing
    SECRET_KEY = 'test-secret-key'
    RATE_LIMIT_DEFAULT = '1000 per hour'  # Higher limit for tests
    
    # Override environment variables for testing
    @staticmethod
    def setup_test_env():
        os.environ['FLASK_DEBUG'] = 'False'
        os.environ['SECRET_KEY'] = 'test-secret-key'
        os.environ['RATE_LIMIT_DEFAULT'] = '1000 per hour'
        os.environ['CSRF_TIME_LIMIT'] = '3600'
