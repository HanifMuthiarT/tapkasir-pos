package repository

import (
	"errors"

	"tapkasir-pos/internal/model"

	"gorm.io/gorm"
)

type bahanBakuRepository struct {
	db *gorm.DB
}

// NewBahanBakuRepository constructor Dependency Injection untuk BahanBakuRepository
func NewBahanBakuRepository(db *gorm.DB) model.BahanBakuRepository {
	return &bahanBakuRepository{db: db}
}

func (r *bahanBakuRepository) GetAll() ([]*model.BahanBaku, error) {
	var list []*model.BahanBaku
	err := r.db.Order("id desc").Find(&list).Error
	return list, err
}

func (r *bahanBakuRepository) GetByID(id uint) (*model.BahanBaku, error) {
	var bahan model.BahanBaku
	if err := r.db.First(&bahan, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrBahanBakuNotFound
		}
		return nil, err
	}
	return &bahan, nil
}

func (r *bahanBakuRepository) Create(bahan *model.BahanBaku) error {
	return r.db.Create(bahan).Error
}

func (r *bahanBakuRepository) Update(bahan *model.BahanBaku) error {
	if _, err := r.GetByID(bahan.ID); err != nil {
		return err
	}

	return r.db.Model(&model.BahanBaku{}).
		Where("id = ?", bahan.ID).
		Select("NamaBahan", "HargaBeli", "StokSistem", "Satuan").
		Updates(bahan).Error
}

func (r *bahanBakuRepository) Delete(id uint) error {
	result := r.db.Delete(&model.BahanBaku{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return model.ErrBahanBakuNotFound
	}
	return nil
}

// UpdateStok memodifikasi stok secara atomic langsung di SQL (mencegah race condition)
func (r *bahanBakuRepository) UpdateStok(id uint, delta float64) error {
	result := r.db.Model(&model.BahanBaku{}).
		Where("id = ?", id).
		Update("stok_sistem", gorm.Expr("stok_sistem + ?", delta))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return model.ErrBahanBakuNotFound
	}
	return nil
}

// SetStok memperbarui nilai stok sistem secara langsung (digunakan saat Stock Opname)
func (r *bahanBakuRepository) SetStok(id uint, newStok float64) error {
	result := r.db.Model(&model.BahanBaku{}).
		Where("id = ?", id).
		Update("stok_sistem", newStok)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return model.ErrBahanBakuNotFound
	}
	return nil
}
