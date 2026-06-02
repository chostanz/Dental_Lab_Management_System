package models

type Login struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
	Role     string `json:"role" db:"role"`
}

type LoginDokter struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type RegisterDokter struct {
	Id       string `json:"id_dokter" db:"id_dokter"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
	NoHp     string `json:"no_hp" db:"no_hp"`
	Nama     string `json:"nama" db:"nama"`
	Klinik   string `json:"klinik" db:"klinik"`
	Alamat   string `json:"alamat" db:"alamat"`
}

type RegisterKaryawan struct {
	Id       string `json:"id_karyawan" db:"id_karyawan"`
	Nama     string `json:"nama" db:"nama"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
	NoHp     string `json:"no_hp" db:"no_hp"`
	Role     string `json:"role" db:"role"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}
