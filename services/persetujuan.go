// services/persetujuan.go
package services

import (
	"RPL/database"
	"RPL/models"
	"time"

	"github.com/google/uuid"
)

type KeputusanPersetujuanRequest struct {
	IdPesanan  string `json:"id_pesanan"`
	IdKaryawan string `json:"id_karyawan"`
	Status     string `json:"status"`
	CatatanBos string `json:"catatan_bos"`
}

func GetAllPersetujuan() ([]models.Persetujuan, error) {
	persetujuan := []models.Persetujuan{}
	err := database.DB.Select(&persetujuan,
		`SELECT id_persetujuan, id_pesanan, id_karyawan, status, catatan_bos, tgl_persetujuan, updated_at
		 FROM persetujuan ORDER BY tgl_persetujuan DESC`,
	)
	if err != nil {
		return nil, err
	}
	return persetujuan, nil
}

func GetPesananPending() ([]models.Pesanan, error) {
	pesanan := []models.Pesanan{}
	err := database.DB.Select(&pesanan,
		`SELECT id_pesanan, id_dokter, status_pesanan, tgl_pesanan, updated_at
		 FROM pesanan
		 WHERE status_pesanan = 'menunggu'
		 ORDER BY tgl_pesanan ASC`,
	)
	if err != nil {
		return nil, err
	}
	return pesanan, nil
}

func GetPersetujuanByPesanan(idPesanan string) (models.Persetujuan, error) {
	var persetujuan models.Persetujuan
	err := database.DB.Get(&persetujuan,
		`SELECT id_persetujuan, id_pesanan, id_karyawan, status, catatan_bos, tgl_persetujuan, updated_at
		 FROM persetujuan WHERE id_pesanan = $1`,
		idPesanan,
	)
	if err != nil {
		return models.Persetujuan{}, err
	}
	return persetujuan, nil
}

func BuatKeputusanPersetujuan(req KeputusanPersetujuanRequest) error {
	validStatus := map[string]bool{
		"disetujui": true,
		"ditolak":   true,
		"revisi":    true,
	}
	if !validStatus[req.Status] {
		return &ValidationError{
			Message: "Status tidak valid, harus disetujui/ditolak/revisi",
			Field:   "status",
			Tag:     "invalid_status",
		}
	}

	if req.Status == "ditolak" && req.CatatanBos == "" {
		return &ValidationError{
			Message: "Catatan wajib diisi untuk penolakan",
			Field:   "catatan_bos",
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

	var count int
	err = database.DB.Get(&count,
		"SELECT COUNT(*) FROM persetujuan WHERE id_pesanan = $1",
		req.IdPesanan,
	)
	if err != nil {
		return err
	}
	if count > 0 {
		return &ValidationError{
			Message: "Pesanan sudah memiliki keputusan",
			Field:   "id_pesanan",
			Tag:     "duplicate",
		}
	}

	if statusPesanan != "menunggu" {
		return &ValidationError{
			Message: "Pesanan sudah diproses sebelumnya",
			Field:   "status_pesanan",
			Tag:     "invalid_status",
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

	persetujuanId := uuid.New().String()
	_, err = tx.Exec(
		`INSERT INTO persetujuan
		 (id_persetujuan, id_pesanan, id_karyawan, status, catatan_bos, tgl_persetujuan, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		persetujuanId,
		req.IdPesanan,
		req.IdKaryawan,
		req.Status,
		req.CatatanBos,
		currentTime,
		currentTime,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"UPDATE pesanan SET status_pesanan = $1, updated_at = $2 WHERE id_pesanan = $3",
		req.Status, currentTime, req.IdPesanan,
	)
	if err != nil {
		return err
	}

	if req.Status == "disetujui" {
		pengerjaanId := uuid.New().String()
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
	}

	return tx.Commit()
}
