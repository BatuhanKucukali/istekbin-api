package handler

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"net/http"
	"time"
)

func CreateHandler(r *redis.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		u := uuid.New()
		rPath := fmt.Sprintf("r/%s", u.String()) // TODO add baseUrl

		err := r.Set(u.String(), nil, time.Hour*24).Err()
		if err != nil {
			log.Error("redis set error.", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "request can not created.")
		}

		c.Response().Header().Add("Location", rPath)
		return c.NoContent(http.StatusCreated)
	}
}
