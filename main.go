package main

import (
	"fmt"
	"github.com/batuhankucukali/binrequest/config"
	"github.com/batuhankucukali/binrequest/handler"
	"github.com/go-redis/redis/v7"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
	"log"
	"strings"
)

func main() {
	initViper()

	var conf config.Config

	err := viper.Unmarshal(&conf)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Pre(middleware.AddTrailingSlash())
	e.Use(middleware.BodyLimit(conf.AppConfig.BodyLimit))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  []string{conf.AppConfig.ClientUrl},
		AllowHeaders:  []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderLocation},
		ExposeHeaders: []string{echo.HeaderLocation},
	}))

	rd := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.RedisConfig.Host, conf.RedisConfig.Port),
		Password: conf.RedisConfig.Password,
		DB:       conf.RedisConfig.DB,
	})

	if err := rd.Ping().Err(); err != nil {
		log.Fatal("Redis connection error.", err)
	}

	e.GET("/", handler.HomeHandler)
	e.POST("/c/", handler.CreateHandler(&conf.AppConfig, rd))
	e.GET("/l/:uuid/", handler.ListHandler(rd))
	e.Any("/r/:uuid/*", handler.RequestHandler(&conf.AppConfig, rd))

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", conf.AppConfig.Port)))
}

func initViper() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigFile("config.yml")
	viper.SetConfigType("yml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Config file not found...")
	}
}
