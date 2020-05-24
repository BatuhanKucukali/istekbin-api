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

		req := c.Request()
		r := new(Request)
		r.Ip = getIp(conf, req)
		r.Cookies = getCookies(req) // maybe we change redis key ip to ip-cookie

		result, _ := rd.Get(r.Ip).Result()

		var items []Item
		json.Unmarshal([]byte(result), &items)
		items = append(items, Item{Key: key, CreatedAt: time.Now()})

		items = deleteItems(items, rd)

		itemBytes, err := json.Marshal(items)
		if err != nil {
			log.Error("json marshal error.")
			return echo.NewHTTPError(http.StatusInternalServerError, "request can not created.")
		}

		if err := rd.Set(r.Ip, itemBytes, conf.RequestStoreTime).Err(); err != nil {
			log.Error("redis set error.")
			return echo.NewHTTPError(http.StatusInternalServerError, "request can not created.")
		}

		c.Response().Header().Add("Location", key)
		return c.NoContent(http.StatusCreated)
	}
}

func deleteItems(items []Item, rd *redis.Client) []Item {
	if len(items) > 5 {
		rd.Del(items[0].Key)
		return append(items[:0], items[1:]...)
	}
	return items
}