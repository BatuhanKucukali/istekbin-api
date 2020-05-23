package handler

import (
	"github.com/batuhankucukali/istekbin/config"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

func CreateHandler(conf *config.App, rd *redis.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		key := uuid.New().String()

		if err := rd.Set(key, nil, conf.RequestStoreTime).Err(); err != nil {
			log.Error("redis set error.")
			return echo.NewHTTPError(http.StatusInternalServerError, "request can not created.")
		}

		c.Response().Header().Add("Location", key)
		return c.NoContent(http.StatusCreated)
	}
}
