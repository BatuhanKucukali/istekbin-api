package handler

import (
	"github.com/batuhankucukali/istekbin/config"
	"github.com/go-redis/redis/v7"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var rd *redis.Client

func initRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: ":6379",
	})
}

func TestCreate(t *testing.T) {
	// setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/create/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	conf := &config.App{}

	initRedis()

	// assertions
	if assert.NoError(t, CreateHandler(conf, rd)(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}
