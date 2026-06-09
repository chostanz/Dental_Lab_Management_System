package services

import (
	"RPL/database"
	"RPL/models"
	"log"
	"time"

	"github.com/google/uuid"
)

type AddPengerjaanRequest struct {
	IdPesanan  string `json:"id_pesanan"`
	IdKaryawan string `json:"id_karyawan"`
}

type UpdateStatusProduksiRequest struct {
	StatusProduksi  string `json:"status_produksi"`
	CatatanKaryawan string `json:"catatan_karyawan"`
	IdRevisi        string `json:"id_revisi"` // opsional
}

type AddRevisiRequest struct {
	IdPesanan       string `json:"id_pesanan"`
	DeskripsiRevisi string `json:"deskripsi_revisi"`
}

func GetAntrianProduksi() ([]models.Pengerjaan, error) {
	pengerjaan := []models.Pengerjaan{}
	err := database.DB.Select(&pengerjaan,
		`SELECT id_pengerjaan, id_pesanan, id_karyawan, id_revisi,
		        status_produksi, catatan_karyawan, tgl_mulai, tgl_selesai,
		        created_at, updated_at
		 FROM pengerjaan
		 WHERE status_produksi IN ('antrian', 'dikerjakan', 'revisi')
		 ORDER BY created_at ASC`,
	)
	if err != nil {
		return nil, err
	}
	return pengerjaan, nil
}

