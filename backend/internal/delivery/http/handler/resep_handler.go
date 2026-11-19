package handler

import (
	"errors"
	"net/http"
	"strconv"

	"tapkasir-pos/internal/delivery/http/middleware"
	"tapkasir-pos/internal/model"

	"github.com/gin-gonic/gin"
)

type ResepHandler struct {
	resepUsecase model.ResepUsecase
}

func NewResepHandler(r *gin.Engine, us model.ResepUsecase, jwtSecret string) {
	handler := &ResepHandler{
		resepUsecase: us,
	}

	// 1. Route Publik: Melihat komposisi resep & takaran menu
	r.GET("/api/menu/:id/resep", handler.GetResepByMenuID)

	// 2. Route Sensitif: Atur Resep, Hapus Item, & Analisis HPP Semua Menu (Wajib PIN Owner / Manager)
	sensitif := r.Group("/api")
	sensitif.Use(middleware.RequireRole(jwtSecret, "owner", "manager"))
	{
		sensitif.POST("/menu/:id/resep", handler.AturResepMenu)
		sensitif.DELETE("/resep/:id", handler.HapusItemResep)
		sensitif.GET("/resep/hpp", handler.GetAllMenuHPP)
	}
}

func (h *ResepHandler) GetResepByMenuID(c *gin.Context) {
	menuID, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	res, err := h.resepUsecase.GetResepByMenuID(menuID)
	if err != nil {
		if errors.Is(err, model.ErrMenuNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Menu tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data resep menu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sukses",
		"data":    res,
	})
}

func (h *ResepHandler) AturResepMenu(c *gin.Context) {
	menuID, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	var req model.AturResepRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data resep tidak valid. Pastikan minimal ada 1 bahan dengan takaran > 0.",
		})
		return
	}

	res, err := h.resepUsecase.AturResepMenu(menuID, &req)
	if err != nil {
		if errors.Is(err, model.ErrMenuNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Menu tidak ditemukan"})
			return
		}
		if errors.Is(err, model.ErrBahanBakuNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Salah satu bahan baku tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan resep menu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Resep menu berhasil diatur dan HPP telah dikalkulasi ulang",
		"data":    res,
	})
}

func (h *ResepHandler) HapusItemResep(c *gin.Context) {
	resepID, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	if err := h.resepUsecase.HapusItemResep(resepID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus item resep"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ResepHandler) GetAllMenuHPP(c *gin.Context) {
	data, err := h.resepUsecase.GetAllMenuHPP()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil laporan HPP menu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sukses",
		"data":    data,
	})
}

func parseParamID(c *gin.Context, paramName string) (uint, bool) {
	id, err := strconv.ParseUint(c.Param(paramName), 10, strconv.IntSize)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter tidak valid"})
		return 0, false
	}
	return uint(id), true
}
