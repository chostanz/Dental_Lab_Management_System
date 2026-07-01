package main

import (
	"RPL/routes"
	"RPL/utils"
	"net/http"
	"os"

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
		AllowOrigins: []string{
			"http://localhost:5173",
			"https://dentallab.up.railway.app",
		},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
	// hash, _ := bcrypt.GenerateFromPassword(
	// 	[]byte("kucinghitam"),
	// 	bcrypt.DefaultCost,
	// )

	// fmt.Println(
	// 	base64.StdEncoding.EncodeToString(hash),
	// )
	// e.Logger.Fatal(e.Start(":8080"))
}
