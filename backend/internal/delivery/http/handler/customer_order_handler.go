package handler

import (
	"errors"
	"net/http"
	"strconv"

	"tapkasir-pos/internal/model"

	"github.com/gin-gonic/gin"
)

type CustomerOrderHandler struct {
	customerUsecase model.CustomerOrderUsecase
}

func NewCustomerOrderHandler(r *gin.Engine, us model.CustomerOrderUsecase) {
	handler := &CustomerOrderHandler{
		customerUsecase: us,
	}

	// Seluruh rute customer bersifat PUBLIK (Pelanggan scan QR meja dari ponsel)
	custGroup := r.Group("/api/customer")
	{
		custGroup.GET("/menu", handler.FetchMenuTersedia)
		custGroup.POST("/order", handler.CheckoutPesanan)
		custGroup.GET("/order/:id", handler.LacakPesanan)
		custGroup.POST("/order/:id/bayar-qris", handler.SimulasiBayarQRIS)
		custGroup.POST("/order/:id/batal", handler.BatalkanPesanan) // Pelanggan batalkan pesanan sendiri
	}
}

func (h *CustomerOrderHandler) FetchMenuTersedia(c *gin.Context) {
	data, err := h.customerUsecase.FetchMenuTersedia()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar menu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sukses",
		"data":    data,
	})
}

func (h *CustomerOrderHandler) CheckoutPesanan(c *gin.Context) {
	var req model.CustomerCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data pesanan tidak valid. Pastikan nomor_meja, nama_pelanggan (min 2 karakter), nomor_telepon (9-15 digit), metode_pembayaran ('qris'/'tunai'), dan items pesanan terisi dengan benar.",
		})
		return
	}

	res, err := h.customerUsecase.CheckoutPesanan(&req)
	if err != nil {
		if errors.Is(err, model.ErrMenuNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Salah satu menu pesanan tidak ditemukan"})
			return
		}
		if errors.Is(err, model.ErrMenuNotAvailable) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Salah satu menu yang dipesan sedang habis/tidak tersedia"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses pesanan meja"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pesanan meja berhasil dibuat",
		"data":    res,
	})
}

func (h *CustomerOrderHandler) LacakPesanan(c *gin.Context) {
	id, ok := parseCustomerID(c)
	if !ok {
		return
	}

	pesanan, err := h.customerUsecase.LacakPesanan(id)
	if err != nil {
		if errors.Is(err, model.ErrPesananNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal melacak pesanan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sukses",
		"data":    pesanan,
	})
}

func (h *CustomerOrderHandler) SimulasiBayarQRIS(c *gin.Context) {
	id, ok := parseCustomerID(c)
	if !ok {
		return
	}

	pesanan, err := h.customerUsecase.SimulasiBayarQRIS(id)
	if err != nil {
		if errors.Is(err, model.ErrPesananNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan tidak ditemukan"})
			return
		}
		if errors.Is(err, model.ErrPesananAlreadyPaid) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Pesanan sudah berstatus lunas sebelumnya"})
			return
		}
		if errors.Is(err, model.ErrBukanPembayaranQRIS) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Pesanan ini bukan menggunakan metode QRIS"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses verifikasi pembayaran QRIS"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pembayaran QRIS berhasil diverifikasi! Status pesanan lunas (PAID), langsung diteruskan ke dapur (dimasak), dan stok bahan baku otomatis dipotong sesuai resep.",
		"data":    pesanan,
	})
}

// BatalkanPesanan memungkinkan pelanggan membatalkan pesanan sendiri (hanya jika belum diproses/dibayar)
func (h *CustomerOrderHandler) BatalkanPesanan(c *gin.Context) {
	id, ok := parseCustomerID(c)
	if !ok {
		return
	}

	var req model.CustomerBatalPesananRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Alasan pembatalan wajib diisi (minimal 3 karakter)"})
		return
	}

	err := h.customerUsecase.BatalkanPesananOlehPelanggan(id, &req)
	if err != nil {
		if errors.Is(err, model.ErrPesananNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan tidak ditemukan"})
			return
		}
		if errors.Is(err, model.ErrPesananAlreadyPaid) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Pesanan sudah dibayar, tidak bisa dibatalkan sendiri. Hubungi kasir."})
			return
		}
		if errors.Is(err, model.ErrPesananAlreadyVoid) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Pesanan ini sudah dibatalkan sebelumnya"})
			return
		}
		if errors.Is(err, model.ErrPesananInvalidStatus) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Pesanan sudah diproses dapur, tidak bisa dibatalkan. Hubungi kasir atau manager."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membatalkan pesanan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pesanan berhasil dibatalkan. Jika ada pertanyaan, silakan hubungi kasir.",
	})
}

func parseCustomerID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, strconv.IntSize)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pesanan tidak valid"})
		return 0, false
	}
	return uint(id), true
}
