package models

import (
	"time"
)

type Dokter struct {
	IdDokter   string    `json:"id_dokter" db:"id_dokter"`
	NamaDokter string    `json:"nama_dokter" db:"nama_dokter"`
	NoHp       string    `json:"no_hp" db:"no_hp"`
	Email      string    `json:"email" db:"email"`
	Password   string    `json:"password" db:"password"` 
	Alamat     *string   `json:"alamat" db:"alamat"`
	Klinik     *string   `json:"klinik" db:"klinik"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
