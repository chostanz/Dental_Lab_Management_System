package services

import (
	"RPL/database"
	"RPL/models"
	"log"
	"time"

	"github.com/google/uuid"
)

type AddPesananRequest struct {
	IdDokter      string                 `json:"id_dokter"`
	DetailPesanan []DetailPesananRequest `json:"detail_pesanan"`
}
type DetailPesananRequest struct {
	IdProduk        string  `json:"id_produk"`
	KodeGigi        string  `json:"kode_gigi"`
	Ukuran          string  `json:"ukuran"`
	Warna           string  `json:"warna"`
	Jumlah          int     `json:"jumlah"`
	HargaSatuan     float64 `json:"harga_satuan"`
	CatatanTambahan string  `json:"catatan_tambahan"`
}

func GetAllPesanan() ([]models.Pesanan, error) {
	pesanan := []models.Pesanan{}
	rows, errSelect := db.Queryx("SELECT id_pesanan, id_dokter, status_pesanan, tgl_pesanan, updated_at FROM pesanan ORDER BY tgl_pesanan DESC")
	if errSelect != nil {
		return nil, errSelect
	}
	for rows.Next() {
		place := models.Pesanan{}
		rows.StructScan(&place)
		pesanan = append(pesanan, place)
	}
	return pesanan, nil

}

func GetPesananById(id string) (models.Pesanan, error) {
	var pesanan models.Pesanan

	err := db.Get(&pesanan, "SELECT * FROM pesanan where id_pesanan = $1", id)
	if err != nil {
		return models.Pesanan{}, err
	}
	return pesanan, nil
}

func GetPesananByDokter(idDokter string) ([]models.Pesanan, error) {
	pesanan := []models.Pesanan{}
	err := db.Select(&pesanan, "SELECT id_pesanan, id_dokter, status_pengerjaan, tgl_pesanan, updated_at FROM pesanan WHERE id_dokter = $1 ORDER BY tgl_pesanan DESC", idDokter)
	if err != nil {
		return nil, err
	}
	return pesanan, nil
}
func AddPesanan(req AddPesananRequest) error {
	if len(req.DetailPesanan) == 0 {
		return &ValidationError{
			Message: "Pilih minimal 1 produk",
			Field:   "detail_pesanan",
			Tag:     "required",
		}
	}

	tx, err := database.DB.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	pesananId := uuid.New().String()
	currentTime := time.Now()

	_, err = tx.Exec(
		`INSERT INTO pesanan (id_pesanan, id_dokter, status_pesanan, tgl_pesanan, updated_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		pesananId, req.IdDokter, "menunggu", currentTime, currentTime,
	)
	if err != nil {
		return err
	}

	var totalHarga float64
	for _, detail := range req.DetailPesanan {
		if detail.Jumlah <= 0 {
			err = &ValidationError{
				Message: "Jumlah harus lebih dari 0",
				Field:   "jumlah",
				Tag:     "min",
			}
			return err
		}

		detailId := uuid.New().String()
		subtotal := float64(detail.Jumlah) * detail.HargaSatuan
		totalHarga += subtotal

		_, err = tx.Exec(
			`INSERT INTO detail_pesanan 
			 (id_detail, id_pesanan, id_produk, kode_gigi, ukuran, warna, jumlah, harga_satuan, catatan_tambahan)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			detailId, pesananId, detail.IdProduk, detail.KodeGigi,
			detail.Ukuran, detail.Warna, detail.Jumlah,
			detail.HargaSatuan, detail.CatatanTambahan,
		)
		if err != nil {
			return err
		}
	}

	transaksiId := uuid.New().String()
	_, err = tx.Exec(
		`INSERT INTO transaksi 
		 (id_transaksi, id_pesanan, id_karyawan, total_harga, status_pembayaran, tgl_transaksi, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		transaksiId,
		pesananId,
		totalHarga,
		"belum bayar",
		currentTime,
		currentTime,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func UpdateStatusPesanan(pesananId string, status string) error {
	// Validasi status yang boleh
	validStatus := map[string]bool{
		"menunggu":  true,
		"disetujui": true,
		"ditolak":   true,
		"produksi":  true,
		"selesai":   true,
		"revisi":    true,
	}
	if !validStatus[status] {
		return &ValidationError{
			Message: "Status tidak valid",
			Field:   "status_pesanan",
			Tag:     "invalid_status",
		}
	}

	_, err := database.DB.Exec(
		"UPDATE pesanan SET status_pesanan = $1, updated_at = $2 WHERE id_pesanan = $3",
		status,
		time.Now(),
		pesananId,
	)
	if err != nil {
		log.Printf("Error update status pesanan: %v", err)
		return err
	}
	return nil
}

func GetDetailPesanan(idPesanan string) ([]models.DetailPesanan, error) {
	details := []models.DetailPesanan{}
	err := database.DB.Select(&details,
		`SELECT id_detail, id_pesanan, id_produk, kode_gigi, ukuran, warna, 
		        jumlah, harga_satuan, subtotal, catatan_tambahan
		 FROM detail_pesanan WHERE id_pesanan = $1`,
		idPesanan,
	)
	if err != nil {
		return nil, err
	}
	return details, nil
}

func UpdateTransaksi(idPesanan string, idKaryawan string, metode string, status string) error {
	validStatus := map[string]bool{
		"lunas":       true,
		"belum bayar": true,
	}
	if !validStatus[status] {
		return &ValidationError{
			Message: "Status pembayaran tidak valid",
			Field:   "status_pembayaran",
			Tag:     "invalid_status",
		}
	}

	_, err := database.DB.Exec(
		`UPDATE transaksi 
		 SET id_karyawan = $1, metode_pembayaran = $2, status_pembayaran = $3, updated_at = $4
		 WHERE id_pesanan = $5`,
		idKaryawan, metode, status, time.Now(), idPesanan,
	)
	if err != nil {
		log.Printf("Error update transaksi: %v", err)
		return err
	}
	return nil
}
