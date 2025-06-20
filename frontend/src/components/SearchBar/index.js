import React, { useState } from 'react';
import styled from 'styled-components';

const SearchContainer = styled.div`
  display: flex;
  justify-content: center;
  margin-bottom: 2rem;
  width: 100%;
`;

const SearchForm = styled.form`
  display: flex;
  width: 100%;
  max-width: 600px;
  gap: 0.5rem;
  position: relative;
`;

const SearchInput = styled.input`
  flex: 1;
  padding: 1rem 1.5rem;
  border-radius: 12px;
  border: 2px solid rgba(139, 92, 246, 0.3);
  background: rgba(26, 26, 26, 0.8);
  color: #ffffff;
  font-size: 1rem;
  font-weight: 500;
  transition: all 0.2s ease-in-out;
  box-shadow: 0 0 0 rgba(139, 92, 246, 0);
  
  &:focus {
    outline: none;
    border-color: rgba(139, 92, 246, 0.8);
    box-shadow: 0 0 20px rgba(139, 92, 246, 0.2);
  }
  
  &::placeholder {
    color: #9ca3af;
  }
`;

const SearchButton = styled.button`
  padding: 0 1.5rem;
  border-radius: 12px;
  border: none;
  background: linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%);
  color: #ffffff;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease-in-out;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 5px 15px rgba(139, 92, 246, 0.4);
  }
  
  &:disabled {
    opacity: 0.7;
    cursor: not-allowed;
    transform: none;
  }
  
  &.loading {
    position: relative;
    
    &:after {
      content: "";
      position: absolute;
      width: 20px;
      height: 20px;
      border: 3px solid rgba(255, 255, 255, 0.3);
      border-top-color: #ffffff;
      border-radius: 50%;
      animation: loader-spin 1s linear infinite;
    }
  }
  
  @keyframes loader-spin {
    to {
      transform: rotate(360deg);
    }
  }
`;

const ExampleTickers = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-top: 0.5rem;
  justify-content: center;
`;

const ExampleTicker = styled.button`
  background: rgba(139, 92, 246, 0.1);
  border: 1px solid rgba(139, 92, 246, 0.3);
  color: #e5e7eb;
  border-radius: 6px;
  padding: 0.25rem 0.75rem;
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &:hover {
    background: rgba(139, 92, 246, 0.2);
  }
`;

const SearchBar = ({ onSearch, isLoading }) => {
  const [ticker, setTicker] = useState('');
  const popularTickers = ['AAPL', 'MSFT', 'GOOGL', 'AMZN', 'META', 'TSLA'];
  
  const handleSubmit = (e) => {
    e.preventDefault();
    if (ticker.trim()) {
      onSearch(ticker.trim().toUpperCase());
    }
  };
  
  const handleExampleClick = (exampleTicker) => {
    setTicker(exampleTicker);
    onSearch(exampleTicker);
  };
  
  return (
    <SearchContainer>
      <div>
        <SearchForm onSubmit={handleSubmit}>
          <SearchInput
            type="text"
            placeholder="Enter stock symbol (e.g., AAPL)"
            value={ticker}
            onChange={(e) => setTicker(e.target.value)}
            disabled={isLoading}
          />
          <SearchButton type="submit" disabled={isLoading || !ticker.trim()} className={isLoading ? 'loading' : ''}>
            {isLoading ? '' : 'Search'}
          </SearchButton>
        </SearchForm>
        
        <ExampleTickers>
          {popularTickers.map(ticker => (
            <ExampleTicker 
              key={ticker} 
              onClick={() => handleExampleClick(ticker)}
              disabled={isLoading}
            >
              {ticker}
            </ExampleTicker>
          ))}
        </ExampleTickers>
      </div>
    </SearchContainer>
  );
};

export default SearchBar;
