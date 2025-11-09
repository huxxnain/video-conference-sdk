package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var (
	rooms      = make(map[string][]*websocket.Conn)
	roomsMutex sync.Mutex
)

// WebRTC signaling via WebSocket
func SignalingHandler(c *gin.Context) {
	r := c.Query("room")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	roomsMutex.Lock()
	rooms[r] = append(rooms[r], conn)
	roomsMutex.Unlock()

	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// Broadcast message to other connections in the same room
		roomsMutex.Lock()
		for _, peer := range rooms[r] {
			if peer != conn {
				peer.WriteMessage(mt, msg)
			}
		}
		roomsMutex.Unlock()
	}
}