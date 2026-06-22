package models

import "time"

type Transaksi struct {
	IdTransaksi      string    `json:"id_transaksi" db:"id_transaksi"`
	IdPesanan        string    `json:"id_pesanan" db:"id_pesanan"`
	IdKaryawan       *string   `json:"id_karyawan" db:"id_karyawan"`
	TotalHarga       float64   `json:"total_harga" db:"total_harga"`
	MetodePembayaran *string   `json:"metode_pembayaran" db:"metode_pembayaran"`
	StatusPembayaran string    `json:"status_pembayaran" db:"status_pembayaran"`
	TglTransaksi     time.Time `json:"tgl_transaksi" db:"tgl_transaksi"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
	NamaDokter       *string   `json:"nama_dokter" db:"nama_dokter"`
}
