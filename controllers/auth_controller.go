package controllers

import (
	"RPL/models"
	"RPL/services"
	"RPL/utils"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	jose "github.com/dvsekhvalnov/jose2go"
	"github.com/labstack/echo/v4"
)

var secretKey = "secretJwToken"

func generateToken(id string, role string) (string, error) {
	claims := services.JwtCustomClaims{
		Id:   id,
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	// Enkripsi JWT ke JWE
	encrypted, err := jose.Encrypt(signed, jose.PBES2_HS256_A128KW, jose.A128CBC_HS256, secretKey)
	if err != nil {
		return "", err
	}

	return encrypted, nil
}

// RegisterDokter — FR-01: hanya dokter yang bisa register sendiri
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

// RegisterKaryawan — FR-01: hanya Bos yang bisa tambah karyawan (CS/Teknisi)
// dipasang AuthBos di route
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

// LoginDokter — FR-01
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

	// Role dikosongkan karena ini dokter
	token, err := generateToken(id, "")
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

// LoginKaryawan — FR-01: untuk CS, Teknisi, Bos
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

	token, err := generateToken(id, role)
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

// ChangePasswordDokter — FR-01
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

// ChangePasswordKaryawan — FR-01
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

// ResetPasswordKaryawan — FR-01: hanya Bos yang bisa reset password karyawan
// dipasang AuthBos di route
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

// Logout — invalidate token
func Logout(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	utils.InvalidTokens[token] = true

	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "Logout berhasil!",
		"status":  true,
	})
}