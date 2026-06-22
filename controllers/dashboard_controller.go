package controllers

import (
	"RPL/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetDashboardStats(c echo.Context) error {
	stats, err := services.GetDashboardStats()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Gagal mengambil data statistik dashboard",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Berhasil mengambil data statistik",
		"data":    stats,
	})
}
