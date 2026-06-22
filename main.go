package main

import (
	"RPL/routes"
	"RPL/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := routes.Route()
	e.Use(middleware.Logger())
	e.Validator = &utils.CustomValidator{
		Validator: validator.New(),
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		// Izinkan metode HTTP yang diperlukan
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		// PENTING: Izinkan header Authorization agar token JWT dari Axios bisa masuk
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	// hash, _ := bcrypt.GenerateFromPassword(
	// 	[]byte("kucinghitam"),
	// 	bcrypt.DefaultCost,
	// )

	// fmt.Println(
	// 	base64.StdEncoding.EncodeToString(hash),
	// )
	e.Logger.Fatal(e.Start(":8080"))
}
