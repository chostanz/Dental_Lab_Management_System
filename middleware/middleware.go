package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	Id   string `json:"id_karyawan"`
	Role string `json:"role"`
	jwt.StandardClaims
}

var secretKey = "secretJwToken"

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
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, false
	}
	return claims, true
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

		if claims.Id == "" || claims.Role != "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token bukan milik dokter!",
				"status":  false,
			})
		}

		c.Set("id_dokter", claims.Id)
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
		fmt.Println("ID:", claims.Id)
		fmt.Println("ROLE:", claims.Role)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak valid atau sudah expired!",
				"status":  false,
			})
		}

		if claims.Id == "" || claims.Role == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token bukan milik karyawan!",
				"status":  false,
			})
		}

		c.Set("id_karyawan", claims.Id)
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
