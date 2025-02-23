package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"todo-backend/internal/api/handler"
	"todo-backend/internal/api/middleware"
	"todo-backend/pkg/db"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db.InitDB()

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
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
