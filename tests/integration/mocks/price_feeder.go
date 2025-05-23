package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"
)

// PriceData represents the price data structure returned by exchanges
type PriceData struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
	Symbol string `json:"symbol"`
	Time   int64  `json:"timestamp"`
}

// ExchangeMock represents a mock exchange API server
type ExchangeMock struct {
	Server *httptest.Server
	Prices map[string]PriceData
}

// NewExchangeMock creates a new mock exchange API server
func NewExchangeMock() *ExchangeMock {
	mock := &ExchangeMock{
		Prices: make(map[string]PriceData),
	}

	// Initialize with some default prices
	mock.SetPrice("CHEQ", "0.123", "1000000")
	mock.SetPrice("ATOM", "10.45", "5000000")
	mock.SetPrice("USDT", "1.001", "10000000")
	mock.SetPrice("BTC", "45000.75", "1000")
	mock.SetPrice("ETH", "2500.50", "3000")

	// Create a mock HTTP server
	mock.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse the request to determine which price data to return
		symbol := r.URL.Query().Get("symbol")
		if symbol == "" {
			// If no symbol is provided, return all prices
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mock.Prices)
			return
		}

		// If a specific symbol is requested
		price, exists := mock.Prices[symbol]
		if !exists {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, `{"error": "Symbol %s not found"}`, symbol)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(price)
	}))

	return mock
}

// SetPrice sets a price for a symbol in the mock
func (m *ExchangeMock) SetPrice(symbol, price, volume string) {
	m.Prices[symbol] = PriceData{
		Price:  price,
		Volume: volume,
		Symbol: symbol,
		Time:   time.Now().Unix(),
	}
}

// Close shuts down the mock server
func (m *ExchangeMock) Close() {
	if m.Server != nil {
		m.Server.Close()
	}
}

// GetURL returns the URL of the mock server
func (m *ExchangeMock) GetURL() string {
	return m.Server.URL
}

// MEXCMock specifically mocks the MEXC API
type MEXCMock struct {
	Server *httptest.Server
	Prices map[string]MEXCResponse
}

// MEXCResponse represents the response structure from MEXC
type MEXCResponse struct {
	Code int `json:"code"`
	Data []struct {
		Symbol    string `json:"symbol"`
		Price     string `json:"last"`
		Volume    string `json:"volume_24h"`
		Bid       string `json:"bid"`
		Ask       string `json:"ask"`
		Timestamp int64  `json:"timestamp"`
	} `json:"data"`
	Message string `json:"msg"`
}

// NewMEXCMock creates a new MEXC API mock
func NewMEXCMock() *MEXCMock {
	mock := &MEXCMock{
		Prices: make(map[string]MEXCResponse),
	}

	// Initialize with default responses
	mock.SetPrice("CHEQ_USDT", "0.123", "1000000")
	mock.SetPrice("ATOM_USDT", "10.45", "5000000")

	// Create a mock HTTP server
	mock.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// MEXC endpoints usually have different formats for different operations
		// Here we're mocking the ticker endpoint
		endpoint := r.URL.Path

		// Handle the ticker endpoint (example: /api/v3/ticker/24h)
		if endpoint == "/api/v3/ticker/24h" {
			// Get symbol from query params
			symbol := r.URL.Query().Get("symbol")

			if symbol == "" {
				// If no symbol provided, return all prices
				allData := []struct {
					Symbol    string `json:"symbol"`
					Price     string `json:"last"`
					Volume    string `json:"volume_24h"`
					Bid       string `json:"bid"`
					Ask       string `json:"ask"`
					Timestamp int64  `json:"timestamp"`
				}{}

				// Combine all price data
				for symbol, response := range mock.Prices {
					fmt.Println("symbol.....", symbol)
					allData = append(allData, response.Data...)
				}

				// Prepare combined response
				combinedResponse := MEXCResponse{
					Code:    200,
					Data:    allData,
					Message: "success",
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(combinedResponse)
				return
			}

			// Handle specific symbol request
			response, exists := mock.Prices[symbol]
			if !exists {
				// Return empty but successful response for unknown symbols
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(MEXCResponse{
					Code: 200,
					Data: []struct {
						Symbol    string `json:"symbol"`
						Price     string `json:"last"`
						Volume    string `json:"volume_24h"`
						Bid       string `json:"bid"`
						Ask       string `json:"ask"`
						Timestamp int64  `json:"timestamp"`
					}{},
					Message: "success",
				})
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}

		// Handle other endpoints with 404
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"code": 404, "data": null, "msg": "Endpoint not found"}`)
	}))

	return mock
}

// SetPrice sets a price for a symbol in the MEXC mock
func (m *MEXCMock) SetPrice(symbol, price, volume string) {
	// Create response structure
	response := MEXCResponse{
		Code: 200,
		Data: []struct {
			Symbol    string `json:"symbol"`
			Price     string `json:"last"`
			Volume    string `json:"volume_24h"`
			Bid       string `json:"bid"`
			Ask       string `json:"ask"`
			Timestamp int64  `json:"timestamp"`
		}{
			{
				Symbol:    symbol,
				Price:     price,
				Volume:    volume,
				Bid:       price, // Simplification, could be slightly lower in real API
				Ask:       price, // Simplification, could be slightly higher in real API
				Timestamp: time.Now().Unix(),
			},
		},
		Message: "success",
	}

	m.Prices[symbol] = response
}

// Close shuts down the mock server
func (m *MEXCMock) Close() {
	if m.Server != nil {
		m.Server.Close()
	}
}

// GetURL returns the URL of the mock server
func (m *MEXCMock) GetURL() string {
	return m.Server.URL
}
