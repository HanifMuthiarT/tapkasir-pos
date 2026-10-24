package repository

import (
	"errors"

	"tapkasir-pos/internal/model"

	"gorm.io/gorm"
)

type menuRepository struct {
	db *gorm.DB
}

// Dependency Injection: Inject database ke dalam repository
func NewMenuRepository(db *gorm.DB) model.MenuRepository {
	return &menuRepository{db}
}

func (r *menuRepository) GetAll() ([]*model.Menu, error) {
	var menu []*model.Menu
	// Query murni tanpa ada logika bisnis
	err := r.db.Find(&menu).Error
	return menu, err
}

func (r *menuRepository) GetByID(id uint) (*model.Menu, error) {
	var menu model.Menu
	if err := r.db.First(&menu, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrMenuNotFound
		}
		return nil, err
	}

	return &menu, nil
}

func (r *menuRepository) Create(menu *model.Menu) error {
	return r.db.Create(menu).Error
}

func (r *menuRepository) Update(menu *model.Menu) error {
	if _, err := r.GetByID(menu.ID); err != nil {
		return err
	}

	return r.db.Model(&model.Menu{}).
		Where("id = ?", menu.ID).
		Select("NamaMenu", "HargaJual", "StatusKetersediaan").
		Updates(menu).Error
}

func (r *menuRepository) Delete(id uint) error {
	result := r.db.Delete(&model.Menu{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return model.ErrMenuNotFound
	}

	return nil
}
