package services

import (
	"RPL/database"
	"RPL/models"
)

func GetDashboardStats() (models.DashboardStats, error) {
	var stats models.DashboardStats

	// Hitung Total Semua Pesanan
	err := database.DB.QueryRow("SELECT COUNT(*) FROM pesanan WHERE status_pesanan != 'ditolak'").Scan(&stats.TotalPesanan)
	if err != nil {
		return stats, err
	}

	// Hitung Pesanan Menunggu (Pending)
	err = database.DB.QueryRow("SELECT COUNT(*) FROM pesanan WHERE status_pesanan = 'menunggu'").Scan(&stats.PesananMenunggu)
	if err != nil {
		return stats, err
	}

	// Hitung Pesanan Diproses (Produksi)
	err = database.DB.QueryRow("SELECT COUNT(*) FROM pesanan WHERE status_pesanan = 'produksi'").Scan(&stats.PesananDiproses)
	if err != nil {
		return stats, err
	}

	// Hitung Pesanan Selesai
	err = database.DB.QueryRow("SELECT COUNT(*) FROM pesanan WHERE status_pesanan = 'selesai'").Scan(&stats.PesananSelesai)
	if err != nil {
		return stats, err
	}

	// Hitung Total Pendapatan (Hanya dari transaksi yang status_pembayaran = 'lunas')
	err = database.DB.QueryRow(`SELECT COALESCE(SUM(t.total_harga), 0) 
    FROM transaksi t
    JOIN pesanan p ON t.id_pesanan = p.id_pesanan
    WHERE t.status_pembayaran = 'lunas' 
      AND p.status_pesanan != 'ditolak'`).Scan(&stats.TotalPendapatan)
	if err != nil {
		return stats, err
	}

	return stats, nil
}
