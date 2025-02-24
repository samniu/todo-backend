package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"todo-backend/internal/model"
	"todo-backend/pkg/db"
	"todo-backend/pkg/ws"

	"github.com/gin-gonic/gin"
)

// CreateTodo 创建待办事项
func CreateTodo(c *gin.Context) {
	userID, _ := c.Get("userID")

	var input model.TodoCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 确保 DueDate 是有效的 RFC3339 格式
	if input.DueDate.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due date"})
		return
	}

	todo := model.Todo{
		UserID:      userID.(uint),
		Title:       input.Title,
		Description: input.Description,
		DueDate:     input.DueDate, // 直接存储指针，支持 null
		RepeatType:  input.RepeatType,
		Note:        input.Note,
		IsCompleted: false,
		IsFavorite:  false,
	}

	if err := db.DB.Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}

	// 发送 WebSocket 通知
	notification, _ := json.Marshal(map[string]interface{}{
		"type": "todo_created",
		"data": todo,
	})
	if clients, ok := ws.Manager.Clients[userID.(uint)]; ok {
		for client := range clients {
			client.Send <- notification
		}
	}

	c.JSON(http.StatusCreated, todo)
}

// UpdateTodo 更新待办事项
func UpdateTodo(c *gin.Context) {
	userID, _ := c.Get("userID")
	todoID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		log.Println("Invalid todo ID")
		return
	}

	var todo model.Todo
	if err := db.DB.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		log.Println("Todo not found")
		return
	}

	var input model.TodoCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("Error binding JSON: %v", err) // 打印错误信息
		log.Printf("Received input: %+v", input)  // 打印输入数据
		return
	}

	todo.Title = input.Title
	todo.Description = input.Description
	todo.DueDate = input.DueDate
	todo.RepeatType = input.RepeatType
	todo.Note = input.Note

	if err := db.DB.Save(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
		return
	}

	// 发送 WebSocket 通知
	notification, _ := json.Marshal(map[string]interface{}{
		"type": "todo_updated",
		"data": todo,
	})
	if clients, ok := ws.Manager.Clients[userID.(uint)]; ok {
		for client := range clients {
			client.Send <- notification
		}
	}

	c.JSON(http.StatusOK, todo)
}

// DeleteTodo 删除待办事项
func DeleteTodo(c *gin.Context) {
	userID, _ := c.Get("userID")
	todoID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}

	var todo model.Todo
	if err := db.DB.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	if err := db.DB.Delete(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
		return
	}

	// 发送 WebSocket 通知
	notification, _ := json.Marshal(map[string]interface{}{
		"type": "todo_deleted",
		"data": map[string]interface{}{
			"id":      todo.ID,
			"user_id": todo.UserID,
		},
	})
	if clients, ok := ws.Manager.Clients[userID.(uint)]; ok {
		for client := range clients {
			client.Send <- notification
		}
	}

	// 返回被删除的 todo 信息
	c.JSON(http.StatusOK, gin.H{
		"message": "Todo deleted successfully",
		"todo": gin.H{
			"id":         todo.ID,
			"user_id":    todo.UserID,
			"title":      todo.Title,
			"created_at": todo.CreatedAt,
			"updated_at": todo.UpdatedAt,
		},
	})
}

// ToggleTodo 切换待办事项的完成状态
func ToggleTodo(c *gin.Context) {
	userID, _ := c.Get("userID")
	todoID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}

	var todo model.Todo
	if err := db.DB.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	todo.IsCompleted = !todo.IsCompleted // 修改字段名

	if err := db.DB.Save(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
		return
	}

	// 发送 WebSocket 通知
	notification, _ := json.Marshal(map[string]interface{}{
		"type": "todo_toggled",
		"data": todo,
	})
	if clients, ok := ws.Manager.Clients[userID.(uint)]; ok {
		for client := range clients {
			client.Send <- notification
		}
	}

	c.JSON(http.StatusOK, todo)
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
