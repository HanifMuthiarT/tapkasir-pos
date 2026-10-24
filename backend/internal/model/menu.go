package model

import "errors"

var ErrMenuNotFound = errors.New("menu tidak ditemukan")

// 1. Struct Entity
type Menu struct {
	ID                 uint    `json:"id" gorm:"primaryKey"`
	NamaMenu           string  `json:"nama_menu"`
	HargaJual          float64 `json:"harga_jual"`
	StatusKetersediaan string  `json:"status_ketersediaan"`
}

func (Menu) TableName() string {
	return "menu"
}

// 2.Kontrak untuk Usecase (Logika Bisnis)
type MenuUsecase interface {
	FetchSemuaMenu() ([]*Menu, error)
	FetchMenuByID(id uint) (*Menu, error)
	TambahMenu(menu *Menu) error
	UpdateMenu(menu *Menu) error
	HapusMenu(id uint) error
}

// 3.Kontrak untuk Repository (Database)
type MenuRepository interface {
	GetAll() ([]*Menu, error)
	GetByID(id uint) (*Menu, error)
	Create(menu *Menu) error
	Update(menu *Menu) error
	Delete(id uint) error
}
