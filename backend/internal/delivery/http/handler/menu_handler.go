package handler

import (
	"errors"
	"net/http"
	"strconv"

	"tapkasir-pos/internal/delivery/http/middleware"
	"tapkasir-pos/internal/model"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	menuUsecase model.MenuUsecase
}

// Dependency Injection: Inject usecase ke handler
func NewMenuHandler(r *gin.Engine, us model.MenuUsecase, jwtSecret string) {
	handler := &MenuHandler{
		menuUsecase: us,
	}
	// 1. Route Publik: Pelanggan scan QR & Layar Kasir standby bebas melihat menu
	r.GET("/api/menu", handler.FetchMenu)
	r.GET("/api/menu/:id", handler.FetchMenuByID)
	// 2. Route Sensitif: Tambah, edit harga/status, dan hapus menu membutuhkan PIN Manager atau Owner
	sensitiveMenu := r.Group("/api/menu")
	sensitiveMenu.Use(middleware.RequireRole(jwtSecret, "owner", "manager"))
	{
		sensitiveMenu.POST("", handler.TambahMenu)
		sensitiveMenu.PUT("/:id", handler.UpdateMenu)
		sensitiveMenu.DELETE("/:id", handler.HapusMenu)
	}
}

func (h *MenuHandler) FetchMenu(c *gin.Context) {
	// Memanggil usecase
	menu, err := h.menuUsecase.FetchSemuaMenu()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data menu"})
		return
	}

	// Mengembalikan response JSON
	c.JSON(http.StatusOK, gin.H{
		"message": "Sukses",
		"data":    menu,
	})
}

type tambahMenuRequest struct {
	NamaMenu           string  `json:"nama_menu" binding:"required"`
	HargaJual          float64 `json:"harga_jual" binding:"gt=0"`
	StatusKetersediaan string  `json:"status_ketersediaan" binding:"required"`
}

func (h *MenuHandler) TambahMenu(c *gin.Context) {
	var request tambahMenuRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data menu tidak valid"})
		return
	}

	menu := &model.Menu{
		NamaMenu:           request.NamaMenu,
		HargaJual:          request.HargaJual,
		StatusKetersediaan: request.StatusKetersediaan,
	}
	if err := h.menuUsecase.TambahMenu(menu); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan menu"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Menu berhasil ditambahkan",
		"data":    menu,
	})
}

func (h *MenuHandler) FetchMenuByID(c *gin.Context) {
	id, ok := parseMenuID(c)
	if !ok {
		return
	}

	menu, err := h.menuUsecase.FetchMenuByID(id)
	if err != nil {
		respondMenuError(c, err, "Gagal mengambil data menu")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sukses",
		"data":    menu,
	})
}

func (h *MenuHandler) UpdateMenu(c *gin.Context) {
	id, ok := parseMenuID(c)
	if !ok {
		return
	}

	var request tambahMenuRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data menu tidak valid"})
		return
	}

	menu := &model.Menu{
		ID:                 id,
		NamaMenu:           request.NamaMenu,
		HargaJual:          request.HargaJual,
		StatusKetersediaan: request.StatusKetersediaan,
	}
	if err := h.menuUsecase.UpdateMenu(menu); err != nil {
		respondMenuError(c, err, "Gagal memperbarui menu")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Menu berhasil diperbarui",
		"data":    menu,
	})
}

func (h *MenuHandler) HapusMenu(c *gin.Context) {
	id, ok := parseMenuID(c)
	if !ok {
		return
	}

	if err := h.menuUsecase.HapusMenu(id); err != nil {
		respondMenuError(c, err, "Gagal menghapus menu")
		return
	}

	c.Status(http.StatusNoContent)
}

func parseMenuID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, strconv.IntSize)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID menu tidak valid"})
		return 0, false
	}

	return uint(id), true
}

func respondMenuError(c *gin.Context, err error, message string) {
	if errors.Is(err, model.ErrMenuNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu tidak ditemukan"})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": message})
}
