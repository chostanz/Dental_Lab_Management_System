package controllers

import (
	"RPL/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetAllPersetujuan(c echo.Context) error {
	data, err := services.GetAllPersetujuan()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Gagal mengambil data persetujuan",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Berhasil mengambil data persetujuan",
		"data":    data,
	})
}

func GetPesananPending(c echo.Context) error {
	data, err := services.GetPesananPending()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Gagal mengambil pesanan pending",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Berhasil mengambil pesanan pending",
		"data":    data,
	})
}

func GetPersetujuanByPesanan(c echo.Context) error {
	idPesanan := c.Param("id_pesanan")
	if idPesanan == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "ID pesanan tidak boleh kosong",
		})
	}

	data, err := services.GetPersetujuanByPesanan(idPesanan)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  false,
			"message": "Data persetujuan tidak ditemukan",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Berhasil mengambil data persetujuan",
		"data":    data,
	})
}

func BuatKeputusanPersetujuan(c echo.Context) error {
	req := new(services.KeputusanPersetujuanRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
	}

	// Validasi field wajib
	if req.IdPesanan == "" || req.IdKaryawan == "" || req.Status == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "id_pesanan, id_karyawan, dan status wajib diisi",
		})
	}

	err := services.BuatKeputusanPersetujuan(*req)
	if err != nil {
		if valErr, ok := err.(*services.ValidationError); ok {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"status":  false,
				"message": valErr.Message,
				"field":   valErr.Field,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Gagal membuat keputusan persetujuan",
			"error":   err.Error(),
		})
	}

	// Response message sesuai keputusan
	message := map[string]string{
		"disetujui": "Pesanan disetujui dan masuk antrian produksi",
		"ditolak":   "Pesanan ditolak",
		"revisi":    "Pesanan dikembalikan untuk revisi",
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status":  true,
		"message": message[req.Status],
	})
}
