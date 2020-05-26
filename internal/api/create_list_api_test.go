package api

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateListHandler(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/cl", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	rd := redisClient()
	defer teardown()

	item := Bin{Key: uuid.New().String(), CreatedAt: time.Now()}
	item2 := Bin{Key: uuid.New().String(), CreatedAt: time.Now()}
	items := []Bin{item, item2}

	itemBytes, _ := json.Marshal(items)
	rd.Set(c.RealIP(), itemBytes, time.Minute*1)

	// Assertions
	if assert.NoError(t, ListBin(rd)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, string(itemBytes), rec.Body.String())
	}
}

func TestCreateListHandlerShouldReturnEmptyResponse(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/cl", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	rd := redisClient()
	defer teardown()

	itemBytes, _ := json.Marshal([]Bin{})

	// Assertions
	if assert.NoError(t, ListBin(rd)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, string(itemBytes), rec.Body.String())
	}
}
