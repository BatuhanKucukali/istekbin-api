package handler

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

	item := Item{Key: uuid.New().String(), CreatedAt: time.Now()}
	item2 := Item{Key: uuid.New().String(), CreatedAt: time.Now()}
	items := []Item{item, item2}

	result, _ := json.Marshal(items)
	rd.Set(c.RealIP(), result, time.Minute*1)

	// Assertions
	if assert.NoError(t, CreateListHandler(rd)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, string(result), rec.Body.String())
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

	result, _ := json.Marshal([]Item{})

	// Assertions
	if assert.NoError(t, CreateListHandler(rd)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, string(result), rec.Body.String())
	}
}
