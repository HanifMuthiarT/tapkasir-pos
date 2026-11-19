package model

import (
	"errors"
)

var (
	ErrResepNotFound = errors.New("resep bahan tidak ditemukan")
)

// 1. Struct Entity Resep (Tabel Penghubung antara Menu dan BahanBaku)
type Resep struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	MenuID      uint       `json:"menu_id" gorm:"not null;index"`
	BahanBakuID uint       `json:"bahan_baku_id" gorm:"not null;index"`
	Takaran     float64    `json:"takaran" gorm:"not null"` // Jumlah takaran bahan yang dipakai per 1 porsi menu
	Menu        *Menu      `json:"menu,omitempty" gorm:"foreignKey:MenuID"`
	BahanBaku   *BahanBaku `json:"bahan_baku,omitempty" gorm:"foreignKey:BahanBakuID"`
}

func (Resep) TableName() string {
	return "resep"
}

// 2. Data Transfer Objects (DTO) untuk Request & Response Kalkulasi HPP
type AturResepItem struct {
	BahanBakuID uint    `json:"bahan_baku_id" binding:"required"`
	Takaran     float64 `json:"takaran" binding:"gt=0"`
}

type AturResepRequest struct {
	Items []AturResepItem `json:"items" binding:"required,min=1,dive"`
}

type ResepItemDetail struct {
	ResepID     uint    `json:"resep_id"`
	BahanBakuID uint    `json:"bahan_baku_id"`
	NamaBahan   string  `json:"nama_bahan"`
	Takaran     float64 `json:"takaran"`
	Satuan      string  `json:"satuan"`
	HargaBeli   float64 `json:"harga_beli"`   // Modal per satuan
	SubtotalHPP float64 `json:"subtotal_hpp"` // Takaran * HargaBeli
}

type MenuHPPResponse struct {
	MenuID         uint              `json:"menu_id"`
	NamaMenu       string            `json:"nama_menu"`
	HargaJual      float64           `json:"harga_jual"`
	TotalHPP       float64           `json:"total_hpp"`       // Total biaya modal bahan per porsi
	EstimasiMargin float64           `json:"estimasi_margin"` // HargaJual - TotalHPP
	MarginPersen   float64           `json:"margin_persen"`   // (EstimasiMargin / HargaJual) * 100
	KomposisiBahan []ResepItemDetail `json:"komposisi_bahan"`
}

// 3. Kontrak Interface untuk Repository (Database)
type ResepRepository interface {
	GetByMenuID(menuID uint) ([]*Resep, error)
	Create(resep *Resep) error
	DeleteByMenuID(menuID uint) error
	Delete(id uint) error
}

// 4. Kontrak Interface untuk Usecase (Logika Bisnis & Perhitungan HPP)
type ResepUsecase interface {
	GetResepByMenuID(menuID uint) (*MenuHPPResponse, error)
	AturResepMenu(menuID uint, req *AturResepRequest) (*MenuHPPResponse, error)
	HapusItemResep(id uint) error
	GetAllMenuHPP() ([]*MenuHPPResponse, error)
}
