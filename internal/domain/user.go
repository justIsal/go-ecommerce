package domain

import (
	"time"
	"gorm.io/gorm"
)

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Email     string         `gorm:"uniqueIndex;type:varchar(100);not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"` 
	Role      string         `gorm:"type:varchar(20);default:'user'" json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` 
}

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index;not null"`
	Token     string    `gorm:"index;not null" json:"token"` 
	ExpiresAt time.Time `json:"expires_at"`
	Revoked   bool      `gorm:"default:false"`
	CreatedAt time.Time
}