package handler

import (
	"log"
	"net/http"

	"todo-backend/pkg/ws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(c *gin.Context) {
	userID, _ := c.Get("userID")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}

	client := &ws.Client{
		ID:     userID.(uint),
		Socket: conn,
		Send:   make(chan []byte, 256),
	}

	ws.Manager.Register <- client

	// 启动读写 goroutines
	go client.WritePump()
	go client.ReadPump()
}
