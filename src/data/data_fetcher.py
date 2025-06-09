import yfinance as yf
import pandas as pd
import requests
import re

class DataFetcher:
    """Fetches historical market data and company info for a given ticker."""
    
    def __init__(self):
        # Common company name to ticker mappings
        self.company_mapping = {
            # Technology
            'apple': 'AAPL',
            'microsoft': 'MSFT',
            'google': 'GOOGL',
            'alphabet': 'GOOGL',
            'amazon': 'AMZN',
            'tesla': 'TSLA',
            'meta': 'META',
            'facebook': 'META',
            'netflix': 'NFLX',
            'nvidia': 'NVDA',
            'salesforce': 'CRM',
            'adobe': 'ADBE',
            'intel': 'INTC',
            'oracle': 'ORCL',
            'ibm': 'IBM',
            'cisco': 'CSCO',
            'uber': 'UBER',
            'airbnb': 'ABNB',
            'spotify': 'SPOT',
            'zoom': 'ZM',
            'paypal': 'PYPL',
            'square': 'SQ',
            'shopify': 'SHOP',
            'snowflake': 'SNOW',
            'palantir': 'PLTR',
            
            # Retail & Consumer
            'walmart': 'WMT',
            'target': 'TGT',
            'costco': 'COST',
            'home depot': 'HD',
            'lowes': 'LOW',
            'starbucks': 'SBUX',
            'mcdonald': 'MCD',
            'nike': 'NKE',
            'adidas': 'ADDYY',
            'coca cola': 'KO',
            'pepsi': 'PEP',
            'general motors': 'GM',
            'ford': 'F',
            'disney': 'DIS',
            
            # Healthcare & Pharmaceuticals
            'johnson & johnson': 'JNJ',
            'pfizer': 'PFE',
            'abbvie': 'ABBV',
            'merck': 'MRK',
            'bristol myers': 'BMY',
            'eli lilly': 'LLY',
            'unitedhealth': 'UNH',
            'cvs': 'CVS',
            'walgreens': 'WBA',
            'rite aid': 'RAD',
            'moderna': 'MRNA',
            'novavax': 'NVAX',
            
            # Financial Services
            'visa': 'V',
            'mastercard': 'MA',
            'jpmorgan': 'JPM',
            'bank of america': 'BAC',
            'wells fargo': 'WFC',
            'goldman sachs': 'GS',
            'morgan stanley': 'MS',
            'american express': 'AXP',
            'berkshire hathaway': 'BRK-B',
            'blackrock': 'BLK',
            
            # Energy & Utilities
            'exxon': 'XOM',
            'chevron': 'CVX',
            'conocophillips': 'COP',
            'marathon petroleum': 'MPC',
            'kinder morgan': 'KMI',
            'nextera energy': 'NEE',
            
            # Consumer Goods
            'procter & gamble': 'PG',
            'unilever': 'UL',
            'colgate': 'CL',
            'kimberly clark': 'KMB',
            'general mills': 'GIS',
            'kraft heinz': 'KHC',
            
            # Telecommunications
            'verizon': 'VZ',
            'at&t': 'T',
            'comcast': 'CMCSA',
            't-mobile': 'TMUS',
            
            # Airlines & Travel
            'american airlines': 'AAL',
            'delta': 'DAL',
            'united airlines': 'UAL',
            'southwest': 'LUV',
            'boeing': 'BA',
            'airbus': 'EADSY',
            
            # Real Estate & REITs
            'realty income': 'O',
            'american tower': 'AMT',
            'prologis': 'PLD',
            'simon property': 'SPG'
        }
    
    def get_ticker(self, company_name: str) -> str:
        """Convert company name to ticker symbol using comprehensive search."""
        if not company_name:
            return ""
            
        # Clean the input
        clean_name = company_name.strip().lower()
        original_input = company_name.strip()
        
        # Remove common words that might interfere with matching
        clean_name = re.sub(r'\b(inc|corp|corporation|company|co|ltd|llc|group|holdings|plc)\b', '', clean_name).strip()
        
        # First, check if it's already a ticker (typically 1-5 uppercase letters)
        if re.match(r'^[A-Z]{1,5}(-[A-Z])?$', original_input.upper()):
            # Validate it's a real ticker by trying to fetch info
            try:
                test_ticker = yf.Ticker(original_input.upper())
                info = test_ticker.info
                if info and info.get('symbol') and info.get('longName'):
                    return original_input.upper()
            except:
                pass
        
        # Apply fuzzy matching for common typos before checking mappings
        clean_name = self._fix_common_typos(clean_name)
        
        # Check our predefined mappings with exact and partial matches
        for company, ticker in self.company_mapping.items():
            # Exact match
            if clean_name == company.lower():
                return ticker
            # Partial match - company name contains our search term
            if clean_name in company.lower() or company.lower() in clean_name:
                # Make sure it's a meaningful match (at least 3 characters)
                if len(clean_name) >= 3 and any(word in company.lower() for word in clean_name.split() if len(word) >= 3):
                    return ticker
        
        # Use Yahoo Finance search to find the ticker dynamically
        return self._search_yahoo_ticker(original_input, clean_name)
    
    def _fix_common_typos(self, name: str) -> str:
        """Fix common typos in company names."""
        typo_corrections = {
            'wallgreens': 'walgreens',
            'wallgreen': 'walgreens',
            'wallgrens': 'walgreens',
            'walgrens': 'walgreens',
            'walgreen': 'walgreens',
            'mcdonalds': 'mcdonald',
            'macdonald': 'mcdonald',
            'macdonalds': 'mcdonald',
            'googel': 'google',
            'gogle': 'google',
            'microsft': 'microsoft',
            'micosoft': 'microsoft',
            'amzon': 'amazon',
            'amazom': 'amazon',
            'teslas': 'tesla',
            'tesela': 'tesla',
            'appel': 'apple',
            'aple': 'apple',
            'starbuck': 'starbucks',
            'starbuks': 'starbucks',
            'walmarts': 'walmart',
            'walmart': 'walmart',
            'pepsi cola': 'pepsi',
            'coca-cola': 'coca cola',
            'coke': 'coca cola',
            'johnson and johnson': 'johnson & johnson',
            'johnson & johnsons': 'johnson & johnson',
            'jp morgan': 'jpmorgan',
            'jpmorgan chase': 'jpmorgan',
            'berkshire': 'berkshire hathaway',
            'berkshire hathway': 'berkshire hathaway',
            'home depo': 'home depot',
            'homedepot': 'home depot',
            'lowe\'s': 'lowes',
            'lowes': 'lowes',
            'target corp': 'target',
            'targets': 'target',
            'costcos': 'costco',
            'cost co': 'costco',
            'exxonmobil': 'exxon',
            'exxon mobil': 'exxon',
            'chevron corp': 'chevron',
            'general motor': 'general motors',
            'gm': 'general motors',
            'ford motor': 'ford',
            'ford motors': 'ford',
            'meta platforms': 'meta',
            'facebook inc': 'meta',
            'alphabet inc': 'alphabet',
            'amazon.com': 'amazon',
            'nike inc': 'nike',
            'the walt disney': 'disney',
            'walt disney': 'disney',
            'disney company': 'disney',
            'american express company': 'american express',
            'amex': 'american express',
            'bank america': 'bank of america',
            'bofa': 'bank of america',
            'wells fargo bank': 'wells fargo',
            'goldman sachs group': 'goldman sachs',
            'morgan stanley group': 'morgan stanley',
            'jpmorgan chase': 'jpmorgan',
            'visa inc': 'visa',
            'mastercard inc': 'mastercard',
            'procter and gamble': 'procter & gamble',
            'p&g': 'procter & gamble',
            'johnson controls': 'johnson & johnson',
            'verizon communications': 'verizon',
            'at&t inc': 'at&t',
            'att': 'at&t',
            'comcast corp': 'comcast',
            't mobile': 't-mobile',
            'tmobile': 't-mobile',
            'american airline': 'american airlines',
            'delta airlines': 'delta',
            'united airline': 'united airlines',
            'southwest airlines': 'southwest',
            'boeing company': 'boeing',
            'cvs health': 'cvs',
            'cvs pharmacy': 'cvs',
            'unitedhealth group': 'unitedhealth',
            'pfizer inc': 'pfizer',
            'merck & co': 'merck',
            'eli lilly and company': 'eli lilly',
            'bristol myers squibb': 'bristol myers',
            'abbvie inc': 'abbvie',
            'moderna inc': 'moderna',
            'novavax inc': 'novavax',
            'rite aid corp': 'rite aid',
            'nvidia corp': 'nvidia',
            'intel corp': 'intel',
            'oracle corp': 'oracle',
            'ibm corp': 'ibm',
            'cisco systems': 'cisco',
            'salesforce.com': 'salesforce',
            'adobe systems': 'adobe',
            'netflix inc': 'netflix',
            'uber technologies': 'uber',
            'airbnb inc': 'airbnb',
            'spotify technology': 'spotify',
            'zoom video': 'zoom',
            'paypal holdings': 'paypal',
            'square inc': 'square',
            'shopify inc': 'shopify',
            'snowflake inc': 'snowflake',
            'palantir technologies': 'palantir'
        }
        
        # Apply corrections
        for typo, correction in typo_corrections.items():
            if typo in name:
                name = name.replace(typo, correction)
        
        return name
    
    def _search_yahoo_ticker(self, original_input: str, clean_name: str) -> str:
        """Search for ticker using Yahoo Finance's comprehensive database."""
        
        # Generate potential ticker variations
        potential_tickers = []
        
        # Add original input variations
        potential_tickers.extend([
            original_input.upper(),
            original_input.upper().replace(' ', ''),
            original_input.upper()[:4],
            original_input.upper()[:3]
        ])
        
        # Add word-based variations
        words = clean_name.split()
        if len(words) > 1:
            # First letters of each word
            potential_tickers.append(''.join([word[0] for word in words]).upper())
            # First few letters of first word + first letter of others
            if len(words[0]) >= 2:
                potential_tickers.append((words[0][:2] + ''.join([w[0] for w in words[1:]])).upper())
            if len(words[0]) >= 3:
                potential_tickers.append((words[0][:3] + ''.join([w[0] for w in words[1:]])).upper())
        
        # Add common patterns for specific companies (using corrected names)
        if 'walgreens' in clean_name:
            potential_tickers.extend(['WBA', 'WG'])
        if 'walmart' in clean_name:
            potential_tickers.extend(['WMT'])
        if 'mcdonald' in clean_name:
            potential_tickers.extend(['MCD'])
        if 'coca cola' in clean_name or 'coke' in clean_name:
            potential_tickers.extend(['KO'])
        if 'johnson' in clean_name and ('&' in clean_name or 'and' in clean_name):
            potential_tickers.extend(['JNJ'])
        if 'berkshire' in clean_name:
            potential_tickers.extend(['BRK-B', 'BRK-A'])
        if 'general motors' in clean_name or clean_name == 'gm':
            potential_tickers.extend(['GM'])
        if 'ford' in clean_name:
            potential_tickers.extend(['F'])
        if 'apple' in clean_name:
            potential_tickers.extend(['AAPL'])
        if 'microsoft' in clean_name:
            potential_tickers.extend(['MSFT'])
        if 'tesla' in clean_name:
            potential_tickers.extend(['TSLA'])
        if 'google' in clean_name or 'alphabet' in clean_name:
            potential_tickers.extend(['GOOGL', 'GOOG'])
        if 'amazon' in clean_name:
            potential_tickers.extend(['AMZN'])
        if 'meta' in clean_name or 'facebook' in clean_name:
            potential_tickers.extend(['META'])
        
        # Test each potential ticker
        for ticker in potential_tickers:
            if len(ticker) >= 1 and len(ticker) <= 5:
                try:
                    test_ticker = yf.Ticker(ticker)
                    info = test_ticker.info
                    
                    # Check if we got valid company data
                    if (info and 
                        info.get('symbol') == ticker and 
                        (info.get('longName') or info.get('shortName')) and
                        info.get('regularMarketPrice') is not None):
                        
                        # Verify the company name matches our search
                        company_names = [
                            info.get('longName', '').lower(),
                            info.get('shortName', '').lower(),
                            info.get('displayName', '').lower()
                        ]
                        
                        # Check if any company name contains words from our search
                        search_words = [word for word in clean_name.split() if len(word) > 2]
                        for company_name_check in company_names:
                            if company_name_check:
                                # Check for word matches
                                if any(search_word in company_name_check for search_word in search_words):
                                    return ticker
                                # Also check reverse - if company name words are in our search
                                company_words = [word for word in company_name_check.split() if len(word) > 2]
                                if any(comp_word in clean_name for comp_word in company_words):
                                    return ticker
                
                except Exception:
                    continue
        
        # If no match found, try a broader search with Yahoo Finance search API
        try:
            # Use yfinance's Ticker.search if available, or try common patterns
            search_terms = [original_input, clean_name]
            for term in search_terms:
                # Try exact term as ticker
                test_ticker = yf.Ticker(term.upper())
                info = test_ticker.info
                if (info and 
                    info.get('symbol') and 
                    info.get('longName') and 
                    info.get('regularMarketPrice') is not None):
                    return term.upper()
        except:
            pass
        
        # Last resort: return the original input uppercased if it looks like a ticker
        if re.match(r'^[A-Z]{1,5}$', original_input.upper()):
            return original_input.upper()
        
        # Return None to indicate not found
        return None

    def fetch_data(self, company_name: str, period: str = '1y', interval: str = '1d') -> pd.DataFrame:
        """Fetch historical stock data for any valid company."""
        ticker = self.get_ticker(company_name)
        
        if not ticker:
            print(f"Could not resolve ticker for: {company_name}")
            return pd.DataFrame()
            
        try:
            # Try to fetch data with the resolved ticker
            data = yf.download(ticker, period=period, interval=interval, progress=False)
            
            # Handle MultiIndex columns from yfinance
            if isinstance(data.columns, pd.MultiIndex):
                # Flatten the columns by taking only the first level (price type)
                data.columns = data.columns.get_level_values(0)
            
            if data.empty:
                # Try with different periods if data is empty
                for alt_period in ['6mo', '3mo', '1mo', '5d']:
                    print(f"Trying alternative period: {alt_period}")
                    data = yf.download(ticker, period=alt_period, interval=interval, progress=False)
                    # Handle MultiIndex columns for retry attempts too
                    if isinstance(data.columns, pd.MultiIndex):
                        data.columns = data.columns.get_level_values(0)
                    if not data.empty:
                        break
                        
            # Validate we have the expected columns
            expected_columns = ['Open', 'High', 'Low', 'Close', 'Volume']
            if data.empty or not all(col in data.columns for col in expected_columns):
                print(f"Invalid data structure for {ticker}")
                return pd.DataFrame()
                
            return data
            
        except Exception as e:
            print(f"Error fetching data for {ticker}: {e}")
            return pd.DataFrame()

    def fetch_company_info(self, company_name: str) -> dict:
        """Fetch comprehensive company information from Yahoo Finance."""
        ticker = self.get_ticker(company_name)
        
        if not ticker:
            return {}
            
        try:
            ticker_obj = yf.Ticker(ticker)
            info = ticker_obj.info
            
            # Validate that we got meaningful data
            if not info or not info.get('symbol'):
                return {}
                
            # Enhance info with additional data if available
            try:
                # Get additional financial data
                financials = ticker_obj.financials
                balance_sheet = ticker_obj.balance_sheet
                cashflow = ticker_obj.cashflow
                
                # Add summary financial metrics if available
                if not financials.empty:
                    info['has_financials'] = True
                if not balance_sheet.empty:
                    info['has_balance_sheet'] = True
                if not cashflow.empty:
                    info['has_cashflow'] = True
                    
            except Exception:
                pass  # Additional data is optional
                
            return info
            
        except Exception as e:
            print(f"Error fetching company info for {ticker}: {e}")
            return {}

    def fetch_news(self, company_name: str) -> list:
        """Fetch recent news for the company."""
        ticker = self.get_ticker(company_name)
        try:
            ticker_obj = yf.Ticker(ticker)
            news = ticker_obj.news
            return news if news else []
        except Exception as e:
            print(f"Error fetching news for {ticker}: {e}")
            return []
