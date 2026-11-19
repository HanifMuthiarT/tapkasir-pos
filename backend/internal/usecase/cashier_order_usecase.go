package usecase

import (
	"tapkasir-pos/internal/model"
)

type cashierOrderUsecase struct {
	pesananRepo model.PesananRepository
	resepRepo   model.ResepRepository
	bahanRepo   model.BahanBakuRepository
}

func NewCashierOrderUsecase(
	pesananRepo model.PesananRepository,
	resepRepo model.ResepRepository,
	bahanRepo model.BahanBakuRepository,
) model.CashierOrderUsecase {
	return &cashierOrderUsecase{
		pesananRepo: pesananRepo,
		resepRepo:   resepRepo,
		bahanRepo:   bahanRepo,
	}
}

func (u *cashierOrderUsecase) GetSemuaPesanan(statusPembayaran string, statusPesanan string) ([]*model.Pesanan, error) {
	return u.pesananRepo.GetAll(statusPembayaran, statusPesanan)
}

func (u *cashierOrderUsecase) GetPesananByID(id uint) (*model.Pesanan, error) {
	return u.pesananRepo.GetByID(id)
}

// KonfirmasiPembayaranTunai: Kasir menerima uang cash, menghitung kembalian, dan OTOMATIS MEMOTONG STOK BAHAN BAKU!
func (u *cashierOrderUsecase) KonfirmasiPembayaranTunai(id uint, req *model.KonfirmasiBayarTunaiRequest) (*model.KonfirmasiBayarTunaiResponse, error) {
	pesanan, err := u.pesananRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if pesanan.StatusPembayaran == "paid" {
		return nil, model.ErrPesananAlreadyPaid
	}

	if pesanan.MetodePembayaran != "tunai" {
		return nil, model.ErrBukanPembayaranTunai
	}

	if req.UangDiterima < pesanan.TotalHarga {
		return nil, model.ErrUangKurang
	}

	kembalian := req.UangDiterima - pesanan.TotalHarga

	// 1. Ubah status pesanan menjadi PAID dan langsung ke status 'dimasak' (dapur)
	if err := u.pesananRepo.UpdateStatus(id, "paid", "dimasak"); err != nil {
		return nil, err
	}

	// 2. OTOMATIS POTONG STOK BAHAN BAKU DI GUDANG SESUAI RESEP!
	for _, item := range pesanan.DetailPesanan {
		reseps, err := u.resepRepo.GetByMenuID(item.MenuID)
		if err != nil {
			continue
		}

		for _, r := range reseps {
			jumlahPakai := r.Takaran * float64(item.Kuantitas)
			_ = u.bahanRepo.UpdateStok(r.BahanBakuID, -jumlahPakai)
		}
	}

	updated, err := u.pesananRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &model.KonfirmasiBayarTunaiResponse{
		Pesanan:      updated,
		TotalHarga:   pesanan.TotalHarga,
		UangDiterima: req.UangDiterima,
		Kembalian:    kembalian,
	}, nil
}

// UpdateStatusPesanan: Kasir/Dapur memperbarui status masakan (menunggu -> dimasak -> diantar -> selesai)
func (u *cashierOrderUsecase) UpdateStatusPesanan(id uint, req *model.UpdateStatusPesananRequest) (*model.Pesanan, error) {
	pesanan, err := u.pesananRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if pesanan.StatusPesanan == "dibatalkan" {
		return nil, model.ErrPesananAlreadyVoid
	}

	if err := u.pesananRepo.UpdateStatus(id, "", req.StatusPesanan); err != nil {
		return nil, err
	}

	return u.pesananRepo.GetByID(id)
}

// GetStrukPesanan: Menghasilkan data nota terformat rapi untuk cetak struk kasir
func (u *cashierOrderUsecase) GetStrukPesanan(id uint, uangDiterima float64) (*model.StrukTransaksiResponse, error) {
	pesanan, err := u.pesananRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	items := make([]model.StrukItemDetail, 0)
	for _, it := range pesanan.DetailPesanan {
		nama := "Menu"
		harga := 0.0
		if it.Menu != nil {
			nama = it.Menu.NamaMenu
			harga = it.Menu.HargaJual
		}
		items = append(items, model.StrukItemDetail{
			NamaMenu:  nama,
			Harga:     harga,
			Kuantitas: it.Kuantitas,
			SubTotal:  it.SubTotal,
		})
	}

	var kembalian float64
	if uangDiterima >= pesanan.TotalHarga {
		kembalian = uangDiterima - pesanan.TotalHarga
	}

	return &model.StrukTransaksiResponse{
		NamaToko:      "TAPKASIR COFFEE & EATERY",
		AlamatToko:    "Jl. Danau Sunter Utara No. 12, Jakarta",
		NomorPesanan:  pesanan.ID,
		NomorMeja:     pesanan.NomorMeja,
		NamaPelanggan: pesanan.NamaPelanggan,
		WaktuPesanan:  pesanan.WaktuPesanan.Format("02 Jan 2006 15:04:05"),
		MetodeBayar:   pesanan.MetodePembayaran,
		StatusBayar:   pesanan.StatusPembayaran,
		Items:         items,
		TotalHarga:    pesanan.TotalHarga,
		UangDiterima:  uangDiterima,
		Kembalian:     kembalian,
	}, nil
}
