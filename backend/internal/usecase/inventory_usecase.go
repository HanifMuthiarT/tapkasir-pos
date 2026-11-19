package usecase

import (
	"errors"

	"tapkasir-pos/internal/model"
)

type inventoryUsecase struct {
	invRepo   model.InventoryRepository
	bahanRepo model.BahanBakuRepository
}

func NewInventoryUsecase(invRepo model.InventoryRepository, bahanRepo model.BahanBakuRepository) model.InventoryUsecase {
	return &inventoryUsecase{
		invRepo:   invRepo,
		bahanRepo: bahanRepo,
	}
}

// CatatStokMasuk mencatat penerimaan bahan dari supplier dan otomatis menambah stok sistem
func (u *inventoryUsecase) CatatStokMasuk(req *model.CatatStokMasukRequest) (*model.LogStokMasuk, error) {
	if req.JumlahMasuk <= 0 {
		return nil, errors.New("jumlah stok masuk harus lebih besar dari 0")
	}

	bahan, err := u.bahanRepo.GetByID(req.BahanBakuID)
	if err != nil {
		return nil, err
	}

	// 1. Tambah stok sistem secara atomic
	if err := u.bahanRepo.UpdateStok(req.BahanBakuID, req.JumlahMasuk); err != nil {
		return nil, err
	}

	// 2. Simpan histori ke tabel log_stok_masuk
	logMasuk := &model.LogStokMasuk{
		BahanBakuID: req.BahanBakuID,
		JumlahMasuk: req.JumlahMasuk,
		Pemasok:     req.Pemasok,
	}
	if err := u.invRepo.CreateLogStokMasuk(logMasuk); err != nil {
		return nil, err
	}

	logMasuk.BahanBaku = bahan
	return logMasuk, nil
}

func (u *inventoryUsecase) GetRiwayatStokMasuk() ([]*model.LogStokMasuk, error) {
	return u.invRepo.GetLogStokMasuk()
}

// CatatBarangRusak mencatat bahan rusak/basi/tumpah dan otomatis memotong stok sistem
func (u *inventoryUsecase) CatatBarangRusak(req *model.CatatBarangRusakRequest) (*model.LogBarangRusak, error) {
	if req.JumlahTerbuang <= 0 {
		return nil, errors.New("jumlah bahan terbuang harus lebih besar dari 0")
	}

	bahan, err := u.bahanRepo.GetByID(req.BahanBakuID)
	if err != nil {
		return nil, err
	}

	if bahan.StokSistem < req.JumlahTerbuang {
		return nil, model.ErrStokTidakCukup
	}

	// 1. Kurangi stok sistem
	if err := u.bahanRepo.UpdateStok(req.BahanBakuID, -req.JumlahTerbuang); err != nil {
		return nil, err
	}

	// 2. Simpan histori ke tabel log_barang_rusak
	logRusak := &model.LogBarangRusak{
		BahanBakuID:    req.BahanBakuID,
		JumlahTerbuang: req.JumlahTerbuang,
		Alasan:         req.Alasan,
	}
	if err := u.invRepo.CreateLogBarangRusak(logRusak); err != nil {
		return nil, err
	}

	logRusak.BahanBaku = bahan
	return logRusak, nil
}

func (u *inventoryUsecase) GetRiwayatBarangRusak() ([]*model.LogBarangRusak, error) {
	return u.invRepo.GetLogBarangRusak()
}

// InputStockOpname mencocokkan stok fisik di toko dengan stok sistem dan menghitung selisih
func (u *inventoryUsecase) InputStockOpname(req *model.InputStockOpnameRequest) (*model.StockOpname, error) {
	if req.StokFisik < 0 {
		return nil, errors.New("stok fisik tidak boleh bernilai negatif")
	}

	bahan, err := u.bahanRepo.GetByID(req.BahanBakuID)
	if err != nil {
		return nil, err
	}

	stokSebelumnya := bahan.StokSistem
	selisih := req.StokFisik - stokSebelumnya

	// 1. Sesuaikan stok sistem di tabel bahan_baku menjadi sama persis dengan hasil hitungan fisik
	if err := u.bahanRepo.SetStok(req.BahanBakuID, req.StokFisik); err != nil {
		return nil, err
	}

	// 2. Simpan catatan audit stock opname beserta nilai selisihnya
	opname := &model.StockOpname{
		BahanBakuID: req.BahanBakuID,
		StokSistem:  stokSebelumnya,
		StokFisik:   req.StokFisik,
		Selisih:     selisih,
	}
	if err := u.invRepo.CreateStockOpname(opname); err != nil {
		return nil, err
	}

	opname.BahanBaku = bahan
	return opname, nil
}

func (u *inventoryUsecase) GetRiwayatStockOpname() ([]*model.StockOpname, error) {
	return u.invRepo.GetStockOpname()
}
