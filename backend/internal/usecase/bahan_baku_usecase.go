package usecase

import (
	"errors"

	"tapkasir-pos/internal/model"
)

type bahanBakuUsecase struct {
	bahanRepo model.BahanBakuRepository
}

func NewBahanBakuUsecase(repo model.BahanBakuRepository) model.BahanBakuUsecase {
	return &bahanBakuUsecase{
		bahanRepo: repo,
	}
}

func (u *bahanBakuUsecase) FetchSemuaBahanBaku() ([]*model.BahanBaku, error) {
	return u.bahanRepo.GetAll()
}

func (u *bahanBakuUsecase) FetchBahanBakuByID(id uint) (*model.BahanBaku, error) {
	return u.bahanRepo.GetByID(id)
}

func (u *bahanBakuUsecase) TambahBahanBaku(bahan *model.BahanBaku) error {
	if bahan.HargaBeli <= 0 {
		return errors.New("harga beli harus lebih besar dari 0")
	}
	if bahan.StokSistem < 0 {
		return errors.New("stok sistem tidak boleh bernilai negatif")
	}
	if bahan.Satuan == "" {
		bahan.Satuan = "gram"
	}

	return u.bahanRepo.Create(bahan)
}

func (u *bahanBakuUsecase) UpdateBahanBaku(bahan *model.BahanBaku) error {
	if bahan.HargaBeli <= 0 {
		return errors.New("harga beli harus lebih besar dari 0")
	}
	if bahan.StokSistem < 0 {
		return errors.New("stok sistem tidak boleh bernilai negatif")
	}

	return u.bahanRepo.Update(bahan)
}

func (u *bahanBakuUsecase) HapusBahanBaku(id uint) error {
	return u.bahanRepo.Delete(id)
}