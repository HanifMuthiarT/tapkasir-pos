package usecase

import (
	"math"

	"tapkasir-pos/internal/model"
)

type resepUsecase struct {
	resepRepo model.ResepRepository
	menuRepo  model.MenuRepository
	bahanRepo model.BahanBakuRepository
}

func NewResepUsecase(resepRepo model.ResepRepository, menuRepo model.MenuRepository, bahanRepo model.BahanBakuRepository) model.ResepUsecase {
	return &resepUsecase{
		resepRepo: resepRepo,
		menuRepo:  menuRepo,
		bahanRepo: bahanRepo,
	}
}

func (u *resepUsecase) GetResepByMenuID(menuID uint) (*model.MenuHPPResponse, error) {
	menu, err := u.menuRepo.GetByID(menuID)
	if err != nil {
		return nil, err
	}

	reseps, err := u.resepRepo.GetByMenuID(menuID)
	if err != nil {
		return nil, err
	}

	var totalHPP float64
	items := make([]model.ResepItemDetail, 0)

	for _, r := range reseps {
		if r.BahanBaku == nil {
			continue
		}
		subtotal := r.Takaran * r.BahanBaku.HargaBeli
		totalHPP += subtotal

		items = append(items, model.ResepItemDetail{
			ResepID:     r.ID,
			BahanBakuID: r.BahanBakuID,
			NamaBahan:   r.BahanBaku.NamaBahan,
			Takaran:     r.Takaran,
			Satuan:      r.BahanBaku.Satuan,
			HargaBeli:   r.BahanBaku.HargaBeli,
			SubtotalHPP: math.Round(subtotal*100) / 100,
		})
	}

	totalHPP = math.Round(totalHPP*100) / 100
	margin := menu.HargaJual - totalHPP
	var marginPersen float64
	if menu.HargaJual > 0 {
		marginPersen = math.Round((margin/menu.HargaJual)*10000) / 100
	}

	return &model.MenuHPPResponse{
		MenuID:         menu.ID,
		NamaMenu:       menu.NamaMenu,
		HargaJual:      menu.HargaJual,
		TotalHPP:       totalHPP,
		EstimasiMargin: margin,
		MarginPersen:   marginPersen,
		KomposisiBahan: items,
	}, nil
}

func (u *resepUsecase) AturResepMenu(menuID uint, req *model.AturResepRequest) (*model.MenuHPPResponse, error) {
	// Pastikan menu valid
	if _, err := u.menuRepo.GetByID(menuID); err != nil {
		return nil, err
	}

	// Validasi semua bahan baku yang diinput benar-benar ada
	for _, item := range req.Items {
		if _, err := u.bahanRepo.GetByID(item.BahanBakuID); err != nil {
			return nil, err
		}
	}

	// Hapus komposisi resep lama untuk menu ini
	if err := u.resepRepo.DeleteByMenuID(menuID); err != nil {
		return nil, err
	}

	// Simpan komposisi resep baru
	for _, item := range req.Items {
		resep := &model.Resep{
			MenuID:      menuID,
			BahanBakuID: item.BahanBakuID,
			Takaran:     item.Takaran,
		}
		if err := u.resepRepo.Create(resep); err != nil {
			return nil, err
		}
	}

	// Kembalikan rincian HPP terupdate
	return u.GetResepByMenuID(menuID)
}

func (u *resepUsecase) HapusItemResep(id uint) error {
	return u.resepRepo.Delete(id)
}

func (u *resepUsecase) GetAllMenuHPP() ([]*model.MenuHPPResponse, error) {
	menus, err := u.menuRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var results []*model.MenuHPPResponse
	for _, m := range menus {
		hppRes, err := u.GetResepByMenuID(m.ID)
		if err != nil {
			continue
		}
		results = append(results, hppRes)
	}

	return results, nil
}
