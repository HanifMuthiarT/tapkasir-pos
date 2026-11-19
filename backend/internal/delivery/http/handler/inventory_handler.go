package handler

import (
	"errors"
	"net/http"

	"tapkasir-pos/internal/delivery/http/middleware"
	"tapkasir-pos/internal/model"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	invUsecase model.InventoryUsecase
}

func NewInventoryHandler(r *gin.Engine, us model.InventoryUsecase, jwtSecret string) {
	handler := &InventoryHandler{
		invUsecase: us,
	}

	// Seluruh operasional gudang dan inventori wajib PIN Manager atau Owner
	invGroup := r.Group("/api/inventory")
	invGroup.Use(middleware.RequireRole(jwtSecret, "owner", "manager"))
	{
		// 1. Stok Masuk (Restock Penerimaan Supplier)
		invGroup.POST("/stok-masuk", handler.CatatStokMasuk)
		invGroup.GET("/stok-masuk", handler.GetRiwayatStokMasuk)

		// 2. Bahan Baku Rusak / Basi / Terbuang (Waste Log)
		invGroup.POST("/barang-rusak", handler.CatatBarangRusak)
		invGroup.GET("/barang-rusak", handler.GetRiwayatBarangRusak)

		// 3. Stock Opname (Audit Fisik vs Sistem)
		invGroup.POST("/stock-opname", handler.InputStockOpname)
		invGroup.GET("/stock-opname", handler.GetRiwayatStockOpname)
	}
}

func (h *InventoryHandler) CatatStokMasuk(c *gin.Context) {
	var req model.CatatStokMasukRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data stok masuk tidak valid. Periksa bahan_baku_id, jumlah_masuk (>0), dan pemasok."})
		return
	}

	res, err := h.invUsecase.CatatStokMasuk(&req)
	if err != nil {
		if errors.Is(err, model.ErrBahanBakuNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bahan baku tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pencatatan stok masuk berhasil. Stok sistem otomatis bertambah.",
		"data":    res,
	})
}

func (h *InventoryHandler) GetRiwayatStokMasuk(c *gin.Context) {
	data, err := h.invUsecase.GetRiwayatStokMasuk()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil riwayat stok masuk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sukses",
		"data":    data,
	})
}

func (h *InventoryHandler) CatatBarangRusak(c *gin.Context) {
	var req model.CatatBarangRusakRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data barang rusak tidak valid. Periksa bahan_baku_id, jumlah_terbuang (>0), dan alasan."})
		return
	}

	res, err := h.invUsecase.CatatBarangRusak(&req)
	if err != nil {
		if errors.Is(err, model.ErrBahanBakuNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bahan baku tidak ditemukan"})
			return
		}
		if errors.Is(err, model.ErrStokTidakCukup) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal: Jumlah bahan terbuang melebihi sisa stok sistem yang ada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pencatatan barang rusak berhasil. Stok sistem otomatis dipotong.",
		"data":    res,
	})
}

func (h *InventoryHandler) GetRiwayatBarangRusak(c *gin.Context) {
	data, err := h.invUsecase.GetRiwayatBarangRusak()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil riwayat barang rusak"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sukses",
		"data":    data,
	})
}

func (h *InventoryHandler) InputStockOpname(c *gin.Context) {
	var req model.InputStockOpnameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data stock opname tidak valid. Periksa bahan_baku_id dan stok_fisik (>=0)."})
		return
	}

	res, err := h.invUsecase.InputStockOpname(&req)
	if err != nil {
		if errors.Is(err, model.ErrBahanBakuNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bahan baku tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Stock opname berhasil disimpan dan stok sistem telah disesuaikan",
		"data":    res,
	})
}

func (h *InventoryHandler) GetRiwayatStockOpname(c *gin.Context) {
	data, err := h.invUsecase.GetRiwayatStockOpname()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil riwayat stock opname"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sukses",
		"data":    data,
	})
}
