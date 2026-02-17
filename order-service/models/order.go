package models

import "time"

type Order struct {
	ID          string    `gorm:"type:uuid;primaryKey"`
	UserID      string    `gorm:"type:uuid;not null;index"`
	ProductName string    `gorm:"type:varchar(255);not null"`
	Quantity    int32     `gorm:"not null"`
	TotalPrice  float64   `gorm:"type:numeric(12,2);not null"`
	Status      string    `gorm:"type:varchar(50);not null;default:pending"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

func (Order) TableName() string {
	return "orders"
}
