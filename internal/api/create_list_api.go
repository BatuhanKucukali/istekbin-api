package api

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

// ListRequest of Bin
// @Summary ListRequest of created bin
// @Produce json
// @Success 200 {object} []handler.Bin
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /cl [get]
func ListBin(rd *redis.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		key := c.RealIP() // TODO list key should be ip+fingerprint

		result, err := rd.Get(key).Result()
		if err == redis.Nil {
			return c.JSON(http.StatusOK, []Bin{})
		} else if err != nil {
			log.Errorf("redis error. %s", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "redis error")
		}

		var items []Bin
		if err := json.Unmarshal([]byte(result), &items); err != nil {
			log.Errorf("deserialize error. %s", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "deserialize error")
		}

		return c.JSON(http.StatusOK, items)
	}
}
