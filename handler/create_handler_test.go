package handler

import (
	"github.com/batuhankucukali/istekbin/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateHandler(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/create/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	conf := &config.App{}

	rd := redisClient()
	defer teardown()

	// Assertions
	if assert.NoError(t, CreateHandler(conf, rd)(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}
