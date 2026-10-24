package model

type Resep struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	MenuID      uint    `json:"menu_id" gorm:"not null"`
	BahanBakuID uint    `json:"bahan_baku_id" gorm:"not null"`
	Takaran     float64 `json:"takaran" gorm:"not null"`
}