package model

import "time"

// Role represents the authorization level of a user.
type Role string

const (
	// RoleAdmin has full access to all features including management.
	RoleAdmin Role = "admin"
	// RoleUser can use the LLM proxy and view the dashboard.
	RoleUser Role = "user"
	// RoleViewer can only view the dashboard and audit logs.
	RoleViewer Role = "viewer"
)

// User represents a gateway user account stored in the database.
type User struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	Username     string    `gorm:"uniqueIndex;not null;size:64" json:"username"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Role         Role      `gorm:"not null;default:user;size:16" json:"role"`
	Email        string    `gorm:"size:128" json:"email"`
	APIKeyHash   string    `gorm:"size:128" json:"-"`
	APIKeyPrefix string    `gorm:"size:12" json:"api_key_prefix"`
	RateLimit    int       `gorm:"default:60" json:"rate_limit"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
