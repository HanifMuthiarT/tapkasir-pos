package model

import "time"

type BahanBaku struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	NamaBahan  string    `json:"nama_bahan" gorm:"type:varchar(100);not null"`
	HargaBeli  float64   `json:"harga_beli" gorm:"not null"`
	StokSistem float64   `json:"stok_sistem" gorm:"not null"` // Menggunakan float untuk takaran desimal (misal: 1.5 Kg)
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type LogStokMasuk struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	BahanBakuID uint      `json:"bahan_baku_id" gorm:"not null"`
	JumlahMasuk float64   `json:"jumlah_masuk" gorm:"not null"`
	Pemasok     string    `json:"pemasok" gorm:"type:varchar(150)"`
	WaktuMasuk  time.Time `json:"waktu_masuk" gorm:"autoCreateTime"`
}

type LogBarangRusak struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	BahanBakuID     uint      `json:"bahan_baku_id" gorm:"not null"`
	JumlahTerbuang  float64   `json:"jumlah_terbuang" gorm:"not null"`
	Alasan          string    `json:"alasan" gorm:"type:text"`
	WaktuPencatatan time.Time `json:"waktu_pencatatan" gorm:"autoCreateTime"`
}

type StockOpname struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	BahanBakuID uint      `json:"bahan_baku_id" gorm:"not null"`
	StokSistem  float64   `json:"stok_sistem" gorm:"not null"`
	StokFisik   float64   `json:"stok_fisik" gorm:"not null"`
	Selisih     float64   `json:"selisih" gorm:"not null"`
	WaktuOpname time.Time `json:"waktu_opname" gorm:"autoCreateTime"`
}