package controllers

import (
	"RPL/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetAllPesanan(c echo.Context) error {
	pesanan, err := services.GetAllPesanan()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Gagal mengambil data pesanan",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Berhasil mengambil data pesanan",
		"data":    pesanan,
	})
}

func GetPesananById(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "ID pesanan tidak boleh kosong",
		})
	}

	pesanan, err := services.GetPesananById(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  false,
			"message": "Pesanan tidak ditemukan",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Berhasil mengambil data pesanan",
		"data":    pesanan,
	})
}

func GetPesananByDokter(c echo.Context) error {
	idDokter := c.Param("id_dokter")
	if idDokter == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "ID dokter tidak boleh kosong",
		})
	}

	pesanan, err := services.GetPesananByDokter(idDokter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Gagal mengambil data pesanan",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Berhasil mengambil data pesanan dokter",
		"data":    pesanan,
	})
}

func GetDetailPesanan(c echo.Context) error {
	idPesanan := c.Param("id")
	if idPesanan == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "ID pesanan tidak boleh kosong",
		})
	}

	details, err := services.GetDetailPesanan(idPesanan)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Gagal mengambil detail pesanan",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Berhasil mengambil detail pesanan",
		"data":    details,
	})
}
func AddPesanan(c echo.Context) error {
	// 1. AMBIL ID DOKTER DENGAN AMAN (Safe Type Assertion)
	idContext := c.Get("id_dokter")
	if idContext == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"status":  false,
			"message": "Akses Ditolak: Anda belum login atau token tidak valid (Bukan Dokter).",
		})
	}

	// Jika aman, baru ubah ke string
	idDokter := idContext.(string)

	// 2. Deklarasikan variabel request (Bukan pointer)
	var req services.AddPesananRequest

	// 3. Bind JSON body dari frontend ke dalam struct req
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
	}

	// 4. SUNTIKKAN ID DOKTER DARI TOKEN
	req.IdDokter = idDokter

	// 5. Panggil service SATU KALI saja
	err := services.AddPesanan(req)

	// 6. Tangani error jika terjadi kegagalan di service
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
			"message": "Gagal menambahkan pesanan",
			"error":   err.Error(),
		})
	}

	// 7. Jika sukses
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status":  true,
		"message": "Pesanan berhasil dibuat, menunggu persetujuan bos",
	})
}

func UpdateStatusPesanan(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "ID pesanan tidak boleh kosong",
		})
	}

	// Ambil status dari body
	body := new(struct {
		Status string `json:"status_pesanan"`
	})
	if err := c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
	}

	if body.Status == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "Status tidak boleh kosong",
		})
	}

	err := services.UpdateStatusPesanan(id, body.Status)
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
			"message": "Gagal mengupdate status pesanan",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Status pesanan berhasil diupdate",
	})
}

func UpdateTransaksi(c echo.Context) error {
	idPesanan := c.Param("id")
	if idPesanan == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "ID pesanan tidak boleh kosong",
		})
	}

	body := new(struct {
		IdKaryawan string `json:"id_karyawan"`
		Metode     string `json:"metode_pembayaran"`
		Status     string `json:"status_pembayaran"`
	})
	if err := c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
	}

	// Validasi field wajib
	if body.IdKaryawan == "" || body.Metode == "" || body.Status == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "id_karyawan, metode_pembayaran, dan status_pembayaran wajib diisi",
		})
	}

	err := services.UpdateTransaksi(idPesanan, body.IdKaryawan, body.Metode, body.Status)
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
			"message": "Gagal mengupdate transaksi",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Transaksi berhasil diupdate",
	})
}


func GetPesananLengkap(c echo.Context) error {

	id := c.Param("id")

	data, err := services.GetPesananLengkap(id)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status": false,
			"message": "Pesanan tidak ditemukan",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": true,
		"message": "Berhasil mengambil detail pesanan lengkap",
		"data": data,
	})
}