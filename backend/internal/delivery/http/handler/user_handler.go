package handler

import (
	"errors"
	"net/http"

	"tapkasir-pos/internal/model"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase model.UserUsecase
}

func NewUserHandler(r *gin.Engine, us model.UserUsecase) {
	handler := &UserHandler{
		userUsecase: us,
	}

	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/verify-pin", handler.VerifyPIN)
		authGroup.POST("/change-pin", handler.ChangePIN)
	}
}

func (h *UserHandler) VerifyPIN(c *gin.Context) {
	var req model.VerifyPINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data input tidak valid. Pastikan role ('owner'/'manager'/'cashier') dan PIN (4-6 digit) terisi.",
		})
		return
	}

	res, err := h.userUsecase.VerifyPIN(&req)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Role pengguna tidak ditemukan"})
			return
		}
		if errors.Is(err, model.ErrInvalidPIN) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "PIN salah! Silakan coba lagi."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan saat memverifikasi PIN"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Verifikasi PIN berhasil",
		"data":    res,
	})
}

func (h *UserHandler) ChangePIN(c *gin.Context) {
	var req model.ChangePINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Data input tidak valid. Periksa role, old_pin, dan new_pin (4-6 digit).",
		})
		return
	}

	if err := h.userUsecase.ChangePIN(&req); err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Role pengguna tidak ditemukan"})
			return
		}
		if errors.Is(err, model.ErrInvalidPIN) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "PIN lama salah!"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengganti PIN"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "PIN berhasil diperbarui",
	})
}