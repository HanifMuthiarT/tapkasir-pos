package usecase

import (
	"tapkasir-pos/internal/model"
)

type voidUsecase struct {
	pesananRepo model.PesananRepository
	resepRepo   model.ResepRepository
	bahanRepo   model.BahanBakuRepository
}

func NewVoidUsecase(
	pesananRepo model.PesananRepository,
	resepRepo model.ResepRepository,
	bahanRepo model.BahanBakuRepository,
) model.VoidUsecase {
	return &voidUsecase{
		pesananRepo: pesananRepo,
		resepRepo:   resepRepo,
		bahanRepo:   bahanRepo,
	}
}

// BatalkanPesanan memproses pembatalan pesanan sensitif, mencatat log void, dan mengembalikan stok jika sudah lunas
func (u *voidUsecase) BatalkanPesanan(pesananID uint, req *model.VoidPesananRequest) (*model.LogVoid, error) {
	pesanan, err := u.pesananRepo.GetByID(pesananID)
	if err != nil {
		return nil, err
	}

	if pesanan.StatusPesanan == "dibatalkan" {
		return nil, model.ErrPesananAlreadyVoid
	}

	// 1. Jika pesanan sebelumnya sudah lunas (PAID), kembalikan stok bahan baku ke gudang (Rollback Stock)!
	if pesanan.StatusPembayaran == "paid" {
		for _, item := range pesanan.DetailPesanan {
			reseps, err := u.resepRepo.GetByMenuID(item.MenuID)
			if err != nil {
				continue
			}

			for _, r := range reseps {
				kembalikanStok := r.Takaran * float64(item.Kuantitas)
				// Nilai positif akan menambah kembali sisa stok di gudang
				_ = u.bahanRepo.UpdateStok(r.BahanBakuID, kembalikanStok)
			}
		}
	}

	// 2. Ubah status pesanan menjadi dibatalkan
	if err := u.pesananRepo.UpdateStatus(pesananID, "unpaid", "dibatalkan"); err != nil {
		return nil, err
	}

	// 3. Simpan rekam jejak pembatalan di tabel log_void
	logVoid := &model.LogVoid{
		PesananID:   pesananID,
		AlasanBatal: req.AlasanBatal,
	}
	if err := u.pesananRepo.CreateLogVoid(logVoid); err != nil {
		return nil, err
	}

	logVoid.Pesanan = pesanan
	return logVoid, nil
}

func (u *voidUsecase) GetRiwayatVoid() ([]*model.LogVoid, error) {
	return u.pesananRepo.GetLogVoid()
}
