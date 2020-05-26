package api

import (
	"github.com/batuhankucukali/istekbin/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateHandler(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/c", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	conf := &config.App{HistoryCount: 5}

	rd := redisClient()
	defer teardown()

	// Assertions
	if assert.NoError(t, CreateBin(conf, rd)(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		ip := c.RealIP()
		items := rd.Get(ip)

		assert.NotEmpty(t, items)
	}
}
