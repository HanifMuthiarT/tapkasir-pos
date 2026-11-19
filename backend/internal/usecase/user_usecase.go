package usecase

import (
	"log"
	"time"

	"tapkasir-pos/internal/model"
	"tapkasir-pos/pkg/security"
)

type userUsecase struct {
	userRepo  model.UserRepository
	jwtSecret string
}

func NewUserUsecase(repo model.UserRepository, jwtSecret string) model.UserUsecase {
	if jwtSecret == "" {
		jwtSecret = "tapkasir-secret-key-2026"
	}
	return &userUsecase{
		userRepo:  repo,
		jwtSecret: jwtSecret,
	}
}

func (u *userUsecase) VerifyPIN(req *model.VerifyPINRequest) (*model.VerifyPINResponse, error) {
	user, err := u.userRepo.GetByRole(req.Role)
	if err != nil {
		return nil, err
	}

	// 1. Cek kecocokan PIN dengan hash di database
	if !security.CheckPINHash(req.PIN, user.PIN) {
		return nil, model.ErrInvalidPIN
	}

	// 2. Buat JWT Token dengan masa aktif 4 jam
	token, exp, err := security.GenerateToken(user.Role, u.jwtSecret, 4*time.Hour)
	if err != nil {
		return nil, err
	}

	return &model.VerifyPINResponse{
		Token:     token,
		Role:      user.Role,
		ExpiresAt: exp,
	}, nil
}

func (u *userUsecase) ChangePIN(req *model.ChangePINRequest) error {
	user, err := u.userRepo.GetByRole(req.Role)
	if err != nil {
		return err
	}

	// Validasi PIN lama
	if !security.CheckPINHash(req.OldPIN, user.PIN) {
		return model.ErrInvalidPIN
	}

	// Hash PIN baru
	hashedNewPIN, err := security.HashPIN(req.NewPIN)
	if err != nil {
		return err
	}

	return u.userRepo.UpdatePIN(req.Role, hashedNewPIN)
}

func (u *userUsecase) SeedDefaultUsers() error {
	count, err := u.userRepo.Count()
	if err != nil {
		return err
	}

	// Hanya isi data awal jika database masih kosong
	if count == 0 {
		defaultUsers := []struct {
			Role string
			PIN  string
		}{
			{Role: "owner", PIN: "111111"},
			{Role: "manager", PIN: "222222"},
			{Role: "cashier", PIN: "333333"},
		}

		for _, item := range defaultUsers {
			hashed, err := security.HashPIN(item.PIN)
			if err != nil {
				return err
			}
			user := &model.User{
				Role: item.Role,
				PIN:  hashed,
			}
			if err := u.userRepo.Create(user); err != nil {
				return err
			}
		}
		log.Println("Seeder berhasil: PIN awal (Owner: 111111, Manager: 222222, Kasir: 333333) siap digunakan.")
	}

	return nil
}