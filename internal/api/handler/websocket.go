package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"todo-backend/pkg/ws"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Claims 结构用于存储 JWT 的声明
type Claims struct {
	UserID               uint `json:"user_id"` // 自定义字段：用户 ID
	jwt.RegisteredClaims      // 标准字段：exp, iat 等
}

// 定义用于签名的密钥
var jwtKey = []byte("your_secret_key") // 替换为你的实际密钥

func HandleWebSocket(c *gin.Context) {
	// 升级为 WebSocket 连接
	log.Println("Handling WebSocket connection")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}

	// 读取第一条消息（认证消息）
	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Printf("Failed to read auth message: %+v", err)
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Failed to read auth message"))
		conn.Close()
		return
	}

	// 解析认证消息
	var authMsg struct {
		Type  string `json:"type"`
		Token string `json:"token"`
	}
	if err := json.Unmarshal(message, &authMsg); err != nil {
		log.Printf("Invalid auth message: %+v", err)
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Invalid auth message"))
		conn.Close()
		return
	}

	// 检查消息类型
	if authMsg.Type != "auth" {
		log.Printf("Invalid message type: %s", authMsg.Type)
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Invalid message type"))
		conn.Close()
		return
	}

	// 验证 Token
	userID, err := validateToken(authMsg.Token)
	if err != nil {
		log.Printf("Unauthorized: %+v", err)
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Unauthorized"))
		conn.Close()
		return
	}

	// 认证成功
	log.Printf("Authentication successful for userID: %d", userID)
	conn.WriteMessage(websocket.TextMessage, []byte(`{"type": "auth_success"}`))

	// 注册客户端
	client := &ws.Client{
		ID:     userID,
		Socket: conn,
		Send:   make(chan []byte, 256),
	}
	ws.Manager.Register <- client

	// 启动读写 goroutines
	go client.WritePump()
	go client.ReadPump()
}

// parseToken 解析 JWT Token 并返回 Claims
func parseToken(tokenString string) (*Claims, error) {
	// 解析 Token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 检查签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtKey, nil
	})

	// 检查解析错误
	if err != nil {
		return nil, err
	}

	// 验证 Claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}

// validateToken 验证 Token 并返回用户 ID
func validateToken(token string) (uint, error) {
	// 解析 Token
	claims, err := parseToken(token)
	if err != nil {
		return 0, err
	}

	// 检查 Token 是否过期
	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		return 0, errors.New("token expired")
	}

	// 返回用户 ID
	return claims.UserID, nil
}
