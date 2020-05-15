package handler

import (
	"fmt"
	"github.com/batuhankucukali/binrequest/config"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"net/http"
)

func CreateHandler(conf *config.App, rd *redis.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		u := uuid.New()

		if err := rd.Set(u.String(), nil, conf.RequestStoreTime).Err(); err != nil {
			log.Error("redis set error.")
			return echo.NewHTTPError(http.StatusInternalServerError, "request can not created.")
		}

		c.Response().Header().Add("Location", fmt.Sprintf("%s/r/%s", conf.BaseUrl, u.String()))
		return c.NoContent(http.StatusCreated)
	}
}
