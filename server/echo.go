package server

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func newEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.HEAD("/", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"service":   "apigateway",
			"version":   "v1",
			"host":      c.Request().Host,
			"timestamp": time.Now().UTC().String(),
		})
	})

	return e
}
