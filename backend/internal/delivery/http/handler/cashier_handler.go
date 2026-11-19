package handler

import (
	"errors"
	"net/http"
	"strconv"

	"tapkasir-pos/internal/model"

	"github.com/gin-gonic/gin"
)

type CashierHandler struct {
	cashierUsecase model.CashierOrderUsecase
}

func NewCashierHandler(r *gin.Engine, us model.CashierOrderUsecase) {
	handler := &CashierHandler{
		cashierUsecase: us,
	}

	// Layar Kasir Standby: Selalu siap digunakan staf tanpa login berulang kali
	cashierGroup := r.Group("/api/cashier")
	{
		cashierGroup.GET("/orders", handler.GetSemuaPesanan)
		cashierGroup.GET("/orders/:id", handler.GetPesananByID)
		cashierGroup.POST("/orders/:id/bayar-tunai", handler.KonfirmasiPembayaranTunai)
		cashierGroup.PATCH("/orders/:id/status", handler.UpdateStatusPesanan)
		cashierGroup.GET("/orders/:id/struk", handler.GetStrukPesanan)
	}
}

func (h *CashierHandler) GetSemuaPesanan(c *gin.Context) {
	statusBayar := c.Query("status_pembayaran")
	statusPesanan := c.Query("status_pesanan")

	data, err := h.cashierUsecase.GetSemuaPesanan(statusBayar, statusPesanan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar pesanan masuk"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sukses",
		"total":   len(data),
		"data":    data,
	})
}

func (h *CashierHandler) GetPesananByID(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	pesanan, err := h.cashierUsecase.GetPesananByID(id)
	if err != nil {
		if errors.Is(err, model.ErrPesananNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil rincian pesanan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sukses",
		"data":    pesanan,
	})
}

func (h *CashierHandler) KonfirmasiPembayaranTunai(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	var req model.KonfirmasiBayarTunaiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nominal uang_diterima tidak valid"})
		return
	}

	res, err := h.cashierUsecase.KonfirmasiPembayaranTunai(id, &req)
	if err != nil {
		if errors.Is(err, model.ErrPesananNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan tidak ditemukan"})
			return
		}
		if errors.Is(err, model.ErrPesananAlreadyPaid) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Pesanan ini sudah berstatus lunas sebelumnya"})
			return
		}
		if errors.Is(err, model.ErrBukanPembayaranTunai) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Metode pembayaran pesanan ini bukan Tunai"})
			return
		}
		if errors.Is(err, model.ErrUangKurang) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Uang yang diterima kurang dari total harga pesanan!"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengonfirmasi pembayaran tunai"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pembayaran tunai berhasil diterima! Status pesanan lunas (PAID), diteruskan ke dapur, dan stok bahan baku otomatis dipotong.",
		"data":    res,
	})
}

func (h *CashierHandler) UpdateStatusPesanan(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	var req model.UpdateStatusPesananRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pilihan status_pesanan tidak valid ('menunggu'/'dimasak'/'diantar'/'selesai')"})
		return
	}

	res, err := h.cashierUsecase.UpdateStatusPesanan(id, &req)
	if err != nil {
		if errors.Is(err, model.ErrPesananNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan tidak ditemukan"})
			return
		}
		if errors.Is(err, model.ErrPesananAlreadyVoid) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Pesanan ini sudah dibatalkan dan tidak dapat diubah statusnya"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui status pesanan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status pesanan berhasil diperbarui",
		"data":    res,
	})
}

func (h *CashierHandler) GetStrukPesanan(c *gin.Context) {
	id, ok := parseParamID(c, "id")
	if !ok {
		return
	}

	uangStr := c.DefaultQuery("uang_diterima", "0")
	uangDiterima, _ := strconv.ParseFloat(uangStr, 64)

	struk, err := h.cashierUsecase.GetStrukPesanan(id, uangDiterima)
	if err != nil {
		if errors.Is(err, model.ErrPesananNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mencetak struk pesanan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data struk berhasil digenerate",
		"data":    struk,
	})
}