func GetAllPengerjaan() ([]models.Pengerjaan, error) {
	pengerjaan := []models.Pengerjaan{}
	err := database.DB.Select(&pengerjaan,
		`SELECT id_pengerjaan, id_pesanan, id_karyawan, id_revisi,
		        status_produksi, catatan_karyawan, tgl_mulai, tgl_selesai,
		        created_at, updated_at
		 FROM pengerjaan ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	return pengerjaan, nil
}

func GetPengerjaanById(id string) (models.Pengerjaan, error) {
	var pengerjaan models.Pengerjaan
	err := database.DB.Get(&pengerjaan,
		`SELECT id_pengerjaan, id_pesanan, id_karyawan, id_revisi,
		        status_produksi, catatan_karyawan, tgl_mulai, tgl_selesai,
		        created_at, updated_at
		 FROM pengerjaan WHERE id_pengerjaan = $1`,
		id,
	)
	if err != nil {
		return models.Pengerjaan{}, err
	}
	return pengerjaan, nil
}

func MulaiPengerjaan(req AddPengerjaanRequest) error {
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
	if statusPesanan != "disetujui" {
		return &ValidationError{
			Message: "Pesanan belum disetujui bos",
			Field:   "status_pesanan",
			Tag:     "invalid_status",
		}
	}

	var count int
	err = database.DB.Get(&count,
		"SELECT COUNT(*) FROM pengerjaan WHERE id_pesanan = $1 AND status_produksi != 'selesai'",
		req.IdPesanan,
	)
	if err != nil {
		return err
	}
	if count > 0 {
		return &ValidationError{
			Message: "Pesanan sudah ada di antrian produksi",
			Field:   "id_pesanan",
			Tag:     "duplicate",
		}
	}

	pengerjaanId := uuid.New().String()
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

	_, err = tx.Exec(
		`INSERT INTO pengerjaan
		 (id_pengerjaan, id_pesanan, id_karyawan, status_produksi, tgl_mulai, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		pengerjaanId,
		req.IdPesanan,
		req.IdKaryawan,
		"antrian",
		currentTime,
		currentTime,
		currentTime,
	)
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		"UPDATE pesanan SET status_pesanan = $1, updated_at = $2 WHERE id_pesanan = $3",
		"produksi", currentTime, req.IdPesanan,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func UpdateStatusProduksi(pengerjaanId string, req UpdateStatusProduksiRequest) error {
	var currentStatus string
	err := database.DB.Get(&currentStatus,
		"SELECT status_produksi FROM pengerjaan WHERE id_pengerjaan = $1",
		pengerjaanId,
	)
	if err != nil {
		return &ValidationError{
			Message: "Data pengerjaan tidak ditemukan",
			Field:   "id_pengerjaan",
			Tag:     "not_found",
		}
	}

	urutan := map[string]int{
		"antrian":    1,
		"dikerjakan": 2,
		"revisi":     3,
		"selesai":    4,
	}

	if urutan[req.StatusProduksi] < urutan[currentStatus] {
		return &ValidationError{
			Message: "Tidak dapat mengubah ke status sebelumnya",
			Field:   "status_produksi",
			Tag:     "invalid_status",
		}
	}

	if req.StatusProduksi == "revisi" && req.CatatanKaryawan == "" {
		return &ValidationError{
			Message: "Catatan wajib diisi saat status revisi",
			Field:   "catatan_karyawan",
			Tag:     "required",
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

	if req.StatusProduksi == "selesai" {
		_, err = tx.Exec(
			`UPDATE pengerjaan
			 SET status_produksi = $1, catatan_karyawan = $2, tgl_selesai = $3, updated_at = $4
			 WHERE id_pengerjaan = $5`,
			req.StatusProduksi, req.CatatanKaryawan, currentTime, currentTime, pengerjaanId,
		)
	} else {
		_, err = tx.Exec(
			`UPDATE pengerjaan
			 SET status_produksi = $1, catatan_karyawan = $2, updated_at = $3
			 WHERE id_pengerjaan = $4`,
			req.StatusProduksi, req.CatatanKaryawan, currentTime, pengerjaanId,
		)
	}
	if err != nil {
		log.Printf("Error update status produksi: %v", err)
		return err
	}
	var statusPesanan string
	switch req.StatusProduksi {
	case "dikerjakan":
		statusPesanan = "produksi"
	case "selesai":
		statusPesanan = "selesai"
	case "revisi":
		statusPesanan = "revisi"
	default:
		statusPesanan = "produksi"
	}

	var idPesanan string
	err = database.DB.Get(&idPesanan,
		"SELECT id_pesanan FROM pengerjaan WHERE id_pengerjaan = $1",
		pengerjaanId,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"UPDATE pesanan SET status_pesanan = $1, updated_at = $2 WHERE id_pesanan = $3",
		statusPesanan, currentTime, idPesanan,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func AjukanRevisi(req AddRevisiRequest) error {
	if req.DeskripsiRevisi == "" {
		return &ValidationError{
			Message: "Deskripsi revisi wajib diisi",
			Field:   "deskripsi_revisi",
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

	tx, err := database.DB.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	revisiId := uuid.New().String()
	currentTime := time.Now()

	// Insert revisi
	_, err = tx.Exec(
		`INSERT INTO revisi (id_revisi, id_pesanan, status_revisi, deskripsi_revisi, tgl_revisi, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		revisiId, req.IdPesanan, "pending", req.DeskripsiRevisi, currentTime, currentTime,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"UPDATE pesanan SET status_pesanan = $1, updated_at = $2 WHERE id_pesanan = $3",
		"revisi", currentTime, req.IdPesanan,
	)
	if err != nil {
		return err
	}

	// Update pengerjaan yang aktif → revisi + link id_revisi
	_, err = tx.Exec(
		`UPDATE pengerjaan
		 SET status_produksi = $1, id_revisi = $2, updated_at = $3
		 WHERE id_pesanan = $4 AND status_produksi != 'selesai'`,
		"revisi", revisiId, currentTime, req.IdPesanan,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func GetRevisiByPesanan(idPesanan string) ([]models.Revisi, error) {
	revisi := []models.Revisi{}
	err := database.DB.Select(&revisi,
		`SELECT id_revisi, id_pesanan, status_revisi, deskripsi_revisi, tgl_revisi, updated_at
		 FROM revisi WHERE id_pesanan = $1 ORDER BY tgl_revisi DESC`,
		idPesanan,
	)
	if err != nil {
		return nil, err
	}
	return revisi, nil
}

func GetTglSelesaiPengerjaan(idPesanan string, tglSelesai **time.Time) error {
	return database.DB.Get(tglSelesai,
		"SELECT tgl_selesai FROM pengerjaan WHERE id_pesanan = $1 AND status_produksi = 'selesai' LIMIT 1",
		idPesanan,
	)
}
