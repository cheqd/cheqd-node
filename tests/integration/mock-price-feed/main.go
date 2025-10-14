package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"

	"github.com/ojo-network/price-feeder/oracle/types"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type PriceData struct {
	Price  string // string format to match provider
	Volume string
}

type MockPriceFeed struct {
	clients   map[*websocket.Conn]bool
	broadcast chan interface{}
	mutex     sync.RWMutex
	prices    map[string]PriceData
}

func NewMockPriceFeed() *MockPriceFeed {
	return &MockPriceFeed{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan interface{}, 100),
		prices: map[string]PriceData{
			"CHEQUSDT": {Price: "0.05", Volume: "1000000"},
			"USDCUSDT": {Price: "1.0", Volume: "5000000"},
		},
	}
}

// --- WebSocket handler ---
func (m *MockPriceFeed) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	m.mutex.Lock()
	m.clients[conn] = true
	clientCount := len(m.clients)
	m.mutex.Unlock()
	log.Printf("New WebSocket client connected. Total clients: %d", clientCount)

	// Send initial prices
	m.mutex.RLock()
	for symbol, data := range m.prices {
		m.sendTicker(conn, symbol, data)
		m.sendCandle(conn, symbol, data)
	}
	m.mutex.RUnlock()

	// Keep connection alive
	for {
		var msg map[string]interface{}
		if err := conn.ReadJSON(&msg); err != nil {
			log.Printf("Client disconnected: %v", err)
			m.mutex.Lock()
			delete(m.clients, conn)
			m.mutex.Unlock()
			break
		}

		if method, ok := msg["method"].(string); ok && method == "SUBSCRIPTION" {
			log.Printf("Received subscription: %v", msg["params"])
		}
	}
}

// --- Helper to send Protobuf ticker ---
func (m *MockPriceFeed) sendTicker(conn *websocket.Conn, symbol string, data PriceData) {
	ticker := &types.BookTicker{
		Symbol: &symbol,
		PublicBookTicker: &types.PublicBookTickerV3Api{
			BidPrice:    data.Price,
			BidQuantity: data.Volume,
		},
	}
	bz, _ := proto.Marshal(ticker)
	conn.WriteMessage(websocket.BinaryMessage, bz)
}

// --- Helper to send Protobuf candle ---
func (m *MockPriceFeed) sendCandle(conn *websocket.Conn, symbol string, data PriceData) {
	now := time.Now().Unix()
	closePrice := data.Price
	volume := data.Volume

	candle := &types.SpotKline{
		Symbol: &symbol,
		PublicSpotKline: &types.PublicSpotKlineV3Api{
			ClosingPrice: closePrice,
			Volume:       volume,
			WindowEnd:    now,
		},
	}
	bz, _ := proto.Marshal(candle)
	conn.WriteMessage(websocket.BinaryMessage, bz)
}

// --- Broadcast updated prices to all clients ---
func (m *MockPriceFeed) broadcastPrices(ctx context.Context) {
	for {
		select {
		case update := <-m.broadcast:
			m.mutex.RLock()
			for client := range m.clients {
				switch msg := update.(type) {
				case *types.BookTicker:
					bz, _ := proto.Marshal(msg)
					client.WriteMessage(websocket.BinaryMessage, bz)
				case *types.SpotKline:
					bz, _ := proto.Marshal(msg)
					client.WriteMessage(websocket.BinaryMessage, bz)
				}
			}
			m.mutex.RUnlock()
		case <-ctx.Done():
			return
		}
	}
}

// --- Periodically update prices ---
func (m *MockPriceFeed) updatePrices(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	iteration := 0

	for {
		select {
		case <-ticker.C:
			iteration++
			m.mutex.Lock()
			for symbol, data := range m.prices {
				// simple variation
				newPrice := fmt.Sprintf("%.6f", parseFloat(data.Price)+float64(iteration%10)*0.0001)
				newVolume := fmt.Sprintf("%.2f", parseFloat(data.Volume)*(1+float64(iteration%10)*0.0001))
				m.prices[symbol] = PriceData{Price: newPrice, Volume: newVolume}

				// broadcast updates
				tickerUpdate := &types.BookTicker{
					Symbol: &symbol,
					PublicBookTicker: &types.PublicBookTickerV3Api{
						BidPrice:    newPrice,
						BidQuantity: newVolume,
					},
				}
				candleUpdate := &types.SpotKline{
					Symbol: &symbol,
					PublicSpotKline: &types.PublicSpotKlineV3Api{
						ClosingPrice: newPrice,
						Volume:       newVolume,
						WindowEnd:    time.Now().Unix(),
					},
				}
				m.broadcast <- tickerUpdate
				m.broadcast <- candleUpdate
			}
			m.mutex.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

// --- Helper ---
func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

// --- HTTP endpoints ---
func (m *MockPriceFeed) setPriceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
		Volume string `json:"volume,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	m.mutex.Lock()
	if req.Volume == "" {
		req.Volume = "1000000"
	}
	m.prices[req.Symbol] = PriceData{Price: req.Price, Volume: req.Volume}
	m.mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (m *MockPriceFeed) getPricesHandler(w http.ResponseWriter, r *http.Request) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	simplePrices := make(map[string]string)
	for s, d := range m.prices {
		simplePrices[s] = d.Price
	}
	json.NewEncoder(w).Encode(simplePrices)
}

func main() {
	port := ":8080"
	feed := NewMockPriceFeed()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go feed.broadcastPrices(ctx)
	go feed.updatePrices(ctx)

	http.HandleFunc("/ws", feed.handleWS)
	http.HandleFunc("/prices", feed.getPricesHandler)
	http.HandleFunc("/set-price", feed.setPriceHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	http.HandleFunc("/api/v3/ticker/price", func(w http.ResponseWriter, r *http.Request) {
		feed.mutex.RLock()
		defer feed.mutex.RUnlock()
		var resp []map[string]string
		for symbol := range feed.prices {
			resp = append(resp, map[string]string{"symbol": symbol})
		}
		json.NewEncoder(w).Encode(resp)
	})

	log.Printf("Mock MEXC Protobuf Price Feed server starting on %s", port)
	log.Printf("WebSocket endpoint: ws://localhost%s/ws", port)
	log.Printf("HTTP GET prices: http://localhost%s/prices", port)
	log.Printf("HTTP POST set price: http://localhost%s/set-price", port)

	log.Fatal(http.ListenAndServe(port, nil))
}
