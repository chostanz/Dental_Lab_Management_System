package main

import (
	"RPL/routes"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := routes.Route()
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
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
