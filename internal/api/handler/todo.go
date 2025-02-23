package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "todo-backend/internal/model"
    "todo-backend/pkg/db"
)

// CreateTodo 创建待办事项
func CreateTodo(c *gin.Context) {
    userID, _ := c.Get("userID")
    
    var input model.TodoCreate
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    todo := model.Todo{
        UserID:      userID.(uint),
        Title:       input.Title,
        Description: input.Description,
        DueDate:     input.DueDate,
        RepeatType:  input.RepeatType,
        Note:        input.Note,
    }

    if err := db.DB.Create(&todo).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
        return
    }

    c.JSON(http.StatusCreated, todo)
}

// GetTodos 获取当前用户的所有待办事项
func GetTodos(c *gin.Context) {
    userID, _ := c.Get("userID")
    
    var todos []model.Todo
    if err := db.DB.Where("user_id = ?", userID).Find(&todos).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos"})
        return
    }

    c.JSON(http.StatusOK, todos)
}
