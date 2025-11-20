package domain

import "time"

type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RoleClient UserRole = "client"
)

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         UserRole  `json:"role"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}
