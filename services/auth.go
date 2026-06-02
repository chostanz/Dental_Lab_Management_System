package services

import (
	"RPL/database"
	"RPL/models"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math/big"

	jose "github.com/dvsekhvalnov/jose2go"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

var db *sqlx.DB = database.Connection()

type ValidationError struct {
	Message string
	Field   string
	Tag     string
}

func (ve *ValidationError) Error() string {
	return ve.Message
}

func GenerateRandomPassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%"
	password := make([]byte, length)
	for i := range password {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[n.Int64()]
	}
	return string(password), nil
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(hashed), nil
}

func verifyPassword(hashedBase64, plain string) error {
	decoded, err := base64.StdEncoding.DecodeString(hashedBase64)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword(decoded, []byte(plain))
}

type JwtCustomClaims struct {
	KaryawanID         string `json:"id_karyawan"`
	Role               string `json:"role"`
	jwt.StandardClaims        // Embed the StandardClaims struct

}

func DecryptJWE(jweToken string, secretKey string) (string, error) {
	// Dekripsi token JWE
	decrypted, _, err := jose.Decode(jweToken, secretKey)
	if err != nil {
		return "", err
	}
	return decrypted, nil
}

func RegisterDokter(dokterRegister models.RegisterDokter) error {
	if len(dokterRegister.Password) < 8 {
		return &ValidationError{
			Message: "Password should be of 8 characters long ",
			Field:   "password",
			Tag:     "strong_password",
		}
	}
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM dokter WHERE email = $1", dokterRegister.Email)
	if err != nil {
		return err
	}
	if count > 0 {
		return &ValidationError{
			Message: "Email sudah digunakan",
			Field:   "email",
			Tag:     "duplicate",
		}
	}

	hashedPassword, err := hashPassword(dokterRegister.Password)
	if err != nil {
		return err
	}

	id := uuid.New().String()

	_, errInsert := db.NamedExec("INSERT INTO dokter (id_dokter, nama_dokter, email, password, no_hp, alamat, klinik) VALUES (:id_dokter, :nama_dokter, :email, :password, :no_hp, :alamat, :klinik)", map[string]interface{}{
		"id_dokter":   id,
		"nama_dokter": dokterRegister.Nama,
		"email":       dokterRegister.Email,
		"password":    hashedPassword,
		"no_hp":       dokterRegister.NoHp,
		"alamat":      dokterRegister.Alamat,
		"klinik":      dokterRegister.Klinik,
	})
	if errInsert != nil {
		log.Print(errInsert)
		return errInsert
	}
	return nil
}

func RegisterKaryawan(karyawanRegister models.RegisterKaryawan) (plainPassword string, err error) {
	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM karyawan WHERE email = $1", karyawanRegister.Email)
	if err != nil {
		return "", err
	}
	if count > 0 {
		return "", &ValidationError{
			Message: "Email sudah digunakan",
			Field:   "email",
			Tag:     "duplicate",
		}
	}
	validRoles := map[string]bool{"cs": true, "teknisi": true, "bos": true}
	if !validRoles[karyawanRegister.Role] {
		return "", &ValidationError{
			Message: "Role tidak valid, harus cs/teknisi/bos",
			Field:   "role",
			Tag:     "invalid_role",
		}
	}

	// Generate password random
	plainPassword, err = GenerateRandomPassword(10)
	if err != nil {
		return "", err
	}

	hashedPassword, err := hashPassword(plainPassword)
	if err != nil {
		return "", err
	}

	id := uuid.New().String()

	_, errInsert := db.NamedExec("INSERT INTO karyawan (id_karyawan, nama, email, password, no_hp, role) VALUES (:id_karyawan, :nama, :email, :password, :no_hp, :role)", map[string]interface{}{
		"id_karyawan": id,
		"nama":        karyawanRegister.Nama,
		"email":       karyawanRegister.Email,
		"password":    hashedPassword,
		"no_hp":       karyawanRegister.NoHp,
		"role":        karyawanRegister.Role,
	})
	if errInsert != nil {
		log.Print(errInsert)
		return "", errInsert
	}
	return plainPassword, nil
}

