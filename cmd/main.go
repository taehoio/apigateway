package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	log := logrus.StandardLogger()

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"Hello": "World"})
	})

	e.GET("/call", func(c echo.Context) error {
		u := c.QueryParam("url")
		resp, err := http.DefaultClient.Get(u)
		if resp != nil {
			defer func() {
				if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil {
					logrus.Error(err)
				}
				if err := resp.Body.Close(); err != nil {
					logrus.Error(err)
				}
			}()
		}
		if err != nil {
			return err
		}

		for k, vv := range resp.Header {
			c.Response().Header().Del(k)
			for _, v := range vv {
				c.Response().Header().Add(k, v)
			}
		}
		if err := c.Stream(resp.StatusCode, resp.Header.Get("Content-Type"), resp.Body); err != nil {
			return err
		}

		return nil
	})

	port := 8080
	log.WithField("port", port).Info("server starting...")
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
