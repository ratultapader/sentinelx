package stream

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	"sentinelx/models"
)

// Store connected clients
var clients = make(map[*websocket.Conn]bool)

// Mutex for thread safety
var mu sync.Mutex

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Handle incoming WebSocket connections
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("DEBUG: WebSocket upgrade failed:", err)
		return
	}

	fmt.Println("DEBUG: new WebSocket client connected")

	mu.Lock()
	clients[ws] = true
	mu.Unlock()

	// ✅ SEND INITIAL MESSAGE (VERY IMPORTANT)
	ws.WriteMessage(websocket.TextMessage, []byte(`{"type":"connected"}`))

	// 🔥 KEEP CONNECTION ALIVE
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("DEBUG: client disconnected")

			mu.Lock()
			delete(clients, ws)
			mu.Unlock()

			ws.Close()
			break
		}
	}
}

// BroadcastAlert sends the full alert object to all connected clients.
func BroadcastAlert(alert models.Alert) {
	message, err := json.Marshal(alert)
	if err != nil {
		fmt.Println("DEBUG: failed to marshal alert:", err)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	fmt.Println("DEBUG: broadcasting to", len(clients), "clients")

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			fmt.Println("DEBUG: removing disconnected client")
			client.Close()
			delete(clients, client)
		}
	}
}