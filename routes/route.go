package routes

import (
	"RPL/controllers"
	"RPL/middleware"

	"github.com/labstack/echo/v4"
)

func Route() *echo.Echo {
	e := echo.New()
	e.POST("/register", controllers.RegisterDokter)
	e.POST("/login", controllers.LoginDokter)
	e.POST("/login/karyawan", controllers.LoginKaryawan)

	dokter := e.Group("", middleware.AuthDokter)
	dokter.POST("/logout", controllers.Logout)
	dokter.PUT("/dokter/change-password", controllers.ChangePasswordDokter)

	karyawan := e.Group("", middleware.AuthKaryawan)
	karyawan.POST("/logout", controllers.Logout)
	karyawan.PUT("/karyawan/change-password", controllers.ChangePasswordKaryawan)

	bos := e.Group("", middleware.AuthBos)
	bos.POST("/karyawan/register", controllers.RegisterKaryawan)
	bos.PUT("/karyawan/reset-password/:id", controllers.ResetPasswordKaryawan)

	produk := e.Group("/api/produk")
	produk.GET("", controllers.GetAllProduk)
	produk.GET("/:id", controllers.GetProdukById)

	produkAdmin := e.Group("/api/produk", middleware.AuthCSOrBos)
	produkAdmin.POST("", controllers.AddProduk)
	produkAdmin.PUT("/:id", controllers.UpdateProduk)
	produkAdmin.DELETE("/:id", controllers.DeleteProduk)

	pesananDokter := e.Group("/api/pesanan/dokter", middleware.AuthDokter)
	pesananDokter.POST("", controllers.AddPesanan)
	pesananDokter.GET("/:id_dokter", controllers.GetPesananByDokter)
	pesananDokter.GET("/by/:id", controllers.GetPesananById)
	pesananDokter.GET("/detail/:id", controllers.GetDetailPesanan)
	pesananDokter.GET("/:id/full", controllers.GetPesananLengkap)

	pesananKaryawan := e.Group("/api/pesanan", middleware.AuthKaryawan)
	pesananKaryawan.GET("", controllers.GetAllPesanan)
	pesananKaryawan.GET("/:id", controllers.GetPesananById)
	pesananKaryawan.GET("/:id/detail", controllers.GetDetailPesanan)
	pesananKaryawan.GET("/:id/full", controllers.GetPesananLengkap)

	pesananCS := e.Group("/api/pesanan", middleware.AuthCSOrBos)
	pesananCS.PUT("/:id/status", controllers.UpdateStatusPesanan)
	pesananCS.PUT("/:id/transaksi", controllers.UpdateTransaksi)

	produksi := e.Group("/api/produksi", middleware.AuthKaryawan)
	produksi.GET("", controllers.GetAllPengerjaan)
	produksi.GET("/antrian", controllers.GetAntrianProduksi)
	produksi.GET("/:id", controllers.GetPengerjaanById)

	produksiTeknisi := e.Group("/api/produksi", middleware.AuthTeknisi)
	produksiTeknisi.PUT("/:id/status", controllers.UpdateStatusProduksi)

	revisi := e.Group("/api/revisi", middleware.AuthDokter)
	revisi.POST("", controllers.AjukanRevisi)
	revisi.GET("/:id_pesanan", controllers.GetRevisiByPesanan)

	revisiView := e.Group("/api/revisi", middleware.AuthKaryawan)
	revisiView.GET("/karyawan/:id_pesanan", controllers.GetRevisiByPesanan)

	persetujuan := e.Group("/api/persetujuan", middleware.AuthBos)
	persetujuan.GET("", controllers.GetAllPersetujuan)
	persetujuan.GET("/pending", controllers.GetPesananPending)
	persetujuan.POST("", controllers.BuatKeputusanPersetujuan)

	persetujuanKaryawan := e.Group("/api/persetujuan/all", middleware.AuthKaryawan)
	persetujuanKaryawan.GET("/:id_pesanan", controllers.GetPersetujuanByPesanan)

	pengiriman := e.Group("/api/pengiriman", middleware.AuthCS)
	pengiriman.GET("", controllers.GetAllPengiriman)
	pengiriman.GET("/siap-kirim", controllers.GetPesananSiapKirim)
	pengiriman.GET("/:id/detail", controllers.GetDetailPengiriman)
	pengiriman.GET("/pesanan/:id_pesanan", controllers.GetPengirimanByPesanan)
	pengiriman.POST("", controllers.AddPengiriman)
	pengiriman.PUT("/:id/status", controllers.UpdateStatusPengiriman)

	pengirimanDokter := e.Group("/api/pengiriman", middleware.AuthDokter)
	pengirimanDokter.GET("/pesanan/:id_pesanan", controllers.GetPengirimanByPesananDokter)
	pengirimanDokter.GET("/:id/detail", controllers.GetDetailPengirimanDokter)

	transaksi := e.Group("/api/transaksi", middleware.AuthCSOrBos)
	transaksi.GET("", controllers.GetAllTransaksi)
	transaksi.GET("/:id", controllers.GetTransaksiById)
	transaksi.GET("/belum-bayar", controllers.GetTransaksiBelumBayar)
	transaksi.GET("/pesanan/:id", controllers.GetTransaksiByPesanan)
	transaksi.GET("/filtered", controllers.GetTransaksiFiltered)
	transaksi.PUT("/pesanan/:id_pesanan/konfirmasi", controllers.KonfirmasiPembayaran)

	dashboard := e.Group("/api/dashboard", middleware.AuthKaryawan)
	dashboard.GET("/statistik", controllers.GetDashboardStats)

	return e
}
