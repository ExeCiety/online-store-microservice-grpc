package models

import "time"

type User struct {
	ID           string    `gorm:"type:uuid;primaryKey"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string    `gorm:"type:varchar(255);not null"`
	Name         string    `gorm:"type:varchar(255);not null"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
}

func (User) TableName() string {
	return "users"
}
