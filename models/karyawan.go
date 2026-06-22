package models

import "time"

type Karyawan struct {
	IdKaryawan string    `json:"id_karyawan" db:"id_karyawan"`
	Nama       string    `json:"nama" db:"nama"`
	NoHp       string    `json:"no_hp" db:"no_hp"`
	Email      string    `json:"email" db:"email"`
	Password   string    `json:"password" db:"password"`
	Role       string    `json:"role" db:"role"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
type UpdateKaryawan struct {
	Nama  string `json:"nama"`
	Email string `json:"email"`
	NoHp  string `json:"no_hp"`
	Role  string `json:"role"`
}
