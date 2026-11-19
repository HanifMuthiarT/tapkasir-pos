package handler

import (
	"errors"
	"net/http"
	"strconv"

	"tapkasir-pos/internal/delivery/http/middleware"
	"tapkasir-pos/internal/model"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	reportUsecase model.ReportUsecase
}

func NewReportHandler(r *gin.Engine, us model.ReportUsecase, jwtSecret string) {
	handler := &ReportHandler{
		reportUsecase: us,
	}

	// Seluruh laporan hanya bisa diakses oleh Owner atau Manager
	reportGroup := r.Group("/api/laporan")
	reportGroup.Use(middleware.RequireRole(jwtSecret, "owner", "manager"))
	{
		// GET /api/laporan?start=2026-11-01&end=2026-11-17
		reportGroup.GET("", handler.GetLaporanLengkap)
		// GET /api/laporan/export?start=2026-11-01&end=2026-11-17
		reportGroup.GET("/export", handler.ExportExcel)
	}
}

func (h *ReportHandler) GetLaporanLengkap(c *gin.Context) {
	startStr := c.DefaultQuery("start", "")
	endStr := c.DefaultQuery("end", "")

	laporan, err := h.reportUsecase.GetLaporanLengkap(startStr, endStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data laporan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sukses",
		"data":    laporan,
	})
}

func (h *ReportHandler) ExportExcel(c *gin.Context) {
	startStr := c.DefaultQuery("start", "")
	endStr := c.DefaultQuery("end", "")

	fileBytes, err := h.reportUsecase.ExportExcel(startStr, endStr)
	if err != nil {
		if errors.Is(err, model.ErrUnauthorized) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengekspor laporan ke Excel"})
		return
	}

	startLabel := startStr
	endLabel := endStr
	if startLabel == "" {
		startLabel = "awal"
	}
	if endLabel == "" {
		endLabel = "sekarang"
	}
	filename := "Laporan_TapKasir_" + startLabel + "_sd_" + endLabel + ".xlsx"

	// Kirim file sebagai response download langsung di browser
	c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Length", strconv.Itoa(len(fileBytes)))
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileBytes)
}
