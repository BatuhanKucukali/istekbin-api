package handler

import (
	"encoding/json"
	"github.com/batuhankucukali/istekbin/config"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"time"
)

type Item struct {
	Key         string
	CreatedAt   time.Time
}


func CreateHandler(conf *config.App, rd *redis.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		key := uuid.New().String()
		if err := rd.Set(key, nil, conf.RequestStoreTime).Err(); err != nil {
			log.Error("uuid set error.")
			return echo.NewHTTPError(http.StatusInternalServerError, "request can not created.")
		}

		ipAddress := c.RealIP()
		result, _ := rd.Get(ipAddress).Result()

		var items []Item
		json.Unmarshal([]byte(result), &items)
		items = append(items, Item{Key: key, CreatedAt: time.Now()})

		items = deleteItems(items, rd, conf.HistoryCount)

		itemBytes, err := json.Marshal(items)
		if err != nil {
			log.Error("json marshal error.")
			return echo.NewHTTPError(http.StatusInternalServerError, "request can not created.")
		}

		if err := rd.Set(ipAddress, itemBytes, conf.RequestStoreTime).Err(); err != nil {
			log.Error("redis set error.")
			return echo.NewHTTPError(http.StatusInternalServerError, "request can not created.")
		}

		c.Response().Header().Add("Location", key)
		return c.NoContent(http.StatusCreated)
	}
}

func deleteItems(items []Item, rd *redis.Client, limit int) []Item {
	if len(items) > limit {
		rd.Del(items[0].Key)
		return append(items[:0], items[1:]...)
	}
	return items
}