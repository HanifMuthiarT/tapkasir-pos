package handler

import (
	"errors"
	"net/http"
	"strconv"

	"tapkasir-pos/internal/delivery/http/middleware"
	"tapkasir-pos/internal/model"

	"github.com/gin-gonic/gin"
)

type BahanBakuHandler struct {
	bahanUsecase model.BahanBakuUsecase
}

func NewBahanBakuHandler(r *gin.Engine, us model.BahanBakuUsecase, jwtSecret string) {
	handler := &BahanBakuHandler{
		bahanUsecase: us,
	}

	// 1. Route Publik / Monitoring Toko
	r.GET("/api/bahan-baku", handler.FetchSemuaBahanBaku)
	r.GET("/api/bahan-baku/:id", handler.FetchBahanBakuByID)

	// 2. Route Sensitif: Tambah, Edit, Hapus Bahan Baku (Wajib PIN Owner / Manager)
	sensitif := r.Group("/api/bahan-baku")
	sensitif.Use(middleware.RequireRole(jwtSecret, "owner", "manager"))
	{
		sensitif.POST("", handler.TambahBahanBaku)
		sensitif.PUT("/:id", handler.UpdateBahanBaku)
		sensitif.DELETE("/:id", handler.HapusBahanBaku)
	}
}

func (h *BahanBakuHandler) FetchSemuaBahanBaku(c *gin.Context) {
	data, err := h.bahanUsecase.FetchSemuaBahanBaku()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data bahan baku"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sukses",
		"data":    data,
	})
}

func (h *BahanBakuHandler) FetchBahanBakuByID(c *gin.Context) {
	id, ok := parseBahanBakuID(c)
	if !ok {
		return
	}

	bahan, err := h.bahanUsecase.FetchBahanBakuByID(id)
	if err != nil {
		respondBahanBakuError(c, err, "Gagal mengambil data bahan baku")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sukses",
		"data":    bahan,
	})
}

type simpanBahanBakuRequest struct {
	NamaBahan  string  `json:"nama_bahan" binding:"required"`
	HargaBeli  float64 `json:"harga_beli" binding:"gt=0"`
	StokSistem float64 `json:"stok_sistem" binding:"gte=0"`
	Satuan     string  `json:"satuan" binding:"required"`
}

func (h *BahanBakuHandler) TambahBahanBaku(c *gin.Context) {
	var req simpanBahanBakuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data bahan baku tidak valid. Periksa nama_bahan, harga_beli (>0), stok_sistem (>=0), dan satuan."})
		return
	}

	bahan := &model.BahanBaku{
		NamaBahan:  req.NamaBahan,
		HargaBeli:  req.HargaBeli,
		StokSistem: req.StokSistem,
		Satuan:     req.Satuan,
	}

	if err := h.bahanUsecase.TambahBahanBaku(bahan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Bahan baku berhasil ditambahkan",
		"data":    bahan,
	})
}

func (h *BahanBakuHandler) UpdateBahanBaku(c *gin.Context) {
	id, ok := parseBahanBakuID(c)
	if !ok {
		return
	}

	var req simpanBahanBakuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data bahan baku tidak valid"})
		return
	}

	bahan := &model.BahanBaku{
		ID:         id,
		NamaBahan:  req.NamaBahan,
		HargaBeli:  req.HargaBeli,
		StokSistem: req.StokSistem,
		Satuan:     req.Satuan,
	}

	if err := h.bahanUsecase.UpdateBahanBaku(bahan); err != nil {
		respondBahanBakuError(c, err, "Gagal memperbarui bahan baku")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Bahan baku berhasil diperbarui",
		"data":    bahan,
	})
}

func (h *BahanBakuHandler) HapusBahanBaku(c *gin.Context) {
	id, ok := parseBahanBakuID(c)
	if !ok {
		return
	}

	if err := h.bahanUsecase.HapusBahanBaku(id); err != nil {
		respondBahanBakuError(c, err, "Gagal menghapus bahan baku")
		return
	}

	c.Status(http.StatusNoContent)
}

func parseBahanBakuID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, strconv.IntSize)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID bahan baku tidak valid"})
		return 0, false
	}
	return uint(id), true
}

func respondBahanBakuError(c *gin.Context, err error, defaultMsg string) {
	if errors.Is(err, model.ErrBahanBakuNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bahan baku tidak ditemukan"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": defaultMsg})
}