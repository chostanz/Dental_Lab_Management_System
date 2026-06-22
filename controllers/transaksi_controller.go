package controllers

import (
	"RPL/services"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetAllTransaksi(c echo.Context) error {
	data, err := services.GetAllTransaksi()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Gagal mengambil data transaksi",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Berhasil mengambil data transaksi",
		"data":    data,
	})
}

func GetTransaksiById(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "ID transaksi tidak boleh kosong",
		})
	}

	data, err := services.GetTransaksiById(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  false,
			"message": "Transaksi tidak ditemukan",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Berhasil mengambil data transaksi",
		"data":    data,
	})
}

func GetTransaksiByDokter(c echo.Context) error {
	idParam := c.Param("id_dokter")

	// Pengecekan Keamanan: Ambil ID dari token JWT
	idContext := c.Get("id_dokter")
	if idContext == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"status":  false,
			"message": "Akses ditolak: Anda belum login.",
		})
	}
	
	idToken := idContext.(string)

	// Validasi: Pastikan ID di URL sama dengan ID di Token
	if idParam != idToken {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"status":  false,
			"message": "Akses ditolak: Anda hanya dapat melihat transaksi Anda sendiri.",
		})
	}

	transaksi, err := services.GetTransaksiByDokter(idParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Gagal mengambil data transaksi",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Berhasil mengambil data transaksi dokter",
		"data":    transaksi,
	})
}

func GetTransaksiByPesanan(c echo.Context) error {
	idPesanan := c.Param("id_pesanan")
	if idPesanan == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "ID pesanan tidak boleh kosong",
		})
	}

	data, err := services.GetTransaksiByPesanan(idPesanan)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  false,
			"message": "Transaksi tidak ditemukan",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Berhasil mengambil data transaksi",
		"data":    data,
	})
}

func GetTransaksiBelumBayar(c echo.Context) error {
	data, err := services.GetTransaksiBelumBayar()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Gagal mengambil data transaksi belum bayar",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Berhasil mengambil data transaksi belum bayar",
		"data":    data,
	})
}

func GetTransaksiFiltered(c echo.Context) error {
	// Ambil query params
	filter := services.FilterTransaksiRequest{}
	filter.Status = c.QueryParam("status")

	if bulan := c.QueryParam("bulan"); bulan != "" {
		b, err := strconv.Atoi(bulan)
		if err == nil {
			filter.Bulan = b
		}
	}

	if tahun := c.QueryParam("tahun"); tahun != "" {
		t, err := strconv.Atoi(tahun)
		if err == nil {
			filter.Tahun = t
		}
	}

	data, err := services.GetTransaksiFiltered(filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Gagal mengambil data transaksi",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Berhasil mengambil data transaksi",
		"data":    data,
	})
}

func KonfirmasiPembayaran(c echo.Context) error {
	idPesanan := c.Param("id_pesanan")
	if idPesanan == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "ID pesanan tidak boleh kosong",
		})
	}

	req := new(services.UpdateTransaksiRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
	}

	if req.IdKaryawan == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  false,
			"message": "ID karyawan tidak boleh kosong",
		})
	}

	err := services.KonfirmasiPembayaran(idPesanan, *req)
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
			"message": "Gagal konfirmasi pembayaran",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  true,
		"message": "Pembayaran berhasil dikonfirmasi",
	})
}
