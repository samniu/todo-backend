package model

import (
	"time"
)

type Todo struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	UserID      uint       `json:"user_id"`
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`     // 允许为空
	IsCompleted bool       `json:"is_completed"` // 修改字段名
	IsFavorite  bool       `json:"is_favorite"`  // 添加字段
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	RepeatType  string     `json:"repeat_type"`
	Note        string     `json:"note"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

type TodoCreate struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"` // 改成指针类型
	RepeatType  string     `json:"repeat_type"`
	Note        string     `json:"note"`
}
