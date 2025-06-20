import React, { useEffect, useState } from 'react';
import { Line } from 'react-chartjs-2';
import styled from 'styled-components';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  TimeScale
} from 'chart.js';

// Register Chart.js components
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  TimeScale
);

const TechnicalContainer = styled.div`
  display: grid;
  grid-template-columns: 1fr;
  gap: 1.5rem;
  margin-bottom: 2rem;
  
  @media (min-width: 1200px) {
    grid-template-columns: repeat(2, 1fr);
  }
`;

const ChartCard = styled.div`
  background: var(--bg-glass);
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 1rem;
  height: 300px;
  
  h3 {
    font-size: 1.25rem;
    margin-bottom: 1rem;
    color: var(--text-secondary);
  }
  
  &.full-width {
    grid-column: 1 / -1;
  }
`;

const NoData = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--text-muted);
  font-style: italic;
`;

const TechnicalIndicators = ({ technicalData }) => {
  const [rsiChartData, setRsiChartData] = useState(null);
  const [macdChartData, setMacdChartData] = useState(null);
  const [bbChartData, setBbChartData] = useState(null);
  
  useEffect(() => {
    if (!technicalData || !technicalData.rsi || !technicalData.macd || !technicalData.bollinger) return;
    
    // RSI Chart Data
    setRsiChartData({
      labels: technicalData.rsi.dates,
      datasets: [{
        label: 'RSI',
        data: technicalData.rsi.values,
        borderColor: '#8b5cf6',
        borderWidth: 2,
        tension: 0.4,
        pointRadius: 0
      }, {
        label: 'Overbought (70)',
        data: Array(technicalData.rsi.dates.length).fill(70),
        borderColor: 'rgba(239, 68, 68, 0.5)',
        borderWidth: 1,
        borderDash: [5, 5],
        pointRadius: 0
      }, {
        label: 'Oversold (30)',
        data: Array(technicalData.rsi.dates.length).fill(30),
        borderColor: 'rgba(16, 185, 129, 0.5)',
        borderWidth: 1,
        borderDash: [5, 5],
        pointRadius: 0
      }]
    });
    
    // MACD Chart Data
    setMacdChartData({
      labels: technicalData.macd.dates,
      datasets: [{
        label: 'MACD',
        data: technicalData.macd.macd,
        borderColor: '#8b5cf6',
        borderWidth: 2,
        tension: 0.4,
        pointRadius: 0
      }, {
        label: 'Signal',
        data: technicalData.macd.signal,
        borderColor: '#f59e0b',
        borderWidth: 2,
        tension: 0.4,
        pointRadius: 0
      }]
    });
    
    // Bollinger Bands Chart Data
    setBbChartData({
      labels: technicalData.bollinger.dates,
      datasets: [{
        label: 'Middle Band (SMA)',
        data: technicalData.bollinger.middle,
        borderColor: '#e5e7eb',
        borderWidth: 2,
        tension: 0.4,
        pointRadius: 0
      }, {
        label: 'Upper Band',
        data: technicalData.bollinger.upper,
        borderColor: 'rgba(16, 185, 129, 0.8)',
        borderWidth: 2,
        borderDash: [5, 5],
        tension: 0.4,
        pointRadius: 0,
        fill: false
      }, {
        label: 'Lower Band',
        data: technicalData.bollinger.lower,
        borderColor: 'rgba(239, 68, 68, 0.8)',
        borderWidth: 2,
        borderDash: [5, 5],
        tension: 0.4,
        pointRadius: 0,
        fill: false
      }, {
        label: 'Price',
        data: technicalData.bollinger.price,
        borderColor: '#8b5cf6',
        borderWidth: 2,
        pointRadius: 0,
        tension: 0.4
      }]
    });
  }, [technicalData]);
  
  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    scales: {
      x: {
        display: true,
        grid: {
          display: false,
          color: 'rgba(255, 255, 255, 0.05)'
        },
        ticks: {
          display: false
        },
        border: {
          display: false
        }
      },
      y: {
        grid: {
          color: 'rgba(255, 255, 255, 0.05)'
        },
        ticks: {
          color: '#9ca3af'
        },
        border: {
          display: false
        }
      }
    },
    plugins: {
      legend: {
        position: 'top',
        labels: {
          color: '#9ca3af',
          usePointStyle: true,
          boxWidth: 6,
          font: {
            size: 10
          }
        }
      },
      tooltip: {
        backgroundColor: 'rgba(26, 26, 26, 0.9)',
        titleColor: '#ffffff',
        bodyColor: '#e5e7eb',
        borderColor: 'rgba(139, 92, 246, 0.3)',
        borderWidth: 1,
        padding: 10
      }
    }
  };
  
  // Specific options for RSI
  const rsiOptions = {
    ...chartOptions,
    scales: {
      ...chartOptions.scales,
      y: {
        ...chartOptions.scales.y,
        min: 0,
        max: 100,
        grid: {
          color: 'rgba(255, 255, 255, 0.05)'
        }
      }
    }
  };
  
  return (
    <TechnicalContainer>
      <ChartCard>
        <h3>Relative Strength Index (RSI)</h3>
        {rsiChartData ? (
          <Line data={rsiChartData} options={rsiOptions} />
        ) : (
          <NoData>No RSI data available</NoData>
        )}
      </ChartCard>
      
      <ChartCard>
        <h3>MACD</h3>
        {macdChartData ? (
          <Line data={macdChartData} options={chartOptions} />
        ) : (
          <NoData>No MACD data available</NoData>
        )}
      </ChartCard>
      
      <ChartCard className="full-width">
        <h3>Bollinger Bands</h3>
        {bbChartData ? (
          <Line data={bbChartData} options={chartOptions} />
        ) : (
          <NoData>No Bollinger Bands data available</NoData>
        )}
      </ChartCard>
    </TechnicalContainer>
  );
};

export default TechnicalIndicators;
