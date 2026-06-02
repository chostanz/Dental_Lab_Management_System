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
	karyawan.PUT("/karyawan/change-password", controllers.ChangePasswordKaryawan)

	// Khusus Bos
	bos := e.Group("", middleware.AuthBos)
	bos.POST("/karyawan/register", controllers.RegisterKaryawan)
	bos.PUT("/karyawan/reset-password/:id", controllers.ResetPasswordKaryawan)

	return e
}
