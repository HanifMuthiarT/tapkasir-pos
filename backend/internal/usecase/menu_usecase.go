package usecase

import "tapkasir-pos/internal/model"

type menuUsecase struct {
	menuRepo model.MenuRepository // Menggunakan interface
}

// Dependency Injection: Inject repository ke dalam usecase
func NewMenuUsecase(repo model.MenuRepository) model.MenuUsecase {
	return &menuUsecase{
		menuRepo: repo,
	}
}

func (u *menuUsecase) FetchSemuaMenu() ([]*model.Menu, error) {
	// Di sini Anda bisa menambahkan logika bisnis.
	// Misalnya: filter hanya menu yang statusnya "tersedia".

	// Memanggil repository
	menus, err := u.menuRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return menus, nil
}

func (u *menuUsecase) TambahMenu(menu *model.Menu) error {
	return u.menuRepo.Create(menu)
}

func (u *menuUsecase) FetchMenuByID(id uint) (*model.Menu, error) {
	return u.menuRepo.GetByID(id)
}

func (u *menuUsecase) UpdateMenu(menu *model.Menu) error {
	return u.menuRepo.Update(menu)
}

func (u *menuUsecase) HapusMenu(id uint) error {
	return u.menuRepo.Delete(id)
}
