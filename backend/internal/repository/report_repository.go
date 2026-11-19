package repository

import (
	"time"

	"tapkasir-pos/internal/model"

	"gorm.io/gorm"
)

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) model.ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) GetLaporanOmzet(start, end time.Time) (*model.RingkasanOmzet, error) {
	var totalOmzet float64
	var totalTransaksi int64
	var omzetTunai float64
	var omzetQRIS float64
	var totalHPP float64
	var totalWaste float64

	// 1. Total Transaksi & Total Omzet
	_ = r.db.Model(&model.Pesanan{}).
		Where("status_pembayaran = 'paid' AND waktu_pesanan BETWEEN ? AND ?", start, end).
		Count(&totalTransaksi).Error

	_ = r.db.Model(&model.Pesanan{}).
		Where("status_pembayaran = 'paid' AND waktu_pesanan BETWEEN ? AND ?", start, end).
		Select("COALESCE(SUM(total_harga), 0)").
		Scan(&totalOmzet).Error

	// 2. Omzet per Metode Pembayaran
	_ = r.db.Model(&model.Pesanan{}).
		Where("status_pembayaran = 'paid' AND metode_pembayaran = 'tunai' AND waktu_pesanan BETWEEN ? AND ?", start, end).
		Select("COALESCE(SUM(total_harga), 0)").
		Scan(&omzetTunai).Error

	_ = r.db.Model(&model.Pesanan{}).
		Where("status_pembayaran = 'paid' AND metode_pembayaran = 'qris' AND waktu_pesanan BETWEEN ? AND ?", start, end).
		Select("COALESCE(SUM(total_harga), 0)").
		Scan(&omzetQRIS).Error

	// 3. Total HPP Bahan Baku dari Menu yang Terjual
	hppQuery := `
		SELECT COALESCE(SUM(dp.kuantitas * r.takaran * b.harga_beli), 0)
		FROM detail_pesanan dp
		JOIN pesanan p ON dp.pesanan_id = p.id
		JOIN resep r ON dp.menu_id = r.menu_id
		JOIN bahan_baku b ON r.bahan_baku_id = b.id
		WHERE p.status_pembayaran = 'paid' AND p.waktu_pesanan BETWEEN ? AND ?
	`
	_ = r.db.Raw(hppQuery, start, end).Scan(&totalHPP).Error

	// 4. Total Kerugian Bahan Rusak / Basi (Waste)
	wasteQuery := `
		SELECT COALESCE(SUM(l.jumlah_terbuang * b.harga_beli), 0)
		FROM log_barang_rusak l
		JOIN bahan_baku b ON l.bahan_baku_id = b.id
		WHERE l.waktu_pencatatan BETWEEN ? AND ?
	`
	_ = r.db.Raw(wasteQuery, start, end).Scan(&totalWaste).Error

	labaKotor := totalOmzet - totalHPP
	labaBersih := labaKotor - totalWaste

	return &model.RingkasanOmzet{
		TotalOmzet:         totalOmzet,
		TotalTransaksi:     totalTransaksi,
		OmzetTunai:         omzetTunai,
		OmzetQRIS:          omzetQRIS,
		TotalHPP:           totalHPP,
		TotalKerugianWaste: totalWaste,
		EstimasiLabaKotor:  labaKotor,
		EstimasiLabaBersih: labaBersih,
		PeriodeMulai:       start.Format("2006-01-02 15:04:05"),
		PeriodeSelesai:     end.Format("2006-01-02 15:04:05"),
	}, nil
}

func (r *reportRepository) GetBestSeller(start, end time.Time, limit int) ([]model.BestSellerMenu, error) {
	var results []model.BestSellerMenu
	if limit <= 0 {
		limit = 10
	}

	bestSellerQuery := `
		SELECT dp.menu_id, m.nama_menu, SUM(dp.kuantitas) as total_terjual, SUM(dp.sub_total) as total_omzet
		FROM detail_pesanan dp
		JOIN pesanan p ON dp.pesanan_id = p.id
		JOIN menu m ON dp.menu_id = m.id
		WHERE p.status_pembayaran = 'paid' AND p.waktu_pesanan BETWEEN ? AND ?
		GROUP BY dp.menu_id, m.nama_menu
		ORDER BY total_terjual DESC
		LIMIT ?
	`

	err := r.db.Raw(bestSellerQuery, start, end, limit).Scan(&results).Error
	return results, err
}
