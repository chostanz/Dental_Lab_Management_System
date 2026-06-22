package controllers

import (
	"RPL/models"
	"RPL/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetAllKaryawan(c echo.Context) error {
	list, err := services.GetAllKaryawan()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Gagal mengambil data karyawan!",
			"status":  false,
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "Berhasil mengambil data karyawan!",
		"status":  true,
		"data":    list,
	})
}

func GetKaryawanByID(c echo.Context) error {
	id := c.Param("id")
	k, err := services.GetKaryawanByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"code":    404,
			"message": "Karyawan tidak ditemukan!",
			"status":  false,
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "Berhasil mengambil detail karyawan!",
		"status":  true,
		"data":    k,
	})
}

func UpdateKaryawan(c echo.Context) error {
	id := c.Param("id")

	var req models.UpdateKaryawan
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "Data tidak valid!",
			"status":  false,
		})
	}

	if err := services.UpdateKaryawan(id, req); err != nil {
		if ve, ok := err.(*services.ValidationError); ok {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code":    400,
				"message": ve.Message,
				"status":  false,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Gagal memperbarui data karyawan!",
			"status":  false,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "Data karyawan berhasil diperbarui!",
		"status":  true,
	})
}

func DeleteKaryawan(c echo.Context) error {
	id := c.Param("id")
	if err := services.DeleteKaryawan(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Gagal menghapus karyawan!",
			"status":  false,
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "Karyawan berhasil dihapus!",
		"status":  true,
	})
}
