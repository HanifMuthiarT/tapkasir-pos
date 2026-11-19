package model

import "time"

// 1. Data Transfer Objects (DTO) untuk Ringkasan Finansial
type RingkasanOmzet struct {
	TotalOmzet         float64 `json:"total_omzet"`
	TotalTransaksi     int64   `json:"total_transaksi"`
	OmzetTunai         float64 `json:"omzet_tunai"`
	OmzetQRIS          float64 `json:"omzet_qris"`
	TotalHPP           float64 `json:"total_hpp"`            // Total modal bahan baku dari pesanan terjual
	TotalKerugianWaste float64 `json:"total_kerugian_waste"` // Total biaya bahan yang rusak/basi
	EstimasiLabaKotor  float64 `json:"estimasi_laba_kotor"`  // TotalOmzet - TotalHPP
	EstimasiLabaBersih float64 `json:"estimasi_laba_bersih"` // EstimasiLabaKotor - TotalKerugianWaste
	PeriodeMulai       string  `json:"periode_mulai"`
	PeriodeSelesai     string  `json:"periode_selesai"`
}

// 2. DTO Produk Terlaris (Best Seller)
type BestSellerMenu struct {
	MenuID       uint    `json:"menu_id"`
	NamaMenu     string  `json:"nama_menu"`
	TotalTerjual int64   `json:"total_terjual"`
	TotalOmzet   float64 `json:"total_omzet"`
}

// 3. DTO Audit Stok & Keamanan
type AuditStokResponse struct {
	TotalLogVoid     int               `json:"total_log_void"`
	TotalLogWaste    int               `json:"total_log_waste"`
	TotalStockOpname int               `json:"total_stock_opname"`
	LogsVoid         []*LogVoid        `json:"logs_void"`
	LogsWaste        []*LogBarangRusak `json:"logs_waste"`
	LogsOpname       []*StockOpname    `json:"logs_opname"`
}

// 4. DTO Response Laporan Lengkap
type LaporanLengkapResponse struct {
	Ringkasan  RingkasanOmzet    `json:"ringkasan"`
	BestSeller []BestSellerMenu  `json:"best_seller"`
	AuditStok  AuditStokResponse `json:"audit_stok"`
}

// 5. Kontrak Interface Repository
type ReportRepository interface {
	GetLaporanOmzet(start, end time.Time) (*RingkasanOmzet, error)
	GetBestSeller(start, end time.Time, limit int) ([]BestSellerMenu, error)
}

// 6. Kontrak Interface Usecase
type ReportUsecase interface {
	GetLaporanLengkap(startStr, endStr string) (*LaporanLengkapResponse, error)
	ExportExcel(startStr, endStr string) ([]byte, error)
}
