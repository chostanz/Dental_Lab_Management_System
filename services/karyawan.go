package services

import (
	"RPL/models"
)

func GetAllKaryawan() ([]models.Karyawan, error) {
	var list []models.Karyawan
	err := db.Select(&list, `
		SELECT id_karyawan, nama, email, no_hp, role, created_at, updated_at
		FROM karyawan
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func GetKaryawanByID(id string) (*models.Karyawan, error) {
	var k models.Karyawan
	err := db.Get(&k, `
		SELECT id_karyawan, nama, email, no_hp, role, created_at, updated_at
		FROM karyawan WHERE id_karyawan = $1
	`, id)
	if err != nil {
		return nil, err
	}
	return &k, nil
}

func UpdateKaryawan(id string, req models.UpdateKaryawan) error {
	validRoles := map[string]bool{"cs": true, "teknisi": true, "bos": true}
	if req.Role != "" && !validRoles[req.Role] {
		return &ValidationError{
			Message: "Role tidak valid, harus cs/teknisi/bos",
			Field:   "role",
			Tag:     "invalid_role",
		}
	}

	if req.Email != "" {
		var count int
		err := db.Get(&count, `
			SELECT COUNT(*) FROM karyawan WHERE email = $1 AND id_karyawan != $2
		`, req.Email, id)
		if err != nil {
			return err
		}
		if count > 0 {
			return &ValidationError{
				Message: "Email sudah digunakan karyawan lain",
				Field:   "email",
				Tag:     "duplicate",
			}
		}
	}

	if req.NoHp != "" {
		var count int
		err := db.Get(&count, `
			SELECT COUNT(*) FROM karyawan WHERE no_hp = $1 AND id_karyawan != $2
		`, req.NoHp, id)
		if err != nil {
			return err
		}
		if count > 0 {
			return &ValidationError{
				Message: "Nomor HP sudah digunakan karyawan lain",
				Field:   "no_hp",
				Tag:     "duplicate",
			}
		}
	}

	_, err := db.Exec(`
		UPDATE karyawan
		SET
			nama       = COALESCE(NULLIF($1, ''), nama),
			email      = COALESCE(NULLIF($2, ''), email),
			no_hp      = COALESCE(NULLIF($3, ''), no_hp),
			role       = COALESCE(NULLIF($4, ''), role),
			updated_at = NOW()
		WHERE id_karyawan = $5
	`, req.Nama, req.Email, req.NoHp, req.Role, id)
	return err
}

func DeleteKaryawan(id string) error {
	_, err := db.Exec(`DELETE FROM karyawan WHERE id_karyawan = $1`, id)
	return err
}