package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

func main() {
    // 加载环境变量
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    // 初始化 Gin 路由
    r := gin.Default()

    // 设置路由
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })

    // 启动服务器
    if err := r.Run(":8080"); err != nil {
        log.Fatal(err)
    }
}
