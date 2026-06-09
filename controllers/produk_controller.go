package controllers

import (
	"RPL/models"
	"RPL/services"
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetAllProduk(c echo.Context) error {
	produk, err := services.GetAllProduk()
	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Terjadi kesalahan pada server!",
			"status":  false,
		})
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, produk)
}

func GetProdukById(c echo.Context) error {
	id := c.Param("id")
	var getProduk models.Produk
	getProduk, err := services.GetProdukById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print(err)
			response := models.Response{
				Code:    404,
				Message: "Division tidak ditemukan!",
				Status:  false,
			}
			return c.JSON(http.StatusNotFound, response)
		} else {
			log.Print(err)
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
				Status:  false,
			})
		}
	}
	return c.JSON(http.StatusOK, getProduk)
}

func AddProduk(c echo.Context) error {
	var addProduk models.Produk
	if err := c.Bind(&addProduk); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak valid!",
			Status:  false,
		})
	}
	errVal := c.Validate(&addProduk)
	if errVal == nil {
		addErr := services.AddProduk(addProduk)
		if addErr != nil {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server.",
				Status:  false,
			})
		}
		return c.JSON(http.StatusCreated, &models.Response{
			Code:    201,
			Message: "Berhasil menambahkan produk!",
			Status:  true,
		})
	} else {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Data tidak boleh kosong!",
			Status:  false,
		})
	}
}
func UpdateProduk(c echo.Context) error {
	id := c.Param("id")

	// Cek dulu apakah produk ada
	_, errGet := services.GetProdukById(id)
	if errGet != nil {
		return c.JSON(http.StatusNotFound, &models.Response{
			Code:    404,
			Message: "Gagal mengupdate produk. Produk tidak ditemukan",
			Status:  false,
		})
	}

	// Bind request body
	var editProduk models.Produk
	if err := c.Bind(&editProduk); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data Invalid!",
			Status:  false,
		})
	}

	// Validasi input
	if errValidate := c.Validate(&editProduk); errValidate != nil {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Data tidak boleh kosong!",
			Status:  false,
		})
	}

	// Jalankan update
	_, updateErr := services.UpdateProduk(editProduk, id)
	if updateErr != nil {
		log.Printf("Error saat update produk: %v", updateErr)
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal pada server.",
			Status:  false,
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Code:    200,
		Message: "Produk berhasil diperbarui!",
		Status:  true,
	})
}

func DeleteProduk(c echo.Context) error {
	id := c.Param("id")

	errService := services.DeleteProduk(id)
	if errService != nil {
		return c.JSON(http.StatusNotFound, &models.Response{
			Code:    404,
			Message: "Gagal menghapus role. Role tidak ditemukan!",
			Status:  false,
		})
	}

	return c.JSON(http.StatusOK, &models.Response{
		Code:    200,
		Message: "Produk berhasil dihapus!",
		Status:  true,
	})

}
