import logging
import os
import re
from logging.handlers import RotatingFileHandler
from datetime import datetime

class SensitiveDataFilter(logging.Filter):
    """Filter to remove sensitive data from log messages"""
    
    def __init__(self):
        super().__init__()
        # Patterns for sensitive data that should be redacted
        self.patterns = [
            (re.compile(r'password[\'"\s]*[:=][\'"\s]*[^\s\'"]+', re.IGNORECASE), 'password=***'),
            (re.compile(r'secret[\'"\s]*[:=][\'"\s]*[^\s\'"]+', re.IGNORECASE), 'secret=***'),
            (re.compile(r'token[\'"\s]*[:=][\'"\s]*[^\s\'"]+', re.IGNORECASE), 'token=***'),
            (re.compile(r'api[_-]?key[\'"\s]*[:=][\'"\s]*[^\s\'"]+', re.IGNORECASE), 'api_key=***'),
            (re.compile(r'bearer [a-zA-Z0-9._-]+', re.IGNORECASE), 'bearer ***'),
            (re.compile(r'[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}'), '***@***.***'),  # Email addresses
            (re.compile(r'\b\d{4}[-\s]?\d{4}[-\s]?\d{4}[-\s]?\d{4}\b'), '****-****-****-****'),  # Credit card numbers
            (re.compile(r'\b\d{3}-\d{2}-\d{4}\b'), '***-**-****'),  # SSN
        ]
    
    def filter(self, record):
        """Filter the log record to remove sensitive data"""
        if hasattr(record, 'msg') and record.msg:
            message = str(record.msg)
            for pattern, replacement in self.patterns:
                message = pattern.sub(replacement, message)
            record.msg = message
        return True

def setup_logging(app=None):
    """Set up application logging with proper security practices."""
    
    # Create logs directory if it doesn't exist
    if not os.path.exists('logs'):
        os.makedirs('logs')
    
    # Configure log level
    log_level = os.getenv('LOG_LEVEL', 'INFO').upper()
    numeric_level = getattr(logging, log_level, logging.INFO)
    
    # Create formatter that doesn't expose sensitive information
    formatter = logging.Formatter(
        '%(asctime)s %(levelname)s: %(message)s [in %(pathname)s:%(lineno)d]'
    )
    
    # File handler with rotation to prevent huge log files
    file_handler = RotatingFileHandler(
        'logs/app.log', 
        maxBytes=10240000,  # 10MB
        backupCount=10
    )
    file_handler.setFormatter(formatter)
    file_handler.setLevel(numeric_level)
    file_handler.addFilter(SensitiveDataFilter())  # Add sensitive data filter
    
    # Console handler for development
    console_handler = logging.StreamHandler()
    console_handler.setFormatter(formatter)
    console_handler.setLevel(numeric_level)
    console_handler.addFilter(SensitiveDataFilter())  # Add sensitive data filter
    
    # Configure root logger
    root_logger = logging.getLogger()
    root_logger.setLevel(numeric_level)
    root_logger.addHandler(file_handler)
    
    # Only add console handler in development
    if os.getenv('FLASK_ENV') != 'production':
        root_logger.addHandler(console_handler)
    
    # Configure Flask app logger if provided
    if app:
        app.logger.addHandler(file_handler)
        if os.getenv('FLASK_ENV') != 'production':
            app.logger.addHandler(console_handler)
        app.logger.setLevel(numeric_level)
    
    # Filter out sensitive information from logs
    def filter_sensitive_data(record):
        """Filter out potentially sensitive information from log records."""
        # List of patterns that might contain sensitive data
        sensitive_patterns = [
            'password', 'token', 'key', 'secret', 'api_key',
            'authorization', 'session', 'cookie'
        ]
        
        message = str(record.getMessage()).lower()
        for pattern in sensitive_patterns:
            if pattern in message:
                record.msg = record.msg.replace(record.args[0] if record.args else '', '[REDACTED]')
                record.args = ()
                break
        
        return True
    
    # Add filter to all handlers
    log_filter = logging.Filter()
    log_filter.filter = filter_sensitive_data
    
    for handler in [file_handler, console_handler]:
        handler.addFilter(log_filter)
    
    return root_logger

def setup_security_logging():
    """Set up security-focused logging for the application"""
    
    # Create logs directory if it doesn't exist
    if not os.path.exists('logs'):
        os.makedirs('logs')
    
    # Configure log level
    log_level = os.getenv('LOG_LEVEL', 'INFO').upper()
    numeric_level = getattr(logging, log_level, logging.INFO)
    
    # Create security-focused formatter
    security_formatter = logging.Formatter(
        '%(asctime)s [%(levelname)s] %(name)s: %(message)s'
    )
    
    # Security file handler with rotation
    security_handler = RotatingFileHandler(
        'logs/security.log',
        maxBytes=int(os.getenv('LOG_MAX_SIZE', '10485760')),  # 10MB default
        backupCount=int(os.getenv('LOG_BACKUP_COUNT', '5'))
    )
    security_handler.setFormatter(security_formatter)
    security_handler.setLevel(numeric_level)
    security_handler.addFilter(SensitiveDataFilter())
    
    # Create security logger
    security_logger = logging.getLogger('moneymaker.security')
    security_logger.setLevel(numeric_level)
    security_logger.addHandler(security_handler)
    
    # Also log to console in development
    if os.getenv('FLASK_ENV') != 'production':
        console_handler = logging.StreamHandler()
        console_handler.setFormatter(security_formatter)
        console_handler.setLevel(numeric_level)
        console_handler.addFilter(SensitiveDataFilter())
        security_logger.addHandler(console_handler)
    
    return security_logger

def log_request(request, response_code=None, error=None):
    """Log HTTP requests in a secure manner."""
    logger = logging.getLogger(__name__)
    
    # Log basic request info without sensitive data
    client_ip = request.environ.get('HTTP_X_FORWARDED_FOR', request.remote_addr)
    method = request.method
    path = request.path
    user_agent = request.headers.get('User-Agent', 'Unknown')[:100]  # Truncate
    
    # Don't log query parameters as they might contain sensitive data
    log_message = f"{client_ip} - {method} {path} - Status: {response_code or 'Unknown'}"
    
    if error:
        logger.error(f"{log_message} - Error: {str(error)[:200]}")  # Truncate error
    else:
        logger.info(log_message)

def log_security_event(event_type, details, client_ip=None):
    """Log security-related events."""
    logger = logging.getLogger('security')
    
    timestamp = datetime.utcnow().isoformat()
    ip_info = f" from {client_ip}" if client_ip else ""
    
    logger.warning(f"SECURITY EVENT [{timestamp}]: {event_type}{ip_info} - {details}")
