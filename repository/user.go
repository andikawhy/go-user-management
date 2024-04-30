package repository

import (
	"time"
)

type User struct {
	ID        uint64 `json:"id" gorm:"primary_key"`
	Username  string `json:"username" gorm:"unique"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RemoveUser struct {
	ID uint64 `json:"username" binding:"required"`
}
