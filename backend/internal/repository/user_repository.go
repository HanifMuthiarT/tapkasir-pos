package repository

import (
	"errors"

	"tapkasir-pos/internal/model"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) model.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByRole(role string) (*model.User, error) {
	var user model.User
	err := r.db.Where("role = ?", role).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) UpdatePIN(role string, hashedPIN string) error {
	result := r.db.Model(&model.User{}).
		Where("role = ?", role).
		Update("pin", hashedPIN)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return model.ErrUserNotFound
	}
	return nil
}

func (r *userRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.User{}).Count(&count).Error
	return count, err
}