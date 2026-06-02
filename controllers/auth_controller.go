package controllers

import (
	"RPL/models"
	"RPL/services"
	"RPL/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

var secretKey = "secretJwToken"

func RegisterDokter(c echo.Context) error {
	var req models.RegisterDokter
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "Data tidak valid!",
			"status":  false,
		})
	}

	if err := services.RegisterDokter(req); err != nil {
		if ve, ok := err.(*services.ValidationError); ok {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code":    400,
				"message": ve.Message,
				"status":  false,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Terjadi kesalahan pada server!",
			"status":  false,
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"code":    201,
		"message": "Registrasi dokter berhasil!",
		"status":  true,
	})
}

func RegisterKaryawan(c echo.Context) error {
	var req models.RegisterKaryawan
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "Data tidak valid!",
			"status":  false,
		})
	}

	plainPassword, err := services.RegisterKaryawan(req)
	if err != nil {
		if ve, ok := err.(*services.ValidationError); ok {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code":    400,
				"message": ve.Message,
				"status":  false,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Terjadi kesalahan pada server!",
			"status":  false,
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"code":    201,
		"message": "Registrasi karyawan berhasil!",
		"status":  true,
		"data": map[string]interface{}{
			"password_sementara": plainPassword, // dikirim sekali, suruh ganti
		},
	})
}

func LoginDokter(c echo.Context) error {
	var req models.LoginDokter
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "Data tidak valid!",
			"status":  false,
		})
	}

	id, nama, err := services.LoginDokter(req)
	if err != nil {
		if ve, ok := err.(*services.ValidationError); ok {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": ve.Message,
				"status":  false,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Terjadi kesalahan pada server!",
			"status":  false,
		})
	}

	token, err := utils.GenerateTokenDokter(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Gagal membuat token!",
			"status":  false,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "Login berhasil!",
		"status":  true,
		"data": map[string]interface{}{
			"token": token,
			"nama":  nama,
			"role":  "dokter",
		},
	})
}

func LoginKaryawan(c echo.Context) error {
	var req models.Login
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "Data tidak valid!",
			"status":  false,
		})
	}

	id, nama, role, err := services.LoginKaryawan(req)
	if err != nil {
		if ve, ok := err.(*services.ValidationError); ok {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": ve.Message,
				"status":  false,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Terjadi kesalahan pada server!",
			"status":  false,
		})
	}

	token, err := utils.GenerateTokenKaryawan(id, role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Gagal membuat token!",
			"status":  false,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "Login berhasil!",
		"status":  true,
		"data": map[string]interface{}{
			"token": token,
			"nama":  nama,
			"role":  role,
		},
	})
}

func ChangePasswordDokter(c echo.Context) error {
	var req models.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "Data tidak valid!",
			"status":  false,
		})
	}

	userID := c.Get("id_dokter").(string)

	if err := services.ChangePasswordDokter(req, userID); err != nil {
		if ve, ok := err.(*services.ValidationError); ok {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code":    400,
				"message": ve.Message,
				"status":  false,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Terjadi kesalahan pada server!",
			"status":  false,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "Password berhasil diubah!",
		"status":  true,
	})
}

func ChangePasswordKaryawan(c echo.Context) error {
	var req models.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "Data tidak valid!",
			"status":  false,
		})
	}

	userID := c.Get("id_karyawan").(string)

	if err := services.ChangePasswordKaryawan(req, userID); err != nil {
		if ve, ok := err.(*services.ValidationError); ok {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"code":    400,
				"message": ve.Message,
				"status":  false,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Terjadi kesalahan pada server!",
			"status":  false,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "Password berhasil diubah!",
		"status":  true,
	})
}

func ResetPasswordKaryawan(c echo.Context) error {
	karyawanID := c.Param("id")

	plainPassword, err := services.ResetPasswordKaryawan(karyawanID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Gagal mereset password!",
			"status":  false,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "Password berhasil direset!",
		"status":  true,
		"data": map[string]interface{}{
			"password_baru": plainPassword,
		},
	})
}

func Logout(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	utils.InvalidateToken(token)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "Logout berhasil!",
		"status":  true,
	})
}
