package repository

import (
	"tapkasir-pos/internal/model"

	"gorm.io/gorm"
)

type resepRepository struct {
	db *gorm.DB
}

func NewResepRepository(db *gorm.DB) model.ResepRepository {
	return &resepRepository{db: db}
}

func (r *resepRepository) GetByMenuID(menuID uint) ([]*model.Resep, error) {
	var list []*model.Resep
	err := r.db.Preload("BahanBaku").Where("menu_id = ?", menuID).Find(&list).Error
	return list, err
}

func (r *resepRepository) Create(resep *model.Resep) error {
	return r.db.Create(resep).Error
}

func (r *resepRepository) DeleteByMenuID(menuID uint) error {
	return r.db.Where("menu_id = ?", menuID).Delete(&model.Resep{}).Error
}

func (r *resepRepository) Delete(id uint) error {
	result := r.db.Delete(&model.Resep{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return model.ErrResepNotFound
	}
	return nil
}
