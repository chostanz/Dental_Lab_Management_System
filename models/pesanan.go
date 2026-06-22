package models

import (
	"time"
)

type Pesanan struct {
	IdPesanan     string    `json:"id_pesanan" db:"id_pesanan"`
	IdDokter      string    `json:"id_dokter" db:"id_dokter"`
	StatusPesanan string    `json:"status_pesanan" db:"status_pesanan"`
	TglPesanan    time.Time `json:"tgl_pesanan"    db:"tgl_pesanan"`
	UpdatedAt     time.Time `json:"updated_at"     db:"updated_at"`
}

type DetailPesanan struct {
	IdDetail        string  `json:"id_detail" db:"id_detail"`
	IdPesanan       string  `json:"id_pesanan" db:"id_pesanan"`
	IdProduk        string  `json:"id_produk" db:"id_produk"`
	KodeGigi        string  `json:"kode_gigi" db:"kode_gigi"`
	Ukuran          string  `json:"ukuran" db:"ukuran"`
	Warna           string  `json:"warna" db:"warna"`
	Jumlah          int     `json:"jumlah" db:"jumlah"`
	HargaSatuan     float64 `json:"harga_satuan" db:"harga_satuan"`
	Subtotal        float64 `json:"subtotal" db:"subtotal"`
	CatatanTambahan string  `json:"catatan_tambahan" db:"catatan_tambahan"`
}

type DetailPesananLengkap struct {
	Pesanan       Pesanan         `json:"pesanan"`
	DetailPesanan []DetailPesanan `json:"detail_pesanan"`
	Transaksi     *Transaksi      `json:"transaksi,omitempty"`
	Pengiriman    *Pengiriman     `json:"pengiriman,omitempty"`
	Revisi        []Revisi        `json:"revisi,omitempty"`
}
