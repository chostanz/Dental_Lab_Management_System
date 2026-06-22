package services

import (
	"RPL/database"
	"RPL/models"
)

func GetDokterProfile(idDokter string) (models.Dokter, error) {
	var dokter models.Dokter
	err := database.DB.Get(&dokter, "SELECT id_dokter, nama_dokter, email, no_hp, klinik, alamat FROM dokter WHERE id_dokter = $1", idDokter)
	return dokter, err
}

func UpdateDokterProfile(idDokter string, nama, noHp, klinik, alamat string) error {
	_, err := database.DB.Exec(
		"UPDATE dokter SET nama_dokter = $1, no_hp = $2, klinik = $3, alamat = $4, updated_at = NOW() WHERE id_dokter = $5",
		nama, noHp, klinik, alamat, idDokter,
	)
	return err
}

func GetKaryawanProfile(idKaryawan string) (models.Karyawan, error) {
	var karyawan models.Karyawan
	err := database.DB.Get(&karyawan, "SELECT id_karyawan, nama, email, no_hp, role FROM karyawan WHERE id_karyawan = $1", idKaryawan)
	return karyawan, err
}

func UpdateKaryawanProfile(idKaryawan, nama, noHp string) error {
	_, err := database.DB.Exec(
		"UPDATE karyawan SET nama = $1, no_hp = $2, updated_at = NOW() WHERE id_karyawan = $3",
		nama, noHp, idKaryawan,
	)
	return err
}
