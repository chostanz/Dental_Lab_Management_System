package models

import (
	"time"
)

type Persetujuan struct {
	IdPersetujuan  string    `json:"id_persetujuan" db:"id_persetujuan"`
	IdPesanan      string    `json:"id_pesanan" db:"id_pesanan"`
	IdKaryawan     string    `json:"id_karyawan" db:"id_karyawan"`
	Status         string    `json:"status" db:"status"`
	CatatanBos     string    `json:"catatan_bos" db:"catatan_bos"`
	TglPersetujuan time.Time `json:"tgl_persetujuan" db:"tgl_persetujuan"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type Pengerjaan struct {
	IdPengerjaan    string     `json:"id_pengerjaan" db:"id_pengerjaan"`
	IdPesanan       string     `json:"id_pesanan" db:"id_pesanan"`
	IdKaryawan      string     `json:"id_karyawan" db:"id_karyawan"`
	IdRevisi        *string    `json:"id_revisi" db:"id_revisi"`
	StatusProduksi  string     `json:"status_produksi" db:"status_produksi"`
	CatatanKaryawan *string    `json:"catatan_karyawan" db:"catatan_karyawan"`
	TglMulai        *time.Time `json:"tgl_mulai" db:"tgl_mulai"`
	TglSelesai      *time.Time `json:"tgl_selesai" db:"tgl_selesai"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

type PengerjaanDetail struct {
	IdPengerjaan    string     `json:"id_pengerjaan" db:"id_pengerjaan"`
	IdPesanan       string     `json:"id_pesanan" db:"id_pesanan"`
	IdKaryawan      string     `json:"id_karyawan" db:"id_karyawan"`
	IdRevisi        *string    `json:"id_revisi" db:"id_revisi"`
	StatusProduksi  string     `json:"status_produksi" db:"status_produksi"`
	CatatanKaryawan *string    `json:"catatan_karyawan" db:"catatan_karyawan"`
	TglMulai        *time.Time `json:"tgl_mulai" db:"tgl_mulai"`
	TglSelesai      *time.Time `json:"tgl_selesai" db:"tgl_selesai"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`

	// Hasil JOIN - info dokter & produk
	NamaDokter string  `json:"nama_dokter" db:"nama_dokter"`
	NamaBahan  *string `json:"nama_bahan" db:"nama_bahan"`
	NamaKaryawan *string `json:"nama" db:"nama"`

	// Hasil JOIN - detail teknis dari detail_pesanan
	KodeGigi        *string `json:"kode_gigi" db:"kode_gigi"`
	Warna           *string `json:"warna" db:"warna"`
	Ukuran          *string `json:"ukuran" db:"ukuran"`
	Jumlah          *int    `json:"jumlah" db:"jumlah"`
	CatatanTambahan *string `json:"catatan_tambahan" db:"catatan_tambahan"` // catatan dari dokter saat pesan
}
type Revisi struct {
	IdRevisi        string    `json:"id_revisi"        db:"id_revisi"`
	IdPesanan       string    `json:"id_pesanan"       db:"id_pesanan"`
	StatusRevisi    string    `json:"status_revisi"    db:"status_revisi"`
	DeskripsiRevisi string    `json:"deskripsi_revisi" db:"deskripsi_revisi"`
	TglRevisi       time.Time `json:"tgl_revisi"       db:"tgl_revisi"`
	UpdatedAt       time.Time `json:"updated_at"       db:"updated_at"`
}
