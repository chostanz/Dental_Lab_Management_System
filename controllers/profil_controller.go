package controllers

import (
	"RPL/services"
	"net/http"
	"github.com/labstack/echo/v4"
)

func GetProfileKaryawan(c echo.Context) error {
	idKaryawan := c.Get("id_karyawan").(string)

	profil, err := services.GetKaryawanProfile(idKaryawan)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Gagal memuat profil karyawan",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": true,
		"data":   profil,
	})
}

func UpdateProfileKaryawan(c echo.Context) error {
	idKaryawan := c.Get("id_karyawan").(string)

	type ReqUpdate struct {
		Nama string `json:"nama"`
		NoHp string `json:"no_hp"`
	}
	var req ReqUpdate
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "message": "Invalid request"})
	}

	err := services.UpdateKaryawanProfile(idKaryawan, req.Nama, req.NoHp)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"status": false, "message": "Gagal update profil"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"status": true, "message": "Profil berhasil diupdate"})
}

func GetProfileDokter(c echo.Context) error {
	idDokter := c.Get("id_dokter").(string)

	profil, err := services.GetDokterProfile(idDokter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": "Gagal memuat profil dokter",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": true,
		"data":   profil,
	})
}

func UpdateProfileDokter(c echo.Context) error {
	idDokter := c.Get("id_dokter").(string)

	type ReqUpdateD struct {
		Nama   string `json:"nama_dokter"`
		NoHp   string `json:"no_hp"`
		Klinik string `json:"klinik"`
		Alamat string `json:"alamat"`
	}
	var req ReqUpdateD
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"status": false, "message": "Invalid request"})
	}

	err := services.UpdateDokterProfile(idDokter, req.Nama, req.NoHp, req.Klinik, req.Alamat)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"status": false, "message": "Gagal update profil"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"status": true, "message": "Profil berhasil diupdate"})
}