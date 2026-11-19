package usecase

import (
	"fmt"
	"time"

	"tapkasir-pos/internal/model"

	"github.com/xuri/excelize/v2"
)

type reportUsecase struct {
	reportRepo  model.ReportRepository
	pesananRepo model.PesananRepository
	invRepo     model.InventoryRepository
}

func NewReportUsecase(
	reportRepo model.ReportRepository,
	pesananRepo model.PesananRepository,
	invRepo model.InventoryRepository,
) model.ReportUsecase {
	return &reportUsecase{
		reportRepo:  reportRepo,
		pesananRepo: pesananRepo,
		invRepo:     invRepo,
	}
}

func parseDateRange(startStr, endStr string) (time.Time, time.Time) {
	now := time.Now()
	// Default 30 hari terakhir
	start := now.AddDate(0, 0, -30)
	end := now

	if startStr != "" {
		if t, err := time.Parse("2006-01-02", startStr); err == nil {
			start = t
		}
	}
	if endStr != "" {
		if t, err := time.Parse("2006-01-02", endStr); err == nil {
			// Jadikan akhir hari (23:59:59)
			end = t.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		}
	}

	return start, end
}

func (u *reportUsecase) GetLaporanLengkap(startStr, endStr string) (*model.LaporanLengkapResponse, error) {
	start, end := parseDateRange(startStr, endStr)

	// 1. Ringkasan Omzet & Finansial
	ringkasan, err := u.reportRepo.GetLaporanOmzet(start, end)
	if err != nil {
		return nil, err
	}

	// 2. Produk Terlaris (Top 10)
	bestSeller, err := u.reportRepo.GetBestSeller(start, end, 10)
	if err != nil {
		return nil, err
	}

	// 3. Audit Stok & Keamanan
	voidLogs, _ := u.pesananRepo.GetLogVoid()
	wasteLogs, _ := u.invRepo.GetLogBarangRusak()
	opnameLogs, _ := u.invRepo.GetStockOpname()

	audit := model.AuditStokResponse{
		TotalLogVoid:     len(voidLogs),
		TotalLogWaste:    len(wasteLogs),
		TotalStockOpname: len(opnameLogs),
		LogsVoid:         voidLogs,
		LogsWaste:        wasteLogs,
		LogsOpname:       opnameLogs,
	}

	return &model.LaporanLengkapResponse{
		Ringkasan:  *ringkasan,
		BestSeller: bestSeller,
		AuditStok:  audit,
	}, nil
}

// ExportExcel mengekspor seluruh data laporan ke file Microsoft Excel (.xlsx)
func (u *reportUsecase) ExportExcel(startStr, endStr string) ([]byte, error) {
	laporan, err := u.GetLaporanLengkap(startStr, endStr)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	defer f.Close()

	// --- SHEET 1: Ringkasan Finansial ---
	sheetFinansial := "Ringkasan Finansial"
	f.SetSheetName("Sheet1", sheetFinansial)

	f.SetCellValue(sheetFinansial, "A1", "LAPORAN KINERJA BISNIS TAPKASIR POS")
	f.SetCellValue(sheetFinansial, "A2", fmt.Sprintf("Periode: %s s/d %s", laporan.Ringkasan.PeriodeMulai, laporan.Ringkasan.PeriodeSelesai))

	headers := [][]interface{}{
		{"Indikator Keuangan", "Nilai (Rp / Transaksi)"},
		{"Total Omzet Kotor", laporan.Ringkasan.TotalOmzet},
		{"Total Transaksi Sukses", laporan.Ringkasan.TotalTransaksi},
		{"Omzet Pembayaran Tunai", laporan.Ringkasan.OmzetTunai},
		{"Omzet Pembayaran QRIS", laporan.Ringkasan.OmzetQRIS},
		{"Total HPP (Modal Bahan Terjual)", laporan.Ringkasan.TotalHPP},
		{"Total Kerugian Bahan Basi/Rusak (Waste)", laporan.Ringkasan.TotalKerugianWaste},
		{"Estimasi Laba Kotor Toko", laporan.Ringkasan.EstimasiLabaKotor},
		{"Estimasi Laba Bersih Toko", laporan.Ringkasan.EstimasiLabaBersih},
	}

	for i, row := range headers {
		f.SetCellValue(sheetFinansial, fmt.Sprintf("A%d", i+4), row[0])
		f.SetCellValue(sheetFinansial, fmt.Sprintf("B%d", i+4), row[1])
	}

	// --- SHEET 2: Produk Terlaris ---
	sheetBestSeller := "Menu Terlaris"
	f.NewSheet(sheetBestSeller)
	f.SetCellValue(sheetBestSeller, "A1", "Peringkat")
	f.SetCellValue(sheetBestSeller, "B1", "Nama Menu")
	f.SetCellValue(sheetBestSeller, "C1", "Total Porsi Terjual")
	f.SetCellValue(sheetBestSeller, "D1", "Total Omzet Penjualan (Rp)")

	for idx, item := range laporan.BestSeller {
		row := idx + 2
		f.SetCellValue(sheetBestSeller, fmt.Sprintf("A%d", row), idx+1)
		f.SetCellValue(sheetBestSeller, fmt.Sprintf("B%d", row), item.NamaMenu)
		f.SetCellValue(sheetBestSeller, fmt.Sprintf("C%d", row), item.TotalTerjual)
		f.SetCellValue(sheetBestSeller, fmt.Sprintf("D%d", row), item.TotalOmzet)
	}

	// --- SHEET 3: Audit Stok & Void ---
	sheetAudit := "Audit Keamanan & Void"
	f.NewSheet(sheetAudit)
	f.SetCellValue(sheetAudit, "A1", "ID Void")
	f.SetCellValue(sheetAudit, "B1", "Nomor Pesanan")
	f.SetCellValue(sheetAudit, "C1", "Alasan Pembatalan")
	f.SetCellValue(sheetAudit, "D1", "Waktu Pembatalan")

	for idx, v := range laporan.AuditStok.LogsVoid {
		row := idx + 2
		f.SetCellValue(sheetAudit, fmt.Sprintf("A%d", row), v.ID)
		f.SetCellValue(sheetAudit, fmt.Sprintf("B%d", row), v.PesananID)
		f.SetCellValue(sheetAudit, fmt.Sprintf("C%d", row), v.AlasanBatal)
		f.SetCellValue(sheetAudit, fmt.Sprintf("D%d", row), v.WaktuBatal.Format("2006-01-02 15:04:05"))
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
