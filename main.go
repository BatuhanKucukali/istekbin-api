package main

import (
	"github.com/batuhankucukali/binrequest/handler"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", handler.HomeHandler)
	e.POST("/c", handler.CreateHandler)
	e.Any("/r/:uuid/*", handler.RequestHandler)

	e.Logger.Fatal(e.Start(":1323"))
}
