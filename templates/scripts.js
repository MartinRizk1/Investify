// Add loading overlay for better UX
function showLoadingOverlay() {
    document.getElementById('loading-overlay').style.display = 'flex';
    document.getElementById('search-btn').classList.add('loading');
    document.getElementById('search-btn-text').innerText = 'Analyzing...';
}

function hideLoadingOverlay() {
    document.getElementById('loading-overlay').style.display = 'none';
    document.getElementById('search-btn').classList.remove('loading');
    document.getElementById('search-btn-text').innerText = 'Search';
}

// Initialize stock chart with data
function initializeCharts() {
    // Only initialize chart if the element and stock data exist
    if (document.getElementById('price-chart') && window.stockData) {
        const ctx = document.getElementById('price-chart').getContext('2d');
        
        // Generate some simulated historical price data
        const basePrice = window.stockData.price;
        const volatility = Math.abs(window.stockData.change) / basePrice * 5; 
        const trend = window.stockData.change > 0 ? 1 : -1;
        const days = 30;
        
        const labels = [];
        const prices = [];
        
        // Generate dates and prices for the past 30 days
        for (let i = days; i >= 0; i--) {
            const date = new Date();
            date.setDate(date.getDate() - i);
            labels.push(date.toISOString().split('T')[0]);
            
            // Simulated price data with trend and randomness
            if (i === 0) {
                // Today's actual price
                prices.push(basePrice);
            } else {
                const randomFactor = (Math.random() - 0.5) * volatility;
                const dayPrice = basePrice - (trend * (i/5) + randomFactor);
                prices.push(Math.max(dayPrice, 0.1 * basePrice).toFixed(2)); // Ensure no negative prices
            }
        }
        
        // Create gradient for chart
        const gradient = ctx.createLinearGradient(0, 0, 0, 300);
        if (trend > 0) {
            gradient.addColorStop(0, 'rgba(16, 185, 129, 0.7)');  
            gradient.addColorStop(1, 'rgba(16, 185, 129, 0)');
        } else {
            gradient.addColorStop(0, 'rgba(239, 68, 68, 0.7)');
            gradient.addColorStop(1, 'rgba(239, 68, 68, 0)');
        }
        
        new Chart(ctx, {
            type: 'line',
            data: {
                labels: labels,
                datasets: [{
                    label: window.stockData.ticker + ' Price',
                    data: prices,
                    borderColor: trend > 0 ? '#10b981' : '#ef4444',
                    backgroundColor: gradient,
                    borderWidth: 2,
                    pointRadius: 1,
                    pointHoverRadius: 5,
                    fill: true,
                    tension: 0.2
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        display: false
                    },
                    tooltip: {
                        mode: 'index',
                        intersect: false,
                        displayColors: false,
                        callbacks: {
                            label: function(context) {
                                return `$${context.parsed.y}`;
                            }
                        }
                    },
                    annotation: {
                        annotations: {
                            prediction: {
                                type: 'line',
                                scaleID: 'y',
                                value: window.stockData.predicted_price,
                                borderColor: '#8b5cf6',
                                borderWidth: 2,
                                borderDash: [5, 5],
                                label: {
                                    backgroundColor: '#8b5cf6',
                                    content: 'AI Prediction: $' + window.stockData.predicted_price.toFixed(2),
                                    enabled: true,
                                    position: 'start'
                                }
                            }
                        }
                    }
                },
                scales: {
                    x: {
                        grid: {
                            display: false,
                            drawBorder: false
                        },
                        ticks: {
                            color: '#9ca3af',
                            font: {
                                size: 10
                            },
                            maxRotation: 0,
                            callback: function(value, index, values) {
                                // Show fewer x-axis labels for better readability
                                return index % 3 === 0 ? this.getLabelForValue(value) : '';
                            }
                        }
                    },
                    y: {
                        grid: {
                            color: 'rgba(255, 255, 255, 0.05)',
                            drawBorder: false
                        },
                        ticks: {
                            color: '#9ca3af',
                            font: {
                                size: 10
                            },
                            callback: function(value) {
                                return '$' + value;
                            }
                        }
                    }
                }
            }
        });
    }
    
    // Initialize confidence meter if element exists
    if (document.getElementById('confidence-meter') && window.stockData) {
        const confidence = window.stockData.prediction_confidence || 70;
        const confidenceMeter = document.getElementById('confidence-meter');
        const needle = document.getElementById('confidence-needle');
        
        // Set the needle position based on confidence
        const rotation = (confidence / 100) * 180 - 90;
        needle.style.transform = `rotate(${rotation}deg)`;
        
        // Set the color based on confidence level
        let color;
        if (confidence >= 80) {
            color = 'var(--success)';
        } else if (confidence >= 60) {
            color = 'var(--purple-primary)';
        } else {
            color = 'var(--danger)';
        }
        
        needle.style.backgroundColor = color;
        document.getElementById('confidence-value').innerText = Math.round(confidence) + '%';
        document.getElementById('confidence-value').style.color = color;
    }
    
    // Initialize and animate the factors list
    if (document.querySelectorAll('.factor-item').length > 0) {
        gsap.from('.factor-item', {
            duration: 0.5,
            y: 20,
            opacity: 0,
            stagger: 0.1,
            ease: 'power2.out'
        });
    }
    
    // Animate the price and prediction if they exist
    if (document.querySelector('.price-big')) {
        gsap.from('.price-big', {
            duration: 0.7,
            scale: 0.9,
            opacity: 0,
            ease: 'back.out'
        });
        
        gsap.from('.predicted-price', {
            duration: 0.7,
            delay: 0.2,
            scale: 0.9,
            opacity: 0,
            ease: 'back.out'
        });
    }
}

// Handle form submission with loading overlay
document.addEventListener('DOMContentLoaded', function() {
    const searchForm = document.getElementById('search-form');
    if (searchForm) {
        searchForm.addEventListener('submit', function() {
            showLoadingOverlay();
            sessionStorage.setItem('loading', 'true');
        });
    }
    
    // Check if we need to show the results animation
    if (window.stockData) {
        if (sessionStorage.getItem('loading') === 'true') {
            // Simulate a short delay for smoother UX
            setTimeout(() => {
                hideLoadingOverlay();
                initializeCharts();
                sessionStorage.removeItem('loading');
            }, 300);
        } else {
            initializeCharts();
        }
    }
    
    // Add the error shake animation if there's an error
    const errorMessage = document.querySelector('.error-message');
    if (errorMessage && errorMessage.innerText.trim() !== '') {
        gsap.from(errorMessage, {
            duration: 0.3,
            x: 10,
            ease: 'power2.inOut',
            repeat: 2,
            yoyo: true
        });
    }
});
