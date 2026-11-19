package model

import (
	"errors"
	"time"
)

var (
	ErrPesananNotFound      = errors.New("pesanan tidak ditemukan")
	ErrMenuNotAvailable     = errors.New("salah satu menu sedang tidak tersedia")
	ErrPesananAlreadyPaid   = errors.New("pesanan sudah berstatus lunas (PAID)")
	ErrBukanPembayaranQRIS  = errors.New("metode pembayaran pesanan ini bukan QRIS")
	ErrBukanPembayaranTunai = errors.New("metode pembayaran pesanan ini bukan tunai")
	ErrPesananInvalidStatus = errors.New("transisi status pesanan tidak valid")
	ErrPesananAlreadyVoid   = errors.New("pesanan ini sudah dibatalkan sebelumnya (VOID)")
	ErrUangKurang           = errors.New("uang yang diterima kurang dari total harga belanja")
)

// 1. Struct Entity Pesanan
type Pesanan struct {
	ID               uint            `json:"id" gorm:"primaryKey"`
	NomorMeja        int             `json:"nomor_meja" gorm:"not null"`
	NamaPelanggan    string          `json:"nama_pelanggan" gorm:"type:varchar(100);not null"`
	NomorTelepon     string          `json:"nomor_telepon" gorm:"type:varchar(20)"`
	MetodePembayaran string          `json:"metode_pembayaran" gorm:"type:varchar(20);not null"`
	StatusPembayaran string          `json:"status_pembayaran" gorm:"type:varchar(30);default:'unpaid'"`
	StatusPesanan    string          `json:"status_pesanan" gorm:"type:varchar(30);default:'menunggu'"`
	TotalHarga       float64         `json:"total_harga" gorm:"not null"`
	WaktuPesanan     time.Time       `json:"waktu_pesanan" gorm:"autoCreateTime"`
	DetailPesanan    []DetailPesanan `json:"detail_pesanan" gorm:"foreignKey:PesananID"`
	LogVoid          *LogVoid        `json:"log_void,omitempty" gorm:"foreignKey:PesananID"`
}

func (Pesanan) TableName() string {
	return "pesanan"
}

// 2. Struct Entity DetailPesanan
type DetailPesanan struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	PesananID uint    `json:"pesanan_id" gorm:"not null;index"`
	MenuID    uint    `json:"menu_id" gorm:"not null;index"`
	Kuantitas int     `json:"kuantitas" gorm:"not null"`
	SubTotal  float64 `json:"sub_total" gorm:"not null"`
	Menu      *Menu   `json:"menu,omitempty" gorm:"foreignKey:MenuID"`
}

func (DetailPesanan) TableName() string {
	return "detail_pesanan"
}

// 3. Struct Entity LogVoid (Tahap 7)
type LogVoid struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	PesananID   uint      `json:"pesanan_id" gorm:"not null;uniqueIndex"`
	AlasanBatal string    `json:"alasan_batal" gorm:"type:text;not null"`
	WaktuBatal  time.Time `json:"waktu_batal" gorm:"autoCreateTime"`
	Pesanan     *Pesanan  `json:"pesanan,omitempty" gorm:"foreignKey:PesananID"`
}

func (LogVoid) TableName() string {
	return "log_void"
}

// 4. Data Transfer Objects (DTO) untuk Pelanggan (Tahap 5)
type CustomerOrderItemRequest struct {
	MenuID    uint `json:"menu_id" binding:"required"`
	Kuantitas int  `json:"kuantitas" binding:"required,gt=0"`
}

type CustomerCheckoutRequest struct {
	NomorMeja        int                        `json:"nomor_meja" binding:"required,gt=0"`
	NamaPelanggan    string                     `json:"nama_pelanggan" binding:"required,min=2"`
	NomorTelepon     string                     `json:"nomor_telepon" binding:"required,min=9,max=15"`
	MetodePembayaran string                     `json:"metode_pembayaran" binding:"required,oneof=qris tunai"`
	Items            []CustomerOrderItemRequest `json:"items" binding:"required,min=1,dive"`
}

