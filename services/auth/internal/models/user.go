package models

import "gorm.io/gorm"

type User struct {
	gorm.Model        // Adds ID, CreatedAt, UpdatedAt, DeletedAt automatically
	Email      string `gorm:"uniqueIndex;not null"`
	Password   string `gorm:"not null"` // We will store the HASH, not plain text
}
