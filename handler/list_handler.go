package handler

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

// List Istekbin
// @Summary List istekbin
// @Accept  json
// @Produce json
// @Success 200 {object} []handler.Request
// @Param uuid path string true "uuid"
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /l/{uuid} [get]
func ListHandler(rd *redis.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		u, err := uuid.Parse(c.Param("uuid"))
		if err != nil {
			return echo.ErrNotFound
		}

		key := u.String()

		reqVal, err := rd.Get(key).Result()
		if err != nil {
			return echo.ErrNotFound
		}

		if len(reqVal) == 0 {
			reqVal = "[]"
		}

		var rl []Request
		if err := json.Unmarshal([]byte(reqVal), &rl); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "request can not deserialized.")
		}

		return c.JSON(http.StatusOK, rl)
	}
}
