package model

import (
	"errors"
	"time"
)

var (
	ErrBahanBakuNotFound = errors.New("bahan baku tidak ditemukan")
	ErrStokTidakCukup    = errors.New("stok bahan baku tidak mencukupi")
)

// 1. Struct Entity BahanBaku
type BahanBaku struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	NamaBahan  string    `json:"nama_bahan" gorm:"type:varchar(100);not null"`
	HargaBeli  float64   `json:"harga_beli" gorm:"not null"`                    // Harga modal beli per satuan
	StokSistem float64   `json:"stok_sistem" gorm:"not null"`                   // Jumlah stok (desimal: misal 1.5 kg / 500 ml)
	Satuan     string    `json:"satuan" gorm:"type:varchar(20);default:'gram'"` // gram, ml, pcs, dll
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (BahanBaku) TableName() string {
	return "bahan_baku"
}

// 2. Struct Entity Log Operasional Gudang (Tahap 4)
type LogStokMasuk struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	BahanBakuID uint       `json:"bahan_baku_id" gorm:"not null;index"`
	JumlahMasuk float64    `json:"jumlah_masuk" gorm:"not null"`
	Pemasok     string     `json:"pemasok" gorm:"type:varchar(150)"`
	WaktuMasuk  time.Time  `json:"waktu_masuk" gorm:"autoCreateTime"`
	BahanBaku   *BahanBaku `json:"bahan_baku,omitempty" gorm:"foreignKey:BahanBakuID"`
}

func (LogStokMasuk) TableName() string {
	return "log_stok_masuk"
}

type LogBarangRusak struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	BahanBakuID     uint       `json:"bahan_baku_id" gorm:"not null;index"`
	JumlahTerbuang  float64    `json:"jumlah_terbuang" gorm:"not null"`
	Alasan          string     `json:"alasan" gorm:"type:text"`
	WaktuPencatatan time.Time  `json:"waktu_pencatatan" gorm:"autoCreateTime"`
	BahanBaku       *BahanBaku `json:"bahan_baku,omitempty" gorm:"foreignKey:BahanBakuID"`
}

func (LogBarangRusak) TableName() string {
	return "log_barang_rusak"
}

type StockOpname struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	BahanBakuID uint       `json:"bahan_baku_id" gorm:"not null;index"`
	StokSistem  float64    `json:"stok_sistem" gorm:"not null"` // Stok sebelum dihitung ulang
	StokFisik   float64    `json:"stok_fisik" gorm:"not null"`  // Hasil hitungan nyata di toko
	Selisih     float64    `json:"selisih" gorm:"not null"`     // StokFisik - StokSistem
	WaktuOpname time.Time  `json:"waktu_opname" gorm:"autoCreateTime"`
	BahanBaku   *BahanBaku `json:"bahan_baku,omitempty" gorm:"foreignKey:BahanBakuID"`
}

func (StockOpname) TableName() string {
	return "stock_opname"
}

// 3. DTO untuk Request Operasional Gudang
type CatatStokMasukRequest struct {
	BahanBakuID uint    `json:"bahan_baku_id" binding:"required"`
	JumlahMasuk float64 `json:"jumlah_masuk" binding:"required,gt=0"`
	Pemasok     string  `json:"pemasok"` // opsional — nama supplier jika ada
	Catatan     string  `json:"catatan"` // opsional — catatan tambahan
}

type CatatBarangRusakRequest struct {
	BahanBakuID    uint    `json:"bahan_baku_id" binding:"required"`
	JumlahTerbuang float64 `json:"jumlah_terbuang" binding:"required,gt=0"`
	Alasan         string  `json:"alasan" binding:"required,min=3"`
}

type InputStockOpnameRequest struct {
	BahanBakuID uint    `json:"bahan_baku_id" binding:"required"`
	StokFisik   float64 `json:"stok_fisik" binding:"required,gte=0"`
	Catatan     string  `json:"catatan"` // opsional — catatan hasil opname
}

// 4. Kontrak Interface Repository
type BahanBakuRepository interface {
	GetAll() ([]*BahanBaku, error)
	GetByID(id uint) (*BahanBaku, error)
	Create(bahan *BahanBaku) error
	Update(bahan *BahanBaku) error
	Delete(id uint) error
	UpdateStok(id uint, delta float64) error
	SetStok(id uint, newStok float64) error
}

type InventoryRepository interface {
	CreateLogStokMasuk(log *LogStokMasuk) error
	GetLogStokMasuk() ([]*LogStokMasuk, error)
	CreateLogBarangRusak(log *LogBarangRusak) error
	GetLogBarangRusak() ([]*LogBarangRusak, error)
	CreateStockOpname(opname *StockOpname) error
	GetStockOpname() ([]*StockOpname, error)
}

// 5. Kontrak Interface Usecase
type BahanBakuUsecase interface {
	FetchSemuaBahanBaku() ([]*BahanBaku, error)
	FetchBahanBakuByID(id uint) (*BahanBaku, error)
	TambahBahanBaku(bahan *BahanBaku) error
	UpdateBahanBaku(bahan *BahanBaku) error
	HapusBahanBaku(id uint) error
}

type InventoryUsecase interface {
	CatatStokMasuk(req *CatatStokMasukRequest) (*LogStokMasuk, error)
	GetRiwayatStokMasuk() ([]*LogStokMasuk, error)
	CatatBarangRusak(req *CatatBarangRusakRequest) (*LogBarangRusak, error)
	GetRiwayatBarangRusak() ([]*LogBarangRusak, error)
	InputStockOpname(req *InputStockOpnameRequest) (*StockOpname, error)
	GetRiwayatStockOpname() ([]*StockOpname, error)
}
