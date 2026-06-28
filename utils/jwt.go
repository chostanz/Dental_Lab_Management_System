package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func getSecretKey() string {
	key := os.Getenv("JWT_SECRET")
	if key == "" {
		key = "secretJwToken"
	}
	return key
}

func GenerateTokenDokter(idDokter string) (string, error) {
	claims := jwt.MapClaims{
		"id_dokter": idDokter,
		"exp":       time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(getSecretKey()))
}

func GenerateTokenKaryawan(idKaryawan string, role string) (string, error) {
	claims := jwt.MapClaims{
		"id_karyawan": idKaryawan,
		"role":        role,
		"exp":         time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(getSecretKey()))
}
