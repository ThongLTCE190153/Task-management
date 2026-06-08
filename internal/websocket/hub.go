package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Hub struct {
	// Lưu theo userID thay vì connection trực tiếp
	// 1 user có thể mở nhiều tab → nhiều connection
	clients map[int][]*websocket.Conn
	mu      sync.RWMutex
}

type Message struct {
	Event     string      `json:"event"`
	Data      interface{} `json:"data"`
	Message   string      `json:"message"`
	Timestamp time.Time   `json:"timestamp"`
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[int][]*websocket.Conn),
	}
}

func (h *Hub) AddClient(userID int, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[userID] = append(h.clients[userID], conn)
	log.Printf("User %d connected, total connections: %d", userID, len(h.clients[userID]))
}

func (h *Hub) RemoveClient(userID int, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	conns := h.clients[userID]

	// Lọc bỏ connection bị disconnect
	var remaining []*websocket.Conn
	for _, c := range conns {
		if c != conn {
			remaining = append(remaining, c)
		}
	}

	if len(remaining) == 0 {
		delete(h.clients, userID)
	} else {
		h.clients[userID] = remaining
	}

	conn.Close()
	log.Printf("User %d disconnected", userID)
}

// Gửi đến 1 user cụ thể
func (h *Hub) SendToUser(userID int, event string, data interface{}, message string) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	conns, exists := h.clients[userID]
	if !exists {
		return
	}

	wsMessage := Message{
		Event:     event,
		Data:      data,
		Message:   message,
		Timestamp: time.Now(),
	}

	jsonData, err := json.Marshal(wsMessage)
	if err != nil {
		log.Println("Failed to marshal websocket message:", err)
		return
	}

	for _, conn := range conns {
		err := conn.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			log.Printf("Failed to send message to user %d: %v", userID, err)
		}
	}
}

// Gửi đến nhiều user cùng lúc (dùng cho notify cả project)
func (h *Hub) SendToUsers(userIDs []int, event string, data interface{}, message string) {
	for _, userID := range userIDs {
		h.SendToUser(userID, event, data, message)
	}
}

// Giữ lại BroadcastJSON để không phải sửa nhiều chỗ
// Nhưng chỉ dùng khi thực sự cần gửi cho tất cả
func (h *Hub) BroadcastJSON(event string, data interface{}, message string) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	wsMessage := Message{
		Event:     event,
		Data:      data,
		Message:   message,
		Timestamp: time.Now(),
	}

	jsonData, err := json.Marshal(wsMessage)
	if err != nil {
		log.Println("Failed to marshal websocket message:", err)
		return
	}

	for _, conns := range h.clients {
		for _, conn := range conns {
			conn.WriteMessage(websocket.TextMessage, jsonData)
		}
	}
}