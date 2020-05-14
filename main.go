package main

import (
	"github.com/batuhankucukali/binrequest/handler"
	"github.com/go-redis/redis/v7"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"log"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("128K"))

	rd := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := rd.Ping().Result()
	if err != nil {
		log.Fatal("Redis connection error.", err)
	}

	e.GET("/", handler.HomeHandler)
	e.POST("/c", handler.CreateHandler(rd))
	e.Any("/r/:uuid/*", handler.RequestHandler(rd))

	e.Logger.Fatal(e.Start(":1323"))
}
