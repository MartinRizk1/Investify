import React, { useState, useEffect } from 'react';
import styled from 'styled-components';
import SearchBar from './components/SearchBar';
import StockChart from './components/StockChart';
import StockInfo from './components/StockInfo';
import TechnicalIndicators from './components/TechnicalIndicators';
import GlobalStyles from './components/GlobalStyles';
import stockService from './services/stockService';

const AppContainer = styled.div`
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
  animation: fadeIn 1s ease-out;
`;

const Header = styled.header`
  display: flex;
  justify-content: center;
  margin-bottom: 2rem;
  text-align: center;
`;

const Logo = styled.h1`
  font-size: 2.5rem;
  font-weight: 900;
  background: linear-gradient(135deg, var(--purple-primary) 0%, var(--purple-secondary) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin-bottom: 0.5rem;
  
  span {
    font-weight: 500;
  }
`;

const Tagline = styled.p`
  color: var(--text-secondary);
  font-size: 1rem;
`;

const LoadingOverlay = styled.div`
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
  animation: fadeIn 0.3s;
  
  .spinner {
    width: 50px;
    height: 50px;
    border: 4px solid rgba(139, 92, 246, 0.3);
    border-top: 4px solid var(--purple-primary);
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }
  
  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }
`;

const ErrorMessage = styled.div`
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  background: var(--danger-bg);
  border: 1px solid var(--danger);
  border-radius: 12px;
  padding: 1rem;
  color: var(--text-primary);
  display: flex;
  align-items: center;
  gap: 1rem;
  z-index: 100;
  animation: slideInDown 0.3s;
  max-width: 90%;
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
  
  .error-icon {
    font-size: 1.5rem;
  }
  
  button {
    background: var(--danger);
    color: white;
    border: none;
    border-radius: 6px;
    padding: 0.5rem 0.75rem;
    margin-left: auto;
    cursor: pointer;
    font-weight: 600;
    font-size: 0.875rem;
  }
  
  @keyframes slideInDown {
    from {
      transform: translate(-50%, -50px);
      opacity: 0;
    }
    to {
      transform: translate(-50%, 0);
      opacity: 1;
    }
  }
`;

function App() {
  const [stockData, setStockData] = useState(null);
  const [technicalData, setTechnicalData] = useState(null);
  const [loading, setLoading] = useState(false);
  const [isLiveConnected, setIsLiveConnected] = useState(false);
  const [currentTicker, setCurrentTicker] = useState('');
  const [error, setError] = useState(null);
  
  useEffect(() => {
    // Clean up WebSocket connection when component unmounts
    return () => {
      stockService.disconnectWebSocket();
    };
  }, []);
  
  // Add event listener for WebSocket events
  useEffect(() => {
    if (!currentTicker) return;
    
    const handleStockUpdate = (event) => {
      if (event.type === 'connection') {
        setIsLiveConnected(event.status);
      } else if (event.type === 'data') {
        // Update stock data from WebSocket
        const data = event.data;
        
        // Update stock price and change information
        if (stockData) {
          setStockData(prevData => ({
            ...prevData,
            price: data.price,
            change: data.change,
            change_pct: data.change_pct,
          }));
        }
        
        // Update technical indicators if available
        if (data.technical) {
          setTechnicalData(data.technical);
        }
      }
    };
    
    stockService.addListener(handleStockUpdate);
    
    // Clean up listener when component unmounts or ticker changes
    return () => {
      stockService.removeListener(handleStockUpdate);
    };
  }, [currentTicker, stockData]);
  
  const handleSearch = async (ticker) => {
    if (!ticker) return;
    
    setLoading(true);
    setCurrentTicker(ticker);
    
    try {
      // Disconnect any existing WebSocket connection
      stockService.disconnectWebSocket();
      setIsLiveConnected(false);
      
      // Fetch initial stock data from API
      const data = await stockService.fetchStockData(ticker);
      setStockData(data);
      
      if (data.technical) {
        setTechnicalData(data.technical);
      }
      
      // Connect to WebSocket for real-time updates
      stockService.connectWebSocket(ticker);
    } catch (error) {
      console.error('Error fetching stock data:', error);
      // Handle error state
      setStockData(null);
      setTechnicalData(null);
      setError(error.message || 'Failed to fetch stock data. Please try again.');
    } finally {
      setLoading(false);
    }
  };
  
  return (
    <>
      <GlobalStyles />
      <AppContainer>
        <Header>
          <div>
            <Logo>Investify<span> AI</span></Logo>
            <Tagline>AI-Powered Stock Analysis & Technical Indicators</Tagline>
          </div>
        </Header>
        
        <SearchBar onSearch={handleSearch} isLoading={loading} />
        
        {stockData && (
          <>
            <StockInfo stockData={stockData} isLiveConnected={isLiveConnected} />
            <StockChart 
              ticker={stockData.ticker} 
              currentPrice={stockData.price}
              predictedPrice={stockData.predicted_price}
            />
            {technicalData && (
              <TechnicalIndicators technicalData={technicalData} />
            )}
          </>
        )}
        
        {loading && (
          <LoadingOverlay>
            <div className="spinner"></div>
          </LoadingOverlay>
        )}
        
        {error && !loading && (
          <ErrorMessage>
            <div className="error-icon">⚠️</div>
            <div>{error}</div>
            <button onClick={() => setError(null)}>Dismiss</button>
          </ErrorMessage>
        )}
      </AppContainer>
    </>
  );
}

export default App;
