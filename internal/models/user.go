package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        string    `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	Name      string    `db:"name" json:"name"`
	Password  string    `db:"password" json:"-"`
	Role      string    `db:"role" json:"role"`
	Active    bool      `db:"active" json:"active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// CreateUserRequest is the request to create a user
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role" binding:"required"`
}

// UpdateUserRequest is the request to update a user
type UpdateUserRequest struct {
	Name   string `json:"name"`
	Email  string `json:"email" binding:"email"`
	Role   string `json:"role"`
	Active bool   `json:"active"`
}