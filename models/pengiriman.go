package models

import "time"

type Pengiriman struct {
	IdPengiriman string     `json:"id_pengiriman" db:"id_pengiriman"`
	IdPesanan    string     `json:"id_pesanan" db:"id_pesanan"`
	NamaJasa     string     `json:"nama_jasa" db:"nama_jasa"`
	NoResi       string     `json:"no_resi" db:"no_resi"`
	TglKirim     *time.Time `json:"tgl_kirim" db:"tgl_kirim"`
	TglDiterima  *time.Time `json:"tgl_diterima" db:"tgl_diterima"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`

	StatusTerakhir *string   `json:"status_terakhir" db:"status_terakhir"`
}

type DetailPengiriman struct {
	IdDetailPengiriman string    `json:"id_detail_pengiriman" db:"id_detail_pengiriman"`
	IdPengiriman       string    `json:"id_pengiriman" db:"id_pengiriman"`
	Status             string    `json:"status" db:"status"`
	Waktu              time.Time `json:"waktu" db:"waktu"`
	Keterangan         string    `json:"keterangan" db:"keterangan"`
}
