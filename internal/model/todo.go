package model

import (
    "time"
)

type Todo struct {
    ID          uint      `gorm:"primarykey" json:"id"`
    UserID      uint      `json:"user_id"`
    Title       string    `json:"title" binding:"required"`
    Description string    `json:"description"`
    DueDate     time.Time `json:"due_date"`
    Completed   bool      `json:"completed" gorm:"default:false"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    RepeatType  string    `json:"repeat_type"`
    Note        string    `json:"note"`

    User        User      `gorm:"foreignKey:UserID" json:"-"`
}

type TodoCreate struct {
    Title       string    `json:"title" binding:"required"`
    Description string    `json:"description"`
    DueDate     time.Time `json:"due_date"`
    RepeatType  string    `json:"repeat_type"`
    Note        string    `json:"note"`
}
