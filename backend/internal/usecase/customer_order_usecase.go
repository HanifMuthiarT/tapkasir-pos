package usecase

import (
	"fmt"

	"tapkasir-pos/internal/model"
)

type customerOrderUsecase struct {
	pesananRepo model.PesananRepository
	menuRepo    model.MenuRepository
	resepRepo   model.ResepRepository
	bahanRepo   model.BahanBakuRepository
}

func NewCustomerOrderUsecase(
	pesananRepo model.PesananRepository,
	menuRepo model.MenuRepository,
	resepRepo model.ResepRepository,
	bahanRepo model.BahanBakuRepository,
) model.CustomerOrderUsecase {
	return &customerOrderUsecase{
		pesananRepo: pesananRepo,
		menuRepo:    menuRepo,
		resepRepo:   resepRepo,
		bahanRepo:   bahanRepo,
	}
}

// FetchMenuTersedia mengambil daftar menu yang aktif/tersedia untuk pelanggan
func (u *customerOrderUsecase) FetchMenuTersedia() ([]*model.Menu, error) {
	menus, err := u.menuRepo.GetAll()
	if err != nil {
		return nil, err
	}

	tersedia := make([]*model.Menu, 0)
	for _, m := range menus {
		if m.StatusKetersediaan == "tersedia" {
			tersedia = append(tersedia, m)
		}
	}
	return tersedia, nil
}

// CheckoutPesanan memproses pesanan meja pelanggan (QRIS / Tunai)
func (u *customerOrderUsecase) CheckoutPesanan(req *model.CustomerCheckoutRequest) (*model.CheckoutResponse, error) {
	var totalHarga float64
	detailItems := make([]model.DetailPesanan, 0)

	// 1. Validasi setiap menu dan hitung total biaya transaksi
	for _, it := range req.Items {
		menu, err := u.menuRepo.GetByID(it.MenuID)
		if err != nil {
			return nil, err
		}

		if menu.StatusKetersediaan != "tersedia" {
			return nil, model.ErrMenuNotAvailable
		}

		subTotal := menu.HargaJual * float64(it.Kuantitas)
		totalHarga += subTotal

		detailItems = append(detailItems, model.DetailPesanan{
			MenuID:    menu.ID,
			Kuantitas: it.Kuantitas,
			SubTotal:  subTotal,
		})
	}

	// 2. Buat entitas pesanan baru
	pesanan := &model.Pesanan{
		NomorMeja:        req.NomorMeja,
		NamaPelanggan:    req.NamaPelanggan,
		NomorTelepon:     req.NomorTelepon,
		MetodePembayaran: req.MetodePembayaran,
		StatusPembayaran: "unpaid",
		StatusPesanan:    "menunggu",
		TotalHarga:       totalHarga,
		DetailPesanan:    detailItems,
	}

	if err := u.pesananRepo.Create(pesanan); err != nil {
		return nil, err
	}

	// 3. Tentukan instruksi alur pembayaran berdasarkan pilihan pelanggan
	var instruksi string
	var qrisPayload string

	if req.MetodePembayaran == "qris" {
		instruksi = "Silakan lakukan pembayaran QRIS melalui e-wallet atau mobile banking Anda. Setelah terverifikasi, pesanan otomatis dikirim ke dapur."
		// Simulasi string QRIS / tautan image QR code
		qrisPayload = fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?size=250x250&data=TAPKASIR-MEJA-%d-PESANAN-%d-TOTAL-%.0f", req.NomorMeja, pesanan.ID, totalHarga)
	} else {
		instruksi = fmt.Sprintf("Pesanan berhasil dicatat! Silakan segera menuju ke meja kasir dengan menyebutkan Nama: '%s' dan Nomor Meja: %d untuk melakukan pembayaran tunai.", req.NamaPelanggan, req.NomorMeja)
	}

	return &model.CheckoutResponse{
		Pesanan:     pesanan,
		Instruksi:   instruksi,
		QRISPayload: qrisPayload,
	}, nil
}

// LacakPesanan memungkinkan pelanggan melihat status pesanannya secara real-time dari HP
func (u *customerOrderUsecase) LacakPesanan(id uint) (*model.Pesanan, error) {
	return u.pesananRepo.GetByID(id)
}

// BatalkanPesananOlehPelanggan memungkinkan pelanggan membatalkan pesanannya sendiri
// Hanya bisa dibatalkan jika masih berstatus "menunggu" dan belum dibayar (unpaid)
func (u *customerOrderUsecase) BatalkanPesananOlehPelanggan(id uint, req *model.CustomerBatalPesananRequest) error {
	pesanan, err := u.pesananRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Hanya boleh batal jika belum dibayar dan masih menunggu
	if pesanan.StatusPembayaran == "paid" {
		return model.ErrPesananAlreadyPaid
	}
	if pesanan.StatusPesanan == "dibatalkan" {
		return model.ErrPesananAlreadyVoid
	}
	// Tidak bisa batal jika sudah diproses dapur
	if pesanan.StatusPesanan == "dimasak" || pesanan.StatusPesanan == "diantar" || pesanan.StatusPesanan == "selesai" {
		return model.ErrPesananInvalidStatus
	}

	// Ubah status jadi dibatalkan
	if err := u.pesananRepo.UpdateStatus(id, "unpaid", "dibatalkan"); err != nil {
		return err
	}

	// Catat ke log void
	logVoid := &model.LogVoid{
		PesananID:   id,
		AlasanBatal: req.AlasanBatal,
	}
	return u.pesananRepo.CreateLogVoid(logVoid)
}

// SimulasiBayarQRIS memvalidasi pembayaran QRIS dan OTOMATIS MEMOTONG STOK BAHAN BAKU di gudang sesuai resep!
func (u *customerOrderUsecase) SimulasiBayarQRIS(id uint) (*model.Pesanan, error) {
	pesanan, err := u.pesananRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if pesanan.StatusPembayaran == "paid" {
		return nil, model.ErrPesananAlreadyPaid
	}

	if pesanan.MetodePembayaran != "qris" {
		return nil, model.ErrBukanPembayaranQRIS
	}

	// 1. Ubah status pesanan menjadi PAID dan langsung ke status 'dimasak' (dapur)
	if err := u.pesananRepo.UpdateStatus(id, "paid", "dimasak"); err != nil {
		return nil, err
	}

	// 2. OTOMATIS POTONG STOK BAHAN BAKU SESUAI TAKARAN RESEP MENU YANG DIPESAN!
	for _, item := range pesanan.DetailPesanan {
		reseps, err := u.resepRepo.GetByMenuID(item.MenuID)
		if err != nil {
			continue
		}

		for _, r := range reseps {
			jumlahPakai := r.Takaran * float64(item.Kuantitas)
			// Kurangi stok sistem bahan baku
			_ = u.bahanRepo.UpdateStok(r.BahanBakuID, -jumlahPakai)
		}
	}

	// Kembalikan data pesanan terbaru
	return u.pesananRepo.GetByID(id)
}
