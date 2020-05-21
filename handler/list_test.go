package handler

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	emptyJSON = `[]`
)

func getRequestListJsonString() string {
	r := new(Request)
	r.Uri = "/1"
	r.Method = "POST"

	r2 := new(Request)
	r2.Uri = "/2"
	r2.Method = "GET"

	var rl []Request

	requests := append([]Request{*r, *r2}, rl...)
	p, err := json.Marshal(requests)
	if err != nil {
		fmt.Println("json convertion error.")
	}

	return string(p)
}

func TestListShouldReturnEmtpyList(t *testing.T) {
	// Setup
	key := uuid.New().String()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/l", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("uuid")
	c.SetParamValues(key)

	rd := redisClient()
	defer teardown()

	// assertions
	if assert.NoError(t, ListHandler(rd)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, emptyJSON, rec.Body.String())
	}
}

func TestListShouldReturnList(t *testing.T) {
	// Setup
	key := uuid.New().String()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/l", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("uuid")
	c.SetParamValues(key)

	listJsonString := getRequestListJsonString()

	rd := redisClient()
	defer teardown()
	rd.Set(key, listJsonString, time.Minute*1)

	// assertions
	if assert.NoError(t, ListHandler(rd)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, listJsonString, rec.Body.String())
	}
}
