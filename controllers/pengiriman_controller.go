package controllers

import (
	"RPL/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetAllPengiriman(c echo.Context) error {
	data, err := services.GetAllPengiriman()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Gagal mengambil data pengiriman",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Berhasil mengambil data pengiriman",
		"data":    data,
	})
}

func GetPengirimanByPesanan(c echo.Context) error {
	idPesanan := c.Param("id_pesanan")
	if idPesanan == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "ID pesanan tidak boleh kosong",
		})
	}

	data, err := services.GetPengirimanByPesanan(idPesanan)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  false,
			"message": "Data pengiriman tidak ditemukan",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Berhasil mengambil data pengiriman",
		"data":    data,
	})
}

func GetDetailPengiriman(c echo.Context) error {
	idPengiriman := c.Param("id")
	if idPengiriman == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "ID pengiriman tidak boleh kosong",
		})
	}

	data, err := services.GetDetailPengiriman(idPengiriman)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Gagal mengambil detail pengiriman",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Berhasil mengambil detail pengiriman",
		"data":    data,
	})
}

func GetPesananSiapKirim(c echo.Context) error {
	data, err := services.GetPesananSiapKirim()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Gagal mengambil pesanan siap kirim",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Berhasil mengambil pesanan siap kirim",
		"data":    data,
	})
}

func AddPengiriman(c echo.Context) error {
	req := new(services.AddPengirimanRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
	}

	if req.IdPesanan == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "ID pesanan tidak boleh kosong",
		})
	}

	err := services.AddPengiriman(*req)
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
			"message": "Gagal menambahkan pengiriman",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status":  true,
		"message": "Pengiriman berhasil ditambahkan",
	})
}

func UpdateStatusPengiriman(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "ID pengiriman tidak boleh kosong",
		})
	}

	req := new(services.UpdateStatusPengirimanRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
	}

	if req.Status == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "Status tidak boleh kosong",
		})
	}

	err := services.UpdateStatusPengiriman(id, *req)
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
			"message": "Gagal update status pengiriman",
			"error":   err.Error(),
		})
	}

	message := "Status pengiriman berhasil diupdate"
	if req.Status == "Diterima" {
		message = "Pesanan telah diterima, status pesanan selesai"
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": message,
	})
}


// controllers/pengiriman.go - tambahkan ini

func GetPengirimanByPesananDokter(c echo.Context) error {
	idPesanan := c.Param("id_pesanan")
	if idPesanan == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "ID pesanan tidak boleh kosong",
		})
	}

	// Ambil id_dokter dari token (di-set oleh middleware AuthDokter)
	idDokter := c.Get("id_dokter").(string)

	// Validasi kepemilikan
	if err := services.ValidatePengirimanOwnership(idPesanan, idDokter); err != nil {
		if valErr, ok := err.(*services.ValidationError); ok {
			status := http.StatusBadRequest
			if valErr.Tag == "forbidden" {
				status = http.StatusForbidden
			}
			return c.JSON(status, map[string]interface{}{
				"status":  false,
				"message": valErr.Message,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Gagal validasi akses",
			"error":   err.Error(),
		})
	}

	data, err := services.GetPengirimanByPesanan(idPesanan)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  false,
			"message": "Data pengiriman tidak ditemukan",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Berhasil mengambil data pengiriman",
		"data":    data,
	})
}

func GetDetailPengirimanDokter(c echo.Context) error {
	idPengiriman := c.Param("id")
	if idPengiriman == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "ID pengiriman tidak boleh kosong",
		})
	}

	idDokter := c.Get("id_dokter").(string)

	if err := services.ValidatePengirimanOwnershipByID(idPengiriman, idDokter); err != nil {
		if valErr, ok := err.(*services.ValidationError); ok {
			status := http.StatusBadRequest
			if valErr.Tag == "forbidden" {
				status = http.StatusForbidden
			}
			return c.JSON(status, map[string]interface{}{
				"status":  false,
				"message": valErr.Message,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Gagal validasi akses",
			"error":   err.Error(),
		})
	}

	data, err := services.GetDetailPengiriman(idPengiriman)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Gagal mengambil detail pengiriman",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Berhasil mengambil detail pengiriman",
		"data":    data,
	})
}