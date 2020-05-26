package main

import (
	"fmt"
	_ "github.com/batuhankucukali/istekbin/docs"
	"github.com/batuhankucukali/istekbin/internal/api"
	"github.com/batuhankucukali/istekbin/internal/config"
	middleware2 "github.com/batuhankucukali/istekbin/internal/middleware"
	"github.com/go-redis/redis/v7"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
	"strings"
)

// @title Istekbin API
// @description Istekbin is a free service that allows you to collect http request.

// @contact.name API Support
// @contact.url https://github.com/BatuhanKucukali/istekbin-api/issues/new

// @license.name Apache 2.0
// @license.url https://github.com/BatuhanKucukali/istekbin-api/blob/master/LICENSE

// @host api.istekbin.com
// @BasePath /
// @schemes https
func main() {
	initViper()

	var conf config.Config

	err := viper.Unmarshal(&conf)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	e := echo.New()

	rd := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.RedisConfig.Host, conf.RedisConfig.Port),
		Password: conf.RedisConfig.Password,
		DB:       conf.RedisConfig.DB,
	})

	if err := rd.Ping().Err(); err != nil {
		log.Fatal("Redis connection error.", err)
	}

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit(conf.AppConfig.BodyLimit))
	e.Use(middleware2.RateLimit(conf.Rate, rd))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  []string{conf.AppConfig.ClientUrl},
		AllowHeaders:  []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderLocation},
		ExposeHeaders: []string{echo.HeaderLocation},
	}))

	e.GET("/", api.HomeHandler)
	e.POST("/c", api.CreateHandler(&conf.AppConfig, rd))
	e.GET("/cl", api.CreateListHandler(rd))
	e.GET("/l/:uuid", api.ListHandler(rd))

	r := e.Group("/r/:uuid")
	r.Any("", api.RequestHandler(&conf.AppConfig, rd))
	r.Any("/", api.RequestHandler(&conf.AppConfig, rd))
	r.Any("/*", api.RequestHandler(&conf.AppConfig, rd))

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", conf.AppConfig.Port)))
}

func initViper() {
	viper.SetConfigFile("configs/config.yml")
	viper.SetConfigType("yml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Config file not found...")
	}
}
