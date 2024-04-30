package repository

import (
	"time"
)

type User struct {
	ID        uint64    `json:"id" gorm:"primary_key"`
	Username  string    `json:"username" gorm:"unique"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`
}

type UserResponse struct {
	ID        uint64    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdat"`
}
