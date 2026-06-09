package controllers

import (
	"RPL/services"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func GetAntrianProduksi(c echo.Context) error {
	data, err := services.GetAntrianProduksi()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Gagal mengambil antrian produksi",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Berhasil mengambil antrian produksi",
		"data":    data,
	})
}

func GetAllPengerjaan(c echo.Context) error {
	data, err := services.GetAllPengerjaan()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Gagal mengambil data pengerjaan",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Berhasil mengambil data pengerjaan",
		"data":    data,
	})
}

func GetPengerjaanById(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "ID pengerjaan tidak boleh kosong",
		})
	}

	data, err := services.GetPengerjaanById(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  false,
			"message": "Data pengerjaan tidak ditemukan",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Berhasil mengambil data pengerjaan",
		"data":    data,
	})
}

func MulaiPengerjaan(c echo.Context) error {
	req := new(services.AddPengerjaanRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
	}

	if req.IdPesanan == "" || req.IdKaryawan == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "id_pesanan dan id_karyawan wajib diisi",
		})
	}

	err := services.MulaiPengerjaan(*req)
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
			"message": "Gagal memulai pengerjaan",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status":  true,
		"message": "Pesanan masuk antrian produksi",
	})
}

func UpdateStatusProduksi(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "ID pengerjaan tidak boleh kosong",
		})
	}

	req := new(services.UpdateStatusProduksiRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
	}

	if req.StatusProduksi == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "status_produksi wajib diisi",
		})
	}

	err := services.UpdateStatusProduksi(id, *req)
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
			"message": "Gagal update status produksi",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Status produksi berhasil diupdate",
	})
}

func AjukanRevisi(c echo.Context) error {
	req := new(services.AddRevisiRequest)
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
			"message": "id_pesanan wajib diisi",
		})
	}

	// Cek apakah lebih dari 7 hari (SRS VR-15)
	// Ambil tgl_selesai dari pengerjaan
	var tglSelesai *time.Time
	err := services.GetTglSelesaiPengerjaan(req.IdPesanan, &tglSelesai)

	warning := ""
	if err == nil && tglSelesai != nil {
		selisih := time.Since(*tglSelesai).Hours() / 24
		if selisih > 7 {
			warning = "Revisi diajukan lebih dari 7 hari setelah pesanan selesai, tetap diproses"
		}
	}

	err = services.AjukanRevisi(*req)
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
			"message": "Gagal mengajukan revisi",
			"error":   err.Error(),
		})
	}

	response := map[string]interface{}{
		"status":  true,
		"message": "Revisi berhasil diajukan",
	}
	if warning != "" {
		response["warning"] = warning
	}

	return c.JSON(http.StatusCreated, response)
}

func GetRevisiByPesanan(c echo.Context) error {
	idPesanan := c.Param("id_pesanan")
	if idPesanan == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "ID pesanan tidak boleh kosong",
		})
	}

	data, err := services.GetRevisiByPesanan(idPesanan)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Gagal mengambil data revisi",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Berhasil mengambil data revisi",
		"data":    data,
	})
}
