package handler

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

func CreateListHandler(rd *redis.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		key := c.RealIP()

		result, err := rd.Get(key).Result()
		if err == redis.Nil {
			return c.JSON(http.StatusOK, []Item{})
		} else if err != nil {
			log.Errorf("redis error. %s", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "redis error")
		}

		var items []Item
		if err := json.Unmarshal([]byte(result), &items); err != nil {
			log.Errorf("deserialize error. %s", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "deserialize error")
		}

		return c.JSON(http.StatusOK, items)
	}
}
