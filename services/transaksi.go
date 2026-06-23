package services

import (
	"RPL/database"
	"RPL/models"
	"fmt"
	"time"
)

type UpdateTransaksiRequest struct {
	IdKaryawan       string  `json:"id_karyawan"`
	MetodePembayaran string  `json:"metode_pembayaran"`
	StatusPembayaran string  `json:"status_pembayaran"`
	JumlahDibayar    float64 `json:"jumlah_dibayar"`
}
type FilterTransaksiRequest struct {
	Status string `query:"status"`
	Bulan  int    `query:"bulan"`
	Tahun  int    `query:"tahun"`
}

func GetAllTransaksi() ([]models.Transaksi, error) {
	transaksi := []models.Transaksi{}
	err := database.DB.Select(&transaksi,
		`SELECT 
            t.id_transaksi, t.id_pesanan, t.id_karyawan, t.total_harga,
            t.metode_pembayaran, t.status_pembayaran, t.tgl_transaksi, t.updated_at,
            d.nama_dokter
         FROM transaksi t
         JOIN pesanan p ON t.id_pesanan = p.id_pesanan
         JOIN dokter d ON p.id_dokter = d.id_dokter
         WHERE p.status_pesanan != 'ditolak' 
         ORDER BY t.tgl_transaksi DESC`,
	)
	if err != nil {
		return nil, err
	}
	return transaksi, nil
}

func GetTransaksiByDokter(idDokter string) ([]models.Transaksi, error) {
	transaksi := []models.Transaksi{}
	query := `
		SELECT 
			t.id_transaksi, t.id_pesanan, t.id_karyawan, t.total_harga,
			t.metode_pembayaran, t.status_pembayaran, t.tgl_transaksi, t.updated_at,
			d.nama_dokter
		FROM transaksi t
		JOIN pesanan p ON t.id_pesanan = p.id_pesanan
		JOIN dokter d ON p.id_dokter = d.id_dokter
		WHERE p.id_dokter = $1
		ORDER BY t.tgl_transaksi DESC
	`

	err := database.DB.Select(&transaksi, query, idDokter)
	if err != nil {
		return nil, err
	}
	return transaksi, nil
}

func GetTransaksiByPesanan(idPesanan string) (models.Transaksi, error) {
	var transaksi models.Transaksi
	err := database.DB.Get(&transaksi,
		`SELECT id_transaksi, id_pesanan, id_karyawan, total_harga,
		        metode_pembayaran, status_pembayaran, tgl_transaksi, updated_at
		 FROM transaksi WHERE id_pesanan = $1`,
		idPesanan,
	)
	if err != nil {
		return models.Transaksi{}, err
	}
	return transaksi, nil
}

func GetTransaksiById(idTransaksi string) (models.Transaksi, error) {
	var transaksi models.Transaksi
	err := database.DB.Get(&transaksi,
		`SELECT id_transaksi, id_pesanan, id_karyawan, total_harga,
		        metode_pembayaran, status_pembayaran, tgl_transaksi, updated_at
		 FROM transaksi WHERE id_transaksi = $1`,
		idTransaksi,
	)
	if err != nil {
		return models.Transaksi{}, err
	}
	return transaksi, nil
}

func GetTransaksiBelumBayar() ([]models.Transaksi, error) {
	transaksi := []models.Transaksi{}
	err := database.DB.Select(&transaksi,
		`SELECT id_transaksi, id_pesanan, id_karyawan, total_harga,
		        metode_pembayaran, status_pembayaran, tgl_transaksi, updated_at
		 FROM transaksi
		 WHERE status_pembayaran = 'belum bayar'
		 ORDER BY tgl_transaksi ASC`,
	)
	if err != nil {
		return nil, err
	}
	return transaksi, nil
}

func KonfirmasiPembayaran(idPesanan string, req UpdateTransaksiRequest) error {
	if req.MetodePembayaran == "" {
		return &ValidationError{
			Message: "Metode pembayaran wajib dipilih",
			Field:   "metode_pembayaran",
			Tag:     "required",
		}
	}

	validMetode := map[string]bool{
		"transfer": true,
		"tunai":    true,
		"qris":     true,
		"gopay":    true,
	}
	if !validMetode[req.MetodePembayaran] {
		return &ValidationError{
			Message: "Metode pembayaran tidak valid",
			Field:   "metode_pembayaran",
			Tag:     "invalid",
		}
	}

	var transaksi models.Transaksi
	err := database.DB.Get(&transaksi,
		"SELECT * FROM transaksi WHERE id_pesanan = $1",
		idPesanan,
	)
	if err != nil {
		return &ValidationError{
			Message: "Data transaksi tidak ditemukan",
			Field:   "id_pesanan",
			Tag:     "not_found",
		}
	}

	if transaksi.StatusPembayaran == "lunas" {
		return &ValidationError{
			Message: "Transaksi sudah lunas",
			Field:   "status_pembayaran",
			Tag:     "already_paid",
		}
	}

	statusPembayaran := req.StatusPembayaran
	if req.JumlahDibayar < transaksi.TotalHarga && statusPembayaran == "lunas" {
		statusPembayaran = "belum bayar"
		return &ValidationError{
			Message: "Jumlah dibayar kurang dari total tagihan, status diubah ke belum bayar",
			Field:   "jumlah_dibayar",
			Tag:     "insufficient",
		}
	}

	currentTime := time.Now()

	_, err = database.DB.Exec(
		`UPDATE transaksi
		 SET id_karyawan       = $1,
		     metode_pembayaran = $2,
		     status_pembayaran = $3,
		     updated_at        = $4
		 WHERE id_pesanan = $5`,
		req.IdKaryawan,
		req.MetodePembayaran,
		statusPembayaran,
		currentTime,
		idPesanan,
	)
	if err != nil {
		return err
	}

	return nil
}

func GetTransaksiFiltered(filter FilterTransaksiRequest) ([]models.Transaksi, error) {
	query := `SELECT id_transaksi, id_pesanan, id_karyawan, total_harga,
			         metode_pembayaran, status_pembayaran, tgl_transaksi, updated_at
			  FROM transaksi WHERE 1=1`
	args := []interface{}{}
	argCount := 1

	if filter.Status != "" {
		query += " AND status_pembayaran = $" + itoa(argCount)
		args = append(args, filter.Status)
		argCount++
	}

	if filter.Bulan > 0 {
		query += " AND EXTRACT(MONTH FROM tgl_transaksi) = $" + itoa(argCount)
		args = append(args, filter.Bulan)
		argCount++
	}

	if filter.Tahun > 0 {
		query += " AND EXTRACT(YEAR FROM tgl_transaksi) = $" + itoa(argCount)
		args = append(args, filter.Tahun)
		argCount++
	}

	query += " ORDER BY tgl_transaksi DESC"

	transaksi := []models.Transaksi{}
	err := database.DB.Select(&transaksi, query, args...)
	if err != nil {
		return nil, err
	}
	return transaksi, nil
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}
