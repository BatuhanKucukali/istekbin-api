package handler

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"net/http"
)

const emptyList = "[]"

func ListHandler(rd *redis.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		u, err := uuid.Parse(c.Param("uuid"))
		if err != nil {
			return echo.ErrNotFound
		}

		reqVal, err := rd.Get(u.String()).Result()
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
