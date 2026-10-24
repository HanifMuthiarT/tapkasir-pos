package model

import "time"

type Pesanan struct {
	ID               uint            `json:"id" gorm:"primaryKey"`
	NomorMeja        int             `json:"nomor_meja" gorm:"not null"`
	NamaPelanggan    string          `json:"nama_pelanggan" gorm:"type:varchar(100)"`
	NomorTelepon     string          `json:"nomor_telepon" gorm:"type:varchar(20)"`
	MetodePembayaran string          `json:"metode_pembayaran" gorm:"type:enum('qris', 'tunai');not null"`
	StatusPembayaran string          `json:"status_pembayaran" gorm:"type:enum('unpaid', 'paid');default:'unpaid'"`
	StatusPesanan    string          `json:"status_pesanan" gorm:"type:enum('menunggu', 'dimasak', 'diantar');default:'menunggu'"`
	TotalHarga       float64         `json:"total_harga" gorm:"not null"`
	WaktuPesanan     time.Time       `json:"waktu_pesanan" gorm:"autoCreateTime"`
	
	// Relasi ke tabel DetailPesanan
	DetailPesanan    []DetailPesanan `json:"detail_pesanan" gorm:"foreignKey:PesananID"`
}

type DetailPesanan struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	PesananID uint    `json:"pesanan_id" gorm:"not null"`
	MenuID    uint    `json:"menu_id" gorm:"not null"`
	Kuantitas int     `json:"kuantitas" gorm:"not null"`
	SubTotal  float64 `json:"sub_total" gorm:"not null"` // Menyimpan harga saat transaksi terjadi
}

type LogVoid struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	PesananID   uint      `json:"pesanan_id" gorm:"not null;unique"` // 1 to 1 relasi dengan Pesanan
	AlasanBatal string    `json:"alasan_batal" gorm:"type:text;not null"`
	WaktuBatal  time.Time `json:"waktu_batal" gorm:"autoCreateTime"`
}