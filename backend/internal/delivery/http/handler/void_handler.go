package handler

import (
	"errors"
	"net/http"

	"tapkasir-pos/internal/delivery/http/middleware"
	"tapkasir-pos/internal/model"

	"github.com/gin-gonic/gin"
)

type VoidHandler struct {
	voidUsecase model.VoidUsecase
}

func NewVoidHandler(r *gin.Engine, us model.VoidUsecase, jwtSecret string) {
	handler := &VoidHandler{
		voidUsecase: us,
	}

	// Operasional Pembatalan Pesanan Wajib PIN Manager atau Owner
	voidGroup := r.Group("/api/void")
	voidGroup.Use(middleware.RequireRole(jwtSecret, "owner", "manager"))
	{
		voidGroup.POST("/orders/:id", handler.BatalkanPesanan)
		voidGroup.GET("/logs", handler.GetRiwayatVoid)
	}
}

func (h *VoidHandler) BatalkanPesanan(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	var req model.VoidPesananRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Alasan pembatalan pesanan wajib diisi (minimal 3 karakter)"})
		return
	}

	res, err := h.voidUsecase.BatalkanPesanan(id, &req)
	if err != nil {
		if errors.Is(err, model.ErrPesananNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan tidak ditemukan"})
			return
		}
		if errors.Is(err, model.ErrPesananAlreadyVoid) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Pesanan ini sudah dibatalkan sebelumnya"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses pembatalan pesanan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pesanan berhasil dibatalkan (VOID), stok bahan baku telah dikembalikan ke gudang, dan log pembatalan telah dicatat.",
		"data":    res,
	})
}

func (h *VoidHandler) GetRiwayatVoid(c *gin.Context) {
	data, err := h.voidUsecase.GetRiwayatVoid()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil log pembatalan pesanan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sukses",
		"total":   len(data),
		"data":    data,
	})
}
