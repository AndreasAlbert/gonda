package storage

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Provider string `gorm:"primaryKey; not null"`
	Name     string `gorm:"primaryKey; not null"`
}

// type ApiKey struct {
// 	gorm.Model

// 	User User
// }
