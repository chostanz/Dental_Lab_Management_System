package models

type DashboardStats struct {
	TotalPesanan    int     `json:"total_pesanan"`
	PesananMenunggu int     `json:"pesanan_menunggu"`
	PesananDiproses int     `json:"pesanan_diproses"`
	PesananSelesai  int     `json:"pesanan_selesai"`
	TotalPendapatan float64 `json:"total_pendapatan"`
}
