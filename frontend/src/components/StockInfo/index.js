import React from 'react';
import styled from 'styled-components';

const StockInfoContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  width: 100%;
  margin-bottom: 2rem;
`;

const StockHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
`;

const StockTitle = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
`;

const TickerSymbol = styled.h1`
  font-size: 2.5rem;
  font-weight: 700;
  margin: 0;
`;

const CompanyName = styled.h2`
  font-size: 1.25rem;
  font-weight: 500;
  color: var(--text-secondary);
  margin: 0;
`;

const LiveBadge = styled.div`
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  background: ${props => props.isLive ? 'rgba(16, 185, 129, 0.1)' : 'rgba(239, 68, 68, 0.1)'};
  border-radius: 2rem;
  
  span {
    color: ${props => props.isLive ? 'var(--success)' : 'var(--danger)'};
    font-weight: 600;
    font-size: 0.875rem;
  }
  
  .indicator {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: ${props => props.isLive ? 'var(--success)' : 'var(--danger)'};
  }
`;

const PriceSection = styled.div`
  display: flex;
  align-items: center;
  gap: 1.5rem;
`;

const CurrentPrice = styled.div`
  font-size: 2.5rem;
  font-weight: 700;
`;

const PriceChange = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  
  .change {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 1.25rem;
    font-weight: 600;
    color: ${props => props.isPositive ? 'var(--success)' : 'var(--danger)'};
    
    .arrow {
      font-size: 1.5rem;
    }
  }
  
  .change-pct {
    font-size: 1rem;
    color: var(--text-secondary);
  }
`;

const AdditionalInfo = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 1rem;
  width: 100%;
`;

const InfoCard = styled.div`
  background: var(--bg-glass);
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  
  .label {
    font-size: 0.875rem;
    color: var(--text-muted);
  }
  
  .value {
    font-size: 1.25rem;
    font-weight: 600;
    color: var(--text-primary);
  }
  
  &.ai-analysis {
    grid-column: 1 / -1;
    
    .value {
      font-size: 1rem;
      font-weight: 500;
      line-height: 1.6;
    }
  }
`;

const StockInfo = ({ stockData, isLiveConnected }) => {
  if (!stockData) return null;
  
  const isPositiveChange = stockData.change >= 0;
  
  const formatCurrency = (value) => {
    if (value === undefined || value === null) return 'N/A';
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    }).format(value);
  };
  
  return (
    <StockInfoContainer>
      <StockHeader>
        <StockTitle>
          <TickerSymbol>{stockData.ticker}</TickerSymbol>
          <CompanyName>{stockData.company_name}</CompanyName>
        </StockTitle>
        <LiveBadge isLive={isLiveConnected}>
          <div className="indicator"></div>
          <span>{isLiveConnected ? 'LIVE' : 'POLLING'}</span>
        </LiveBadge>
      </StockHeader>
      
      <PriceSection>
        <CurrentPrice>{formatCurrency(stockData.price)}</CurrentPrice>
        <PriceChange isPositive={isPositiveChange}>
          <div className="change">
            <span className="arrow">{isPositiveChange ? '↑' : '↓'}</span>
            {formatCurrency(Math.abs(stockData.change))}
          </div>
          <div className="change-pct">
            ({stockData.change_pct})
          </div>
        </PriceChange>
      </PriceSection>
      
      <AdditionalInfo>
        <InfoCard>
          <div className="label">Open</div>
          <div className="value">{formatCurrency(stockData.open)}</div>
        </InfoCard>
        <InfoCard>
          <div className="label">Day Range</div>
          <div className="value">{formatCurrency(stockData.low)} - {formatCurrency(stockData.high)}</div>
        </InfoCard>
        <InfoCard>
          <div className="label">Volume</div>
          <div className="value">{stockData.volume}</div>
        </InfoCard>
        <InfoCard>
          <div className="label">Market Cap</div>
          <div className="value">{stockData.market_cap}</div>
        </InfoCard>
        <InfoCard>
          <div className="label">P/E Ratio</div>
          <div className="value">{stockData.pe_ratio}</div>
        </InfoCard>
        <InfoCard>
          <div className="label">52 Week Range</div>
          <div className="value">{formatCurrency(stockData.low_52w)} - {formatCurrency(stockData.high_52w)}</div>
        </InfoCard>
        <InfoCard>
          <div className="label">Dividend Yield</div>
          <div className="value">{stockData.dividend_yield}</div>
        </InfoCard>
        <InfoCard>
          <div className="label">AI Prediction</div>
          <div className="value">
            {formatCurrency(stockData.predicted_price)} 
            <span style={{ fontSize: '0.875rem', marginLeft: '0.5rem', color: 'var(--text-muted)' }}>
              ({stockData.prediction_confidence.toFixed(1)}% confidence)
            </span>
          </div>
        </InfoCard>
        
        {stockData.ai_analysis && (
          <InfoCard className="ai-analysis">
            <div className="label">AI Analysis</div>
            <div className="value">{stockData.ai_analysis}</div>
          </InfoCard>
        )}
      </AdditionalInfo>
    </StockInfoContainer>
  );
};

export default StockInfo;
