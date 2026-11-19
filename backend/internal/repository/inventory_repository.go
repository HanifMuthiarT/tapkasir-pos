package repository

import (
	"tapkasir-pos/internal/model"

	"gorm.io/gorm"
)

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) model.InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) CreateLogStokMasuk(log *model.LogStokMasuk) error {
	return r.db.Create(log).Error
}

func (r *inventoryRepository) GetLogStokMasuk() ([]*model.LogStokMasuk, error) {
	var list []*model.LogStokMasuk
	err := r.db.Preload("BahanBaku").Order("id desc").Find(&list).Error
	return list, err
}

func (r *inventoryRepository) CreateLogBarangRusak(log *model.LogBarangRusak) error {
	return r.db.Create(log).Error
}

func (r *inventoryRepository) GetLogBarangRusak() ([]*model.LogBarangRusak, error) {
	var list []*model.LogBarangRusak
	err := r.db.Preload("BahanBaku").Order("id desc").Find(&list).Error
	return list, err
}

func (r *inventoryRepository) CreateStockOpname(opname *model.StockOpname) error {
	return r.db.Create(opname).Error
}

func (r *inventoryRepository) GetStockOpname() ([]*model.StockOpname, error) {
	var list []*model.StockOpname
	err := r.db.Preload("BahanBaku").Order("id desc").Find(&list).Error
	return list, err
}