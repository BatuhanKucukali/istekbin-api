package handler

import (
	"github.com/batuhankucukali/istekbin/config"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreate(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/create/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	conf := &config.App{}

	rd := redisClient()
	defer teardown()

	// assertions
	if assert.NoError(t, CreateHandler(conf, rd)(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}
