package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	Id         string `json:"id"`          // generic, bisa dokter atau karyawan
	IdKaryawan string `json:"id_karyawan"` // Tambahan untuk token karyawan
	IdDokter   string `json:"id_dokter"`   // Tambahan untuk token dokter
	Role       string `json:"role"`        // dokter = "", cs/teknisi/bos = rolenya
	jwt.StandardClaims
}

func getSecretKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secretJwToken"
	}
	return []byte(secret)
}

func extractToken(c echo.Context) (string, bool) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", false
	}
	return strings.TrimPrefix(authHeader, "Bearer "), true
}

func parseToken(tokenStr string) (*JwtCustomClaims, bool) {
	claims := &JwtCustomClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return getSecretKey(), nil
	})
	if err != nil || !token.Valid {
		return nil, false
	}
	return claims, true
}

func AuthAny(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenStr, ok := extractToken(c)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak ditemukan atau tidak valid!",
				"status":  false,
			})
		}

		claims, ok := parseToken(tokenStr)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak valid atau sudah expired!",
				"status":  false,
			})
		}

		if claims.Id == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak valid!",
				"status":  false,
			})
		}

		c.Set("id", claims.Id)
		c.Set("role", claims.Role)
		return next(c)
	}
}

func AuthDokter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenStr, ok := extractToken(c)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak ditemukan atau tidak valid!",
				"status":  false,
			})
		}

		claims, ok := parseToken(tokenStr)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak valid atau sudah expired!",
				"status":  false,
			})
		}
		dokterID := claims.Id
		if dokterID == "" {
			dokterID = claims.IdDokter
		}
		if dokterID == "" || claims.Role != "" {
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"code":    403,
				"message": "Akses ditolak! Hanya untuk Dokter.",
				"status":  false,
			})
		}

		c.Set("id_dokter", dokterID)
		return next(c)
	}
}

func AuthKaryawan(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenStr, ok := extractToken(c)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak ditemukan atau tidak valid!",
				"status":  false,
			})
		}

		claims, ok := parseToken(tokenStr)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak valid atau sudah expired!",
				"status":  false,
			})
		}
		karyawanID := claims.Id
		if karyawanID == "" {
			karyawanID = claims.IdKaryawan
		}

		// karyawan harus punya id dan role
		if karyawanID == "" || claims.Role == "" {
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"code":    403,
				"message": "Akses ditolak! Hanya untuk Karyawan.",
				"status":  false,
			})
		}

		c.Set("id_karyawan", karyawanID)
		c.Set("role", claims.Role)
		return next(c)
	}
}

func AuthCS(next echo.HandlerFunc) echo.HandlerFunc {
	return AuthKaryawan(func(c echo.Context) error {
		if c.Get("role") != "cs" {
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"code":    403,
				"message": "Akses ditolak! Hanya untuk CS.",
				"status":  false,
			})
		}
		return next(c)
	})
}

func AuthTeknisi(next echo.HandlerFunc) echo.HandlerFunc {
	return AuthKaryawan(func(c echo.Context) error {
		if c.Get("role") != "teknisi" {
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"code":    403,
				"message": "Akses ditolak! Hanya untuk Teknisi.",
				"status":  false,
			})
		}
		return next(c)
	})
}

func AuthBos(next echo.HandlerFunc) echo.HandlerFunc {
	return AuthKaryawan(func(c echo.Context) error {
		if c.Get("role") != "bos" {
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"code":    403,
				"message": "Akses ditolak! Hanya untuk Bos.",
				"status":  false,
			})
		}
		return next(c)
	})
}

func AuthCSOrBos(next echo.HandlerFunc) echo.HandlerFunc {
	return AuthKaryawan(func(c echo.Context) error {
		role := c.Get("role")
		if role != "cs" && role != "bos" {
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"code":    403,
				"message": "Akses ditolak! Hanya untuk CS atau Bos.",
				"status":  false,
			})
		}
		return next(c)
	})
}

// AuthCSOrTeknisi untuk halaman produksi yang bisa diakses CS dan Teknisi
func AuthCSOrTeknisi(next echo.HandlerFunc) echo.HandlerFunc {
	return AuthKaryawan(func(c echo.Context) error {
		role := c.Get("role")
		if role != "cs" && role != "teknisi" {
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"code":    403,
				"message": "Akses ditolak! Hanya untuk CS atau Teknisi.",
				"status":  false,
			})
		}
		return next(c)
	})
}
