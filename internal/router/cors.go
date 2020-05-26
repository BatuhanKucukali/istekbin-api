package router

import (
	"github.com/batuhankucukali/istekbin-api/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func cors(conf *config.App) echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  []string{conf.ClientUrl},
		AllowHeaders:  []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderLocation},
		ExposeHeaders: []string{echo.HeaderLocation},
	})
}
