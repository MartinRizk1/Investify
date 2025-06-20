import React, { useState, useEffect, useRef } from 'react';
import { Line } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  TimeScale,
  Filler
} from 'chart.js';
import 'chartjs-adapter-luxon';
import annotationPlugin from 'chartjs-plugin-annotation';
import styled from 'styled-components';

// Register Chart.js components
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  TimeScale,
  Filler,
  annotationPlugin
);

const ChartContainer = styled.div`
  position: relative;
  width: 100%;
  height: 400px;
  margin-bottom: 2rem;
  background: var(--bg-glass);
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 1rem;
`;

const TimeframeSelector = styled.div`
  display: flex;
  justify-content: center;
  gap: 0.5rem;
  margin-bottom: 1rem;
`;

const TimeframeButton = styled.button`
  background: ${props => props.active ? 'var(--purple-accent)' : 'var(--bg-secondary)'};
  border: 1px solid ${props => props.active ? 'var(--purple-secondary)' : 'transparent'};
  color: var(--text-primary);
  border-radius: 6px;
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  
  &:hover {
    background: ${props => props.active ? 'var(--purple-accent)' : 'var(--bg-tertiary)'};
  }
`;

// Generate mock historical data for demo purposes
const generateHistoricalData = (days, basePrice, volatility) => {
  const data = [];
  const now = new Date();
  
  for (let i = days; i >= 0; i--) {
    const date = new Date(now);
    date.setDate(date.getDate() - i);
    
    // Generate random price movement with some trend
    const priceChange = (Math.random() - 0.48) * volatility;
    const newPrice = i === days ? basePrice : data[days - i - 1].price * (1 + priceChange);
    
    data.push({
      time: date.toISOString(),
      price: parseFloat(newPrice.toFixed(2))
    });
  }
  
  return data;
};

const StockChart = ({ ticker, currentPrice, predictedPrice }) => {
  const [timeframe, setTimeframe] = useState('1M'); // Default to 1 month
  const [chartData, setChartData] = useState(null);
  const chartRef = useRef(null);
  
  useEffect(() => {
    if (!ticker || !currentPrice) return;
    
    // In a real application, you would fetch historical data from API
    // For this demo, we'll generate mock data
    const daysMap = {
      '1D': 1,
      '1W': 7,
      '1M': 30,
      '3M': 90,
      '6M': 180,
      '1Y': 365,
    };
    
    const days = daysMap[timeframe];
    const volatilityMap = {
      '1D': 0.005,
      '1W': 0.01,
      '1M': 0.02,
      '3M': 0.03,
      '6M': 0.04,
      '1Y': 0.05,
    };
    
    const historicalData = generateHistoricalData(days, currentPrice, volatilityMap[timeframe]);
    
    // Format data for Chart.js
    const formattedData = {
      labels: historicalData.map(data => data.time),
      datasets: [{
        label: ticker,
        data: historicalData.map(data => data.price),
        borderColor: '#8b5cf6',
        borderWidth: 2,
        pointRadius: 0,
        tension: 0.2,
        fill: true,
        backgroundColor: 'rgba(139, 92, 246, 0.1)',
      }]
    };
    
    setChartData(formattedData);
  }, [ticker, timeframe, currentPrice]);
  
  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    scales: {
      x: {
        type: 'time',
        time: {
          unit: timeframe === '1D' ? 'hour' : timeframe === '1W' ? 'day' : 'day',
          tooltipFormat: timeframe === '1D' ? 'HH:mm' : 'MMM d',
          displayFormats: {
            hour: 'HH:mm',
            day: 'MMM d',
            week: 'MMM d',
            month: 'MMM yy'
          }
        },
        grid: {
          display: true,
          color: 'rgba(255, 255, 255, 0.05)',
        },
        ticks: {
          color: '#9ca3af',
          maxRotation: 0,
        },
        border: {
          display: false
        }
      },
      y: {
        position: 'right',
        grid: {
          display: true,
          color: 'rgba(255, 255, 255, 0.05)',
        },
        ticks: {
          color: '#9ca3af',
          callback: function(value) {
            return '$' + value.toFixed(2);
          }
        },
        border: {
          display: false
        }
      }
    },
    plugins: {
      legend: {
        display: false,
      },
      tooltip: {
        mode: 'index',
        intersect: false,
        backgroundColor: 'rgba(26, 26, 26, 0.9)',
        titleColor: '#ffffff',
        bodyColor: '#e5e7eb',
        borderColor: 'rgba(139, 92, 246, 0.3)',
        borderWidth: 1,
        padding: 10,
        displayColors: false,
        callbacks: {
          label: function(context) {
            return `$${context.raw.toFixed(2)}`;
          }
        }
      },
      annotation: predictedPrice ? {
        annotations: {
          prediction: {
            type: 'line',
            yMin: predictedPrice,
            yMax: predictedPrice,
            borderColor: 'rgba(16, 185, 129, 0.7)',
            borderWidth: 2,
            borderDash: [5, 5],
            label: {
              display: true,
              content: `AI Prediction: $${predictedPrice.toFixed(2)}`,
              position: 'start',
              backgroundColor: 'rgba(16, 185, 129, 0.8)',
              color: '#ffffff',
              font: {
                weight: 'bold'
              },
              padding: 6
            }
          }
        }
      } : {}
    },
    interaction: {
      mode: 'index',
      intersect: false,
    },
    animation: {
      duration: 1000
    }
  };
  
  return (
    <ChartContainer>
      <TimeframeSelector>
        {['1D', '1W', '1M', '3M', '6M', '1Y'].map(tf => (
          <TimeframeButton
            key={tf}
            active={timeframe === tf}
            onClick={() => setTimeframe(tf)}
          >
            {tf}
          </TimeframeButton>
        ))}
      </TimeframeSelector>
      
      {chartData && (
        <Line
          ref={chartRef}
          data={chartData}
          options={chartOptions}
        />
      )}
    </ChartContainer>
  );
};

export default StockChart;
