package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Start(serverPort string) {
	router := echo.New()

	router.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method == http.MethodOptions {
				c.Response().Header().Add("Access-Control-Allow-Origin", "*")
				c.Response().Header().Add("Access-Control-Allow-Methods", "GET, DELETE, POST, PUT, OPTIONS, HEAD")
				c.Response().Header().Add("Access-Control-Allow-Headers", "Authorization, Origin, X-Requested-With, Content-Type, Accept, ngrok-skip-browser-warning")
				c.NoContent(http.StatusOK)
			} else {
				c.Response().Header().Add("Access-Control-Allow-Origin", "*")
				c.Response().Header().Add("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS, HEAD")
				c.Response().Header().Add("Access-Control-Allow-Headers", "Authorization, Origin, X-Requested-With, Content-Type, Accept, ngrok-skip-browser-warning")
			}

			return next(c)
		}
	})

	router.Static("/files", "./data/chunks/")

	err := router.Start(serverPort)
	if err != nil {
		panic(fmt.Errorf("server crashed %s", err))
	}
}
