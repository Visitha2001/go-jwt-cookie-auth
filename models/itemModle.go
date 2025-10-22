package models

import (
	"time"

	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	ID        uint      `gorm:"primarykey;autoIncrement:true;not null" json:"id"`
	UserID    uint      `gorm:"not null;references:users(id);default:0" json:"user_id"`
	Name      string    `gorm:"not null" json:"name"`
	Price     float64   `gorm:"not null" json:"price"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}

func MigrateItems(db *gorm.DB) error {
	err := db.AutoMigrate(&Item{})
	if err != nil {
		return err
	}
	return nil
}