func LoginDokter(dokterLogin models.LoginDokter) (id string, nama string, err error) {
	var dbPassword string
	var namaDokter string
	var idDokter string

	err = db.QueryRow("SELECT id_dokter, nama_dokter, password from dokter where email = $1", dokterLogin.Email).Scan(
		&idDokter, &namaDokter, &dbPassword)
	if err != nil {
		fmt.Println("Error querying users:", err)
		return "", "", err
	}
	if err = verifyPassword(dbPassword, dokterLogin.Password); err != nil {
		return "", "", &ValidationError{
			Message: "Email atau password salah",
			Field:   "password",
			Tag:     "invalid_credentials",
		}
	}
	return idDokter, namaDokter, nil
}

func LoginKaryawan(req models.Login) (id string, nama string, role string, err error) {
	var dbPassword string
	var namaKaryawan string
	var idKaryawan string
	var dbRole string

	err = db.QueryRow(
		`SELECT id_karyawan, nama, password, role FROM karyawan WHERE email = $1`,
		req.Email,
	).Scan(&idKaryawan, &namaKaryawan, &dbPassword, &dbRole)

	if err != nil {
		return "", "", "", &ValidationError{
			Message: "Email atau password salah",
			Field:   "email",
			Tag:     "invalid_credentials",
		}
	}

	if err = verifyPassword(dbPassword, req.Password); err != nil {
		return "", "", "", &ValidationError{
			Message: "Email atau password salah",
			Field:   "password",
			Tag:     "invalid_credentials",
		}
	}

	return idKaryawan, namaKaryawan, dbRole, nil
}

func ChangePasswordDokter(changePassword models.ChangePasswordRequest, userID string) error {
	if len(changePassword.NewPassword) < 8 {
		return &ValidationError{
			Message: "Password should be of 8 characters long",
			Field:   "password",
			Tag:     "strong_password",
		}
	}

	if changePassword.OldPassword == changePassword.NewPassword {
		return errors.New("new password must be different from old password")
	}

	var dbPassword string
	err := db.Get(&dbPassword, "SELECT password FROM dokter WHERE id_dokter = $1", userID)
	if err != nil {
		return err
	}

	if err = verifyPassword(dbPassword, changePassword.OldPassword); err != nil {
		return &ValidationError{
			Message: "Password lama salah",
			Field:   "old_password",
			Tag:     "invalid_credentials",
		}
	}

	hashedPassword, err := hashPassword(changePassword.NewPassword)
	if err != nil {
		return err
	}

	_, err = database.DB.Exec(
		"UPDATE dokter SET password = $1, updated_at = NOW() WHERE id_dokter = $2",
		hashedPassword, userID,
	)
	return err
}

func ChangePasswordKaryawan(changePassword models.ChangePasswordRequest, userID string) error {
	if len(changePassword.NewPassword) < 8 {
		return &ValidationError{
			Message: "Password should be of 8 characters long",
			Field:   "password",
			Tag:     "strong_password",
		}
	}
	if changePassword.OldPassword == changePassword.NewPassword {
		return errors.New("password baru harus berbeda dari password lama")
	}

	var dbPassword string
	err := database.DB.Get(&dbPassword,
		"SELECT password FROM karyawan WHERE id_karyawan = $1", userID, // fix: query ke karyawan
	)
	if err != nil {
		return err
	}

	if err = verifyPassword(dbPassword, changePassword.OldPassword); err != nil {
		return &ValidationError{
			Message: "Password lama salah",
			Field:   "old_password",
			Tag:     "invalid_credentials",
		}
	}

	hashedPassword, err := hashPassword(changePassword.NewPassword)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		"UPDATE karyawan SET password = $1, updated_at = NOW() WHERE id_karyawan = $2",
		hashedPassword, userID,
	)
	return err
}

func ResetPasswordKaryawan(karyawanID string) (plainPassword string, err error) {
	plainPassword, err = GenerateRandomPassword(10)
	if err != nil {
		return "", err
	}

	hashedPassword, err := hashPassword(plainPassword)
	if err != nil {
		return "", err
	}

	_, err = database.DB.Exec(
		"UPDATE karyawan SET password = $1, updated_at = NOW() WHERE id_karyawan = $2",
		hashedPassword, karyawanID,
	)
	if err != nil {
		return "", err
	}

	fmt.Printf("Password reset untuk karyawan %s\n", karyawanID)
	return plainPassword, nil
}
