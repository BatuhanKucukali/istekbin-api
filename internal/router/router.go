package router

import (
	"github.com/batuhankucukali/istekbin/internal/api"
	"github.com/batuhankucukali/istekbin/internal/config"
	"github.com/go-redis/redis/v7"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Init(conf *config.Configuration, rd *redis.Client) *echo.Echo {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit(conf.AppConfig.BodyLimit))
	e.Use(RateLimit(conf.Rate, rd))
	e.Use(cors(&conf.AppConfig))

	e.GET("/", api.Home)
	e.POST("/c", api.CreateBin(&conf.AppConfig, rd))
	e.GET("/cl", api.ListBin(rd))
	e.GET("/l/:uuid", api.ListRequest(rd))

	rg := e.Group("/r/:uuid")
	rg.Any("", api.CreateRequest(&conf.AppConfig, rd))
	rg.Any("/", api.CreateRequest(&conf.AppConfig, rd))
	rg.Any("/*", api.CreateRequest(&conf.AppConfig, rd))

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	return e
}
