package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"todo-backend/internal/api/handler"
	"todo-backend/internal/api/middleware"
	"todo-backend/pkg/db"
	"todo-backend/pkg/ws"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// 初始化数据库
	db.InitDB()

	// 初始化 WebSocket 管理器
	ws.InitManager()

	r := gin.Default()

	// 公开路由
	r.POST("/api/register", handler.Register)
	r.POST("/api/login", handler.Login)

	// 需要认证的路由
	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		// Todo 路由
		auth.POST("/todos", handler.CreateTodo)
		auth.GET("/todos", handler.GetTodos)
		auth.PUT("/todos/:id", handler.UpdateTodo)
		auth.DELETE("/todos/:id", handler.DeleteTodo)
		auth.PATCH("/todos/:id/toggle", handler.ToggleTodo)

		// WebSocket 路由
		auth.GET("/ws", handler.HandleWebSocket)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
