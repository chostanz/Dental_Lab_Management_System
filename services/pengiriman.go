package services

import (
	"RPL/database"
	"RPL/models"
	"time"

	"github.com/google/uuid"
)

type AddPengirimanRequest struct {
	IdPesanan string `json:"id_pesanan"`
	NamaJasa  string `json:"nama_jasa"`
	NoResi    string `json:"no_resi"`
}

type UpdateStatusPengirimanRequest struct {
	Status     string `json:"status"`
	Keterangan string `json:"keterangan"`
}

func GetAllPengiriman() ([]models.Pengiriman, error) {
	pengiriman := []models.Pengiriman{}
	err := database.DB.Select(&pengiriman,
		`SELECT id_pengiriman, id_pesanan, nama_jasa, no_resi,
		        tgl_kirim, tgl_diterima, created_at, updated_at
		 FROM pengiriman ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	return pengiriman, nil
}

func GetPengirimanByPesanan(idPesanan string) (models.Pengiriman, error) {
	var pengiriman models.Pengiriman
	err := database.DB.Get(&pengiriman,
		`SELECT id_pengiriman, id_pesanan, nama_jasa, no_resi,
		        tgl_kirim, tgl_diterima, created_at, updated_at
		 FROM pengiriman WHERE id_pesanan = $1`,
		idPesanan,
	)
	if err != nil {
		return models.Pengiriman{}, err
	}
	return pengiriman, nil
}

func GetDetailPengiriman(idPengiriman string) ([]models.DetailPengiriman, error) {
	details := []models.DetailPengiriman{}
	err := database.DB.Select(&details,
		`SELECT id_detail_pengiriman, id_pengiriman, status, waktu, keterangan
		 FROM detail_pengiriman
		 WHERE id_pengiriman = $1
		 ORDER BY waktu ASC`,
		idPengiriman,
	)
	if err != nil {
		return nil, err
	}
	return details, nil
}
func AddPengiriman(req AddPengirimanRequest) error {
	if req.NoResi == "" {
		return &ValidationError{
			Message: "Nomor resi wajib diisi",
			Field:   "no_resi",
			Tag:     "required",
		}
	}
	if req.NamaJasa == "" {
		return &ValidationError{
			Message: "Pilih jasa kurir",
			Field:   "nama_jasa",
			Tag:     "required",
		}
	}

	var statusPesanan string
	err := database.DB.Get(&statusPesanan,
		"SELECT status_pesanan FROM pesanan WHERE id_pesanan = $1",
		req.IdPesanan,
	)
	if err != nil {
		return &ValidationError{
			Message: "Pesanan tidak ditemukan",
			Field:   "id_pesanan",
			Tag:     "not_found",
		}
	}
	if statusPesanan != "selesai" {
		return &ValidationError{
			Message: "Pesanan belum selesai diproduksi",
			Field:   "status_pesanan",
			Tag:     "invalid_status",
		}
	}

	var count int
	err = database.DB.Get(&count,
		"SELECT COUNT(*) FROM pengiriman WHERE id_pesanan = $1",
		req.IdPesanan,
	)
	if err != nil {
		return err
	}
	if count > 0 {
		return &ValidationError{
			Message: "Pengiriman untuk pesanan ini sudah ada",
			Field:   "id_pesanan",
			Tag:     "duplicate",
		}
	}

	currentTime := time.Now()

	tx, err := database.DB.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	pengirimanId := uuid.New().String()
	_, err = tx.Exec(
		`INSERT INTO pengiriman
		 (id_pengiriman, id_pesanan, nama_jasa, no_resi, tgl_kirim, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		pengirimanId,
		req.IdPesanan,
		req.NamaJasa,
		req.NoResi,
		currentTime,
		currentTime,
		currentTime,
	)
	if err != nil {
		return err
	}

	detailId := uuid.New().String()
	_, err = tx.Exec(
		`INSERT INTO detail_pengiriman
		 (id_detail_pengiriman, id_pengiriman, status, waktu, keterangan)
		 VALUES ($1, $2, $3, $4, $5)`,
		detailId,
		pengirimanId,
		"Menunggu",
		currentTime,
		"Pesanan siap dikirim, menunggu kurir",
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"UPDATE pesanan SET status_pesanan = $1, updated_at = $2 WHERE id_pesanan = $3",
		"dikirim", currentTime, req.IdPesanan,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func UpdateStatusPengiriman(idPengiriman string, req UpdateStatusPengirimanRequest) error {
	validStatus := map[string]bool{
		"Menunggu":         true,
		"Dijemput Kurir":   true,
		"Dalam Pengiriman": true,
		"Diterima":         true,
	}
	if !validStatus[req.Status] {
		return &ValidationError{
			Message: "Status tidak valid",
			Field:   "status",
			Tag:     "invalid_status",
		}
	}

	// Cek pengiriman ada
	var idPesanan string
	err := database.DB.Get(&idPesanan,
		"SELECT id_pesanan FROM pengiriman WHERE id_pengiriman = $1",
		idPengiriman,
	)
	if err != nil {
		return &ValidationError{
			Message: "Data pengiriman tidak ditemukan",
			Field:   "id_pengiriman",
			Tag:     "not_found",
		}
	}

	currentTime := time.Now()

	tx, err := database.DB.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	detailId := uuid.New().String()
	_, err = tx.Exec(
		`INSERT INTO detail_pengiriman
		 (id_detail_pengiriman, id_pengiriman, status, waktu, keterangan)
		 VALUES ($1, $2, $3, $4, $5)`,
		detailId,
		idPengiriman,
		req.Status,
		currentTime,
		req.Keterangan,
	)
	if err != nil {
		return err
	}

	if req.Status == "Diterima" {
		_, err = tx.Exec(
			"UPDATE pengiriman SET tgl_diterima = $1, updated_at = $2 WHERE id_pengiriman = $3",
			currentTime, currentTime, idPengiriman,
		)
		if err != nil {
			return err
		}

		_, err = tx.Exec(
			"UPDATE pesanan SET status_pesanan = $1, updated_at = $2 WHERE id_pesanan = $3",
			"selesai", currentTime, idPesanan,
		)
		if err != nil {
			return err
		}
	} else {
		_, err = tx.Exec(
			"UPDATE pengiriman SET updated_at = $1 WHERE id_pengiriman = $2",
			currentTime, idPengiriman,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func GetPesananSiapKirim() ([]models.Pesanan, error) {
	pesanan := []models.Pesanan{}
	err := database.DB.Select(&pesanan,
		`SELECT p.id_pesanan, p.id_dokter, p.status_pesanan, p.tgl_pesanan, p.updated_at
		 FROM pesanan p
		 LEFT JOIN pengiriman pg ON p.id_pesanan = pg.id_pesanan
		 WHERE p.status_pesanan = 'selesai'
		 AND pg.id_pengiriman IS NULL
		 ORDER BY p.updated_at ASC`,
	)
	if err != nil {
		return nil, err
	}
	return pesanan, nil
}
