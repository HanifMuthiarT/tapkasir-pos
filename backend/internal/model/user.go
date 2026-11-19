package model

import (
	"errors"
	"time"
)

var (
	ErrInvalidPIN      = errors.New("PIN yang dimasukkan salah")
	ErrUserNotFound    = errors.New("user/role tidak ditemukan")
	ErrUnauthorized    = errors.New("akses tidak diizinkan")
	ErrPINAlreadyExist = errors.New("user dengan role ini sudah terdaftar")
)

// 1. Struct Entity User
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Role      string    `json:"role" gorm:"type:enum('owner', 'manager', 'cashier');not null"`
	PIN       string    `json:"pin" gorm:"type:varchar(255)"` // Akan disimpan dalam bentuk hash bcrypt
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "user"
}

// 2. Data Transfer Objects (DTO) untuk Request & Response
type VerifyPINRequest struct {
	Role string `json:"role" binding:"required,oneof=owner manager cashier"`
	PIN  string `json:"pin" binding:"required,min=4,max=6"`
}

type VerifyPINResponse struct {
	Token     string `json:"token"`
	Role      string `json:"role"`
	ExpiresAt int64  `json:"expires_at"`
}
type ChangePINRequest struct {
	Role   string `json:"role" binding:"required,oneof=owner manager cashier"`
	OldPIN string `json:"old_pin" binding:"required,min=4,max=6"`
	NewPIN string `json:"new_pin" binding:"required,min=4,max=6"`
}
// 3. Kontrak untuk Repository (Database)
type UserRepository interface {
	GetByRole(role string) (*User, error)
	Create(user *User) error
	UpdatePIN(role string, hashedPIN string) error
	Count() (int64, error)
}
// 4. Kontrak untuk Usecase (Logika Bisnis)
type UserUsecase interface {
	VerifyPIN(req *VerifyPINRequest) (*VerifyPINResponse, error)
	ChangePIN(req *ChangePINRequest) error
	SeedDefaultUsers() error
}