type CheckoutResponse struct {
	Pesanan     *Pesanan `json:"pesanan"`
	Instruksi   string   `json:"instruksi"`
	QRISPayload string   `json:"qris_payload,omitempty"`
}

// 5. Data Transfer Objects (DTO) untuk Kasir (Tahap 6)
type KonfirmasiBayarTunaiRequest struct {
	UangDiterima float64 `json:"uang_diterima" binding:"required,gt=0"`
}

type KonfirmasiBayarTunaiResponse struct {
	Pesanan      *Pesanan `json:"pesanan"`
	TotalHarga   float64  `json:"total_harga"`
	UangDiterima float64  `json:"uang_diterima"`
	Kembalian    float64  `json:"kembalian"`
}

type UpdateStatusPesananRequest struct {
	StatusPesanan string `json:"status_pesanan" binding:"required,oneof=menunggu dimasak diantar selesai"`
}

type StrukItemDetail struct {
	NamaMenu  string  `json:"nama_menu"`
	Harga     float64 `json:"harga"`
	Kuantitas int     `json:"kuantitas"`
	SubTotal  float64 `json:"sub_total"`
}

type StrukTransaksiResponse struct {
	NamaToko      string            `json:"nama_toko"`
	AlamatToko    string            `json:"alamat_toko"`
	NomorPesanan  uint              `json:"nomor_pesanan"`
	NomorMeja     int               `json:"nomor_meja"`
	NamaPelanggan string            `json:"nama_pelanggan"`
	WaktuPesanan  string            `json:"waktu_pesanan"`
	MetodeBayar   string            `json:"metode_bayar"`
	StatusBayar   string            `json:"status_bayar"`
	Items         []StrukItemDetail `json:"items"`
	TotalHarga    float64           `json:"total_harga"`
	UangDiterima  float64           `json:"uang_diterima,omitempty"`
	Kembalian     float64           `json:"kembalian,omitempty"`
}

// 6. Data Transfer Objects (DTO) untuk Void / Pembatalan Pesanan (Tahap 7)
type VoidPesananRequest struct {
	AlasanBatal string `json:"alasan_batal" binding:"required,min=3"`
}

// DTO untuk pembatalan oleh pelanggan sendiri (lebih sederhana, tanpa alasan panjang)
type CustomerBatalPesananRequest struct {
	AlasanBatal string `json:"alasan_batal" binding:"required,min=3"`
}

// 7. Kontrak Interface untuk Repository (Database)
type PesananRepository interface {
	Create(pesanan *Pesanan) error
	GetByID(id uint) (*Pesanan, error)
	GetAll(statusPembayaran string, statusPesanan string) ([]*Pesanan, error)
	UpdateStatus(id uint, statusPembayaran string, statusPesanan string) error
	CreateLogVoid(log *LogVoid) error
	GetLogVoid() ([]*LogVoid, error)
	GetLogVoidByPesananID(pesananID uint) (*LogVoid, error)
}

// 8. Kontrak Interface untuk Usecase
type CustomerOrderUsecase interface {
	FetchMenuTersedia() ([]*Menu, error)
	CheckoutPesanan(req *CustomerCheckoutRequest) (*CheckoutResponse, error)
	LacakPesanan(id uint) (*Pesanan, error)
	SimulasiBayarQRIS(id uint) (*Pesanan, error)
	BatalkanPesananOlehPelanggan(id uint, req *CustomerBatalPesananRequest) error
}

type CashierOrderUsecase interface {
	GetSemuaPesanan(statusPembayaran string, statusPesanan string) ([]*Pesanan, error)
	GetPesananByID(id uint) (*Pesanan, error)
	KonfirmasiPembayaranTunai(id uint, req *KonfirmasiBayarTunaiRequest) (*KonfirmasiBayarTunaiResponse, error)
	UpdateStatusPesanan(id uint, req *UpdateStatusPesananRequest) (*Pesanan, error)
	GetStrukPesanan(id uint, uangDiterima float64) (*StrukTransaksiResponse, error)
}

type VoidUsecase interface {
	BatalkanPesanan(pesananID uint, req *VoidPesananRequest) (*LogVoid, error)
	GetRiwayatVoid() ([]*LogVoid, error)
}
