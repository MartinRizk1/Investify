import axios from 'axios';

/**
 * Stock Service - Handles API calls and WebSocket connections
 */
class StockService {
  constructor() {
    this.websocket = null;
    this.reconnectAttempts = 0;
    this.reconnectInterval = null;
    this.listeners = [];
    this.isConnected = false;
  }

  /**
   * Connect to the WebSocket for real-time updates
   * @param {string} ticker - The stock symbol
   */
  connectWebSocket(ticker) {
    if (!ticker) return;
    
    // Close existing connection if any
    this.disconnectWebSocket();
    
    try {
      // Use secure WebSocket if on HTTPS
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      const wsUrl = `${protocol}//${window.location.host}/ws/stocks/${ticker}`;
      
      this.websocket = new WebSocket(wsUrl);
      
      this.websocket.onopen = () => {
        console.log('WebSocket connected');
        this.isConnected = true;
        this.reconnectAttempts = 0;
        this.notifyListeners({ type: 'connection', status: true });
      };
      
      this.websocket.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          this.notifyListeners({ type: 'data', data });
        } catch (error) {
          console.error('Error parsing WebSocket data:', error);
        }
      };
      
      this.websocket.onclose = () => {
        console.log('WebSocket disconnected');
        this.isConnected = false;
        this.notifyListeners({ type: 'connection', status: false });
        
        // Try to reconnect with exponential backoff
        if (this.reconnectAttempts < 5) {
          const backoffTime = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000);
          this.reconnectInterval = setTimeout(() => {
            this.reconnectAttempts++;
            console.log(`Attempting to reconnect (${this.reconnectAttempts})...`);
            this.connectWebSocket(ticker);
          }, backoffTime);
        }
      };
      
      this.websocket.onerror = (error) => {
        console.error('WebSocket error:', error);
        this.isConnected = false;
        this.notifyListeners({ type: 'error', error });
      };
      
      return true;
    } catch (error) {
      console.error('Error establishing WebSocket connection:', error);
      return false;
    }
  }

  /**
   * Disconnect from the WebSocket
   */
  disconnectWebSocket() {
    if (this.websocket) {
      this.websocket.close();
      this.websocket = null;
    }
    
    if (this.reconnectInterval) {
      clearTimeout(this.reconnectInterval);
      this.reconnectInterval = null;
    }
    
    this.isConnected = false;
  }

  /**
   * Add a listener for WebSocket events
   * @param {Function} listener - The callback function to be called when events occur
   */
  addListener(listener) {
    if (typeof listener === 'function' && !this.listeners.includes(listener)) {
      this.listeners.push(listener);
    }
  }

  /**
   * Remove a listener
   * @param {Function} listener - The listener to remove
   */
  removeListener(listener) {
    const index = this.listeners.indexOf(listener);
    if (index > -1) {
      this.listeners.splice(index, 1);
    }
  }

  /**
   * Notify all listeners of an event
   * @param {Object} event - The event object
   */
  notifyListeners(event) {
    this.listeners.forEach(listener => {
      try {
        listener(event);
      } catch (error) {
        console.error('Error in WebSocket listener:', error);
      }
    });
  }

  /**
   * Fetch stock data from the API as a fallback
   * @param {string} ticker - The stock symbol
   * @returns {Promise} - A promise that resolves to the stock data
   */
  async fetchStockData(ticker) {
    try {
      const response = await axios.get(`/api/stocks/${ticker}`);
      return response.data;
    } catch (error) {
      console.error('Error fetching stock data:', error);
      throw error;
    }
  }
}

// Create and export a singleton instance
const stockService = new StockService();
export default stockService;
