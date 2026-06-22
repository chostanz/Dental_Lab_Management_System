package services

import (
	"RPL/models"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
)

func GetAllProduk() ([]models.Produk, error) {
	produk := []models.Produk{}
	rows, errSelect := db.Queryx("SELECT id_produk, nama_bahan, spesifikasi, harga, created_at, updated_at FROM produk")
	if errSelect != nil {
		return nil, errSelect
	}
	for rows.Next() {
		place := models.Produk{}
		rows.StructScan(&place)
		produk = append(produk, place)
	}
	return produk, nil

}

func GetProdukById(id string) (models.Produk, error) {
	var produkId models.Produk

	err := db.Get(&produkId, "SELECT * FROM produk where id_produk = $1", id)
	if err != nil {
		return models.Produk{}, err
	}
	return produkId, nil
}

func AddProduk(req models.Produk) error {
	uuid := uuid.New()
	produkId := uuid.String()

	_, err := db.NamedExec("INSERT INTO produk (id_produk, nama_bahan, spesifikasi, harga) VALUES (:id_produk, :nama_bahan, :spesifikasi, :harga)", map[string]interface{}{
		"id_produk":   produkId,
		"nama_bahan":  req.NamaBahan,
		"spesifikasi": req.Spesifikasi,
		"harga":       req.Harga,
	})
	if err != nil {
		log.Printf("Error saat menambahkan produk: ")
		return err
	}
	return nil
}

func UpdateProduk(req models.Produk, produkId string) (models.Produk, error) {
	currentTime := time.Now()
	_, err := db.NamedExec("UPDATE produk SET nama_bahan = :nama_bahan, spesifikasi = :spesifikasi, harga=:harga, updated_at = :updated_at WHERE id_produk = :produk_id", map[string]interface{}{
		"nama_bahan":  req.NamaBahan,
		"spesifikasi": req.Spesifikasi,
		"harga":       req.Harga,
		"updated_at":  currentTime,
		"produk_id":   produkId,
	})
	if err != nil {
		log.Printf("Error saat mengupdate data produk, ", err)
		return models.Produk{}, err
	}
	return req, nil
}

func DeleteProduk(produkId string) error {
	result, err := db.Exec("DELETE FROM produk WHERE id_produk = $1", produkId)
	if err != nil {
		log.Printf("Error saat hapus data produk: %v ", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
