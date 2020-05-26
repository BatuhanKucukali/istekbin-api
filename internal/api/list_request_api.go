package api

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

// ListRequest
// @Summary List of created request
// @Accept  json
// @Produce json
// @Success 200 {object} []handler.Request
// @Param uuid path string true "uuid"
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /l/{uuid} [get]
func ListRequest(rd *redis.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		u, err := uuid.Parse(c.Param("uuid"))
		if err != nil {
			return echo.ErrNotFound
		}

		key := u.String()

		result, err := rd.Get(key).Result()
		if redis.Nil == err {
			return echo.ErrNotFound
		} else if err != nil {
			log.Errorf("redis error. %s", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "redis error")
		}

		if len(result) == 0 {
			return c.JSON(http.StatusOK, []Request{})
		}

		var rl []Request
		if err := json.Unmarshal([]byte(result), &rl); err != nil {
			log.Errorf("deserialize error. %s", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "deserialize error")
		}

		return c.JSON(http.StatusOK, rl)
	}
}
