package routes

import (
	"RPL/controllers"
	"RPL/middleware"

	"github.com/labstack/echo/v4"
)

func Route() *echo.Echo {
	e := echo.New()
	e.POST("/register/dokter", controllers.RegisterDokter)
	e.POST("/login/dokter", controllers.LoginDokter)
	e.POST("/login/karyawan", controllers.LoginKaryawan)

	// Dokter
	dokter := e.Group("", middleware.AuthDokter)
	dokter.POST("/logout", controllers.Logout)
	dokter.PUT("/dokter/change-password", controllers.ChangePasswordDokter)

	// Karyawan (semua role)
	karyawan := e.Group("", middleware.AuthKaryawan)
	karyawan.POST("/logout", controllers.Logout)
	karyawan.PUT("/karyawan/change/password", controllers.ChangePasswordKaryawan)

	//khusus cs dan bos
	bosOrCs := e.Group("", middleware.AuthCSOrBos)
	bosOrCs.POST("/add/produk", controllers.AddProduk)
	bosOrCs.PUT("/produk/edit/:id", controllers.UpdateProduk)
	bosOrCs.DELETE("produk/hapus/:id", controllers.DeleteProduk)

	// Khusus Bos
	bos := e.Group("", middleware.AuthBos)
	bos.POST("/karyawan/register", controllers.RegisterKaryawan)
	bos.PUT("/karyawan/reset/password/:id", controllers.ResetPasswordKaryawan)

	pesanan := e.Group("/api/pesanan")
	pesanan.GET("", controllers.AddPesanan)
	pesanan.GET("/:id", controllers.GetPesananById)
	pesanan.GET("/:id/detail", controllers.GetDetailPesanan)
	pesanan.GET("/dokter/:id_dokter", controllers.GetPesananByDokter)
	pesanan.POST("/pesanan", controllers.AddPesanan)
	pesanan.PUT("/:id/status", controllers.UpdateStatusPesanan)
	pesanan.PUT("/:id/transaksi", controllers.UpdateTransaksi)

	produksi := e.Group("/api/produksi")
	produksi.GET("", controllers.GetAllPengerjaan)
	produksi.GET("/antrian", controllers.GetAntrianProduksi)
	produksi.GET("/:id", controllers.GetPengerjaanById)
	produksi.POST("", controllers.MulaiPengerjaan)
	produksi.PUT("/:id/status", controllers.UpdateStatusProduksi)

	revisi := e.Group("/api/revisi")
	revisi.POST("", controllers.AjukanRevisi)
	revisi.GET("/:id_pesanan", controllers.GetRevisiByPesanan)

	// routes/routes.go
	persetujuan := e.Group("/api/persetujuan")
	persetujuan.GET("", controllers.GetAllPersetujuan)
	persetujuan.GET("/pending", controllers.GetPesananPending)
	persetujuan.GET("/:id_pesanan", controllers.GetPersetujuanByPesanan)
	persetujuan.POST("", controllers.BuatKeputusanPersetujuan)

	// routes/routes.go
	pengiriman := e.Group("/api/pengiriman")
	pengiriman.GET("", controllers.GetAllPengiriman)
	pengiriman.GET("/siap-kirim", controllers.GetPesananSiapKirim)
	pengiriman.GET("/:id/detail", controllers.GetDetailPengiriman)
	pengiriman.GET("/pesanan/:id_pesanan", controllers.GetPengirimanByPesanan)
	pengiriman.POST("", controllers.AddPengiriman)
	pengiriman.PUT("/:id/status", controllers.UpdateStatusPengiriman)
	return e
}
