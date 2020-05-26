package api

import (
	"encoding/json"
	"github.com/batuhankucukali/istekbin/internal/config"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"time"
)

type BinItem struct {
	Key       string    `json:"key"`
	CreatedAt time.Time `json:"created_at"`
}

// Create Istekbin
// @Summary Create istekbin
// @Accept  json
// @Header 201 {string} Location "uuid"
// @Failure 500 {object} echo.HTTPError
// @Router /c [post]
func Create(conf *config.App, rd *redis.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		key := uuid.New().String()

		if err := rd.Set(key, nil, conf.RequestStoreTime).Err(); err != nil {
			log.Errorf("redis error. %s", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "redis error")
		}

		go saveList(key, c.RealIP(), *conf, *rd) // TODO list key should be ip+fingerprint

		c.Response().Header().Add("Location", key)
		return c.NoContent(http.StatusCreated)
	}
}

func saveList(key string, ipAddress string, conf config.App, rd redis.Client) {
	result, _ := rd.Get(ipAddress).Result()

	var items []BinItem
	if err := json.Unmarshal([]byte(result), &items); err != nil {
		log.Errorf("deserialize error. %s", err)
	}

	items = append([]BinItem{{Key: key, CreatedAt: time.Now()}}, items...)
	items = deleteItemsIfNeeded(items, conf.HistoryCount)

	itemBytes, err := json.Marshal(items)
	if err != nil {
		log.Errorf("serialize error. %s", err)
	}

	if err := rd.Set(ipAddress, itemBytes, conf.RequestStoreTime).Err(); err != nil {
		log.Errorf("redis error. %s", err)
	}
}

func deleteItemsIfNeeded(items []BinItem, limit int) []BinItem {
	if len(items) > limit {
		return items[0:limit]
	}
	return items
}
