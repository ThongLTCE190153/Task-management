package handlers

import (
	"net/http"

	appwebsocket "trithong.com/task-golang/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	hub *appwebsocket.Hub
}

func NewWebSocketHandler(hub *appwebsocket.Hub) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *WebSocketHandler) HandleConnection(c *gin.Context) {
	// Lấy userID từ middleware đã xác thực
	userID := c.GetInt("user_id")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	// Truyền userID vào Hub
	h.hub.AddClient(userID, conn)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			h.hub.RemoveClient(userID, conn)
			break
		}
	}
}