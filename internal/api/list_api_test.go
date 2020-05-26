package api

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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
	requestBytes, err := json.Marshal(requests)
	if err != nil {
		fmt.Println("marshaling error")
	}

	return string(requestBytes)
}

func TestListHandlerShouldReturnNotFoundWhenUuidIsNotValid(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/l", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("uuid")
	c.SetParamValues("abc")

	rd := redisClient()
	defer teardown()

	// Assertions
	err := ListRequest(rd)(c)
	if assert.NotNil(t, err) {
		rec, ok := err.(*echo.HTTPError)
		if ok {
			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	}
}

func TestListHandlerShouldReturnNotFoundWhenKeyIsNotFound(t *testing.T) {
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

	// Assertions
	err := ListRequest(rd)(c)
	if assert.NotNil(t, err) {
		rec, ok := err.(*echo.HTTPError)
		if ok {
			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	}
}

func TestListHandlerShouldReturnEmptyResponse(t *testing.T) {
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

	rd.Set(key, nil, time.Minute*1)

	// Assertions
	if assert.NoError(t, ListRequest(rd)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, "[]", rec.Body.String())
	}
}

func TestListHandlerShouldReturnList(t *testing.T) {
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

	// Assertions
	if assert.NoError(t, ListRequest(rd)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, listJsonString, rec.Body.String())
	}
}
