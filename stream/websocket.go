package stream

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	"sentinelx/models"
)

// ===============================
// 🔥 GLOBAL CLIENT STORE
// ===============================
var clients = make(map[*websocket.Conn]bool)
var mu sync.Mutex

// ===============================
// 🔥 WEBSOCKET UPGRADER
// ===============================
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all (dev mode)
	},
}

// ===============================
// 🔥 HANDLE CONNECTIONS
// ===============================
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("❌ WebSocket upgrade failed:", err)
		return
	}

	fmt.Println("🔌 New WebSocket client connected")

	// add client
	mu.Lock()
	clients[ws] = true
	fmt.Println("👥 Total clients:", len(clients))
	mu.Unlock()

	// send initial message
	err = ws.WriteMessage(websocket.TextMessage, []byte(`{"type":"connected"}`))
	if err != nil {
		ws.Close()
		return
	}

	// keep connection alive
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("❌ Client disconnected")

			mu.Lock()
			delete(clients, ws)
			fmt.Println("👥 Total clients:", len(clients))
			mu.Unlock()

			ws.Close()
			break
		}
	}
}

// ===============================
// 🔥 BROADCAST ALERT (REAL-TIME)
// ===============================
func BroadcastAlert(alert models.Alert) {
	message, err := json.Marshal(alert)
	if err != nil {
		fmt.Println("❌ Failed to marshal alert:", err)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	fmt.Println("📡 Broadcasting alert to", len(clients), "clients")

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			fmt.Println("⚠️ Removing disconnected client")
			client.Close()
			delete(clients, client)
		}
	}
}