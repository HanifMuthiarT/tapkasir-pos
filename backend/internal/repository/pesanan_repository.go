package repository

import (
	"errors"

	"tapkasir-pos/internal/model"

	"gorm.io/gorm"
)

type pesananRepository struct {
	db *gorm.DB
}

func NewPesananRepository(db *gorm.DB) model.PesananRepository {
	return &pesananRepository{db: db}
}

func (r *pesananRepository) Create(pesanan *model.Pesanan) error {
	// GORM otomatis menyimpan entitas Pesanan beserta seluruh DetailPesanan di dalamnya
	return r.db.Create(pesanan).Error
}

func (r *pesananRepository) GetByID(id uint) (*model.Pesanan, error) {
	var pesanan model.Pesanan
	err := r.db.Preload("DetailPesanan.Menu").
		Preload("LogVoid").
		First(&pesanan, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrPesananNotFound
		}
		return nil, err
	}

	return &pesanan, nil
}

func (r *pesananRepository) GetAll(statusPembayaran string, statusPesanan string) ([]*model.Pesanan, error) {
	var list []*model.Pesanan
	query := r.db.Preload("DetailPesanan.Menu").Preload("LogVoid").Order("id desc")

	if statusPembayaran != "" {
		query = query.Where("status_pembayaran = ?", statusPembayaran)
	}
	if statusPesanan != "" {
		query = query.Where("status_pesanan = ?", statusPesanan)
	}

	err := query.Find(&list).Error
	return list, err
}

func (r *pesananRepository) UpdateStatus(id uint, statusPembayaran string, statusPesanan string) error {
	updates := make(map[string]interface{})
	if statusPembayaran != "" {
		updates["status_pembayaran"] = statusPembayaran
	}
	if statusPesanan != "" {
		updates["status_pesanan"] = statusPesanan
	}

	result := r.db.Model(&model.Pesanan{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return model.ErrPesananNotFound
	}

	return nil
}

func (r *pesananRepository) CreateLogVoid(log *model.LogVoid) error {
	return r.db.Create(log).Error
}

func (r *pesananRepository) GetLogVoid() ([]*model.LogVoid, error) {
	var list []*model.LogVoid
	err := r.db.Preload("Pesanan.DetailPesanan.Menu").Order("id desc").Find(&list).Error
	return list, err
}

func (r *pesananRepository) GetLogVoidByPesananID(pesananID uint) (*model.LogVoid, error) {
	var log model.LogVoid
	err := r.db.Where("pesanan_id = ?", pesananID).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}
