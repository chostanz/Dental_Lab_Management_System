package models

import "time"

type Produk struct {
	Id          string    `json:"id_produk" db:"id_produk"`
	NamaBahan   string    `json:"nama_bahan" db:"nama_bahan"`
	Spesifikasi string    `json:"spesifikasi" db:"spesifikasi"`
	Harga       float64   `json:"harga" db:"harga"`
	CreatedAt   time.Time `json:"created_at"  db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"  db:"updated_at"`
}
