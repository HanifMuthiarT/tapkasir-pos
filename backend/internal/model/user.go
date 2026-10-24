package model

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Role      string    `json:"role" gorm:"type:enum('owner', 'manager', 'cashier');not null"`
	PIN       string    `json:"pin" gorm:"type:varchar(255)"` // Akan disimpan dalam bentuk hash bcrypt
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}