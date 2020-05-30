package api

import (
	"encoding/json"
	"fmt"
	"github.com/batuhankucukali/istekbin-api/internal/config"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCreateRequestShouldReturnNotFoundWhenUuidIsNotValid(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/r/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("uuid")
	c.SetParamValues("abc")

	rd := redisClient()
	defer teardown()

	conf := &config.App{}

	// Assertions
	err := CreateRequest(conf, rd)(c)
	if assert.NotNil(t, err) {
		rec, ok := err.(*echo.HTTPError)
		if ok {
			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	}
}

func TestCreateRequestShouldReturnNotFoundWhenKeyIsNotFound(t *testing.T) {
	// Setup
	key := uuid.New().String()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/r/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("uuid")
	c.SetParamValues(key)

	rd := redisClient()
	defer teardown()

	conf := &config.App{}

	// Assertions
	err := CreateRequest(conf, rd)(c)
	if assert.NotNil(t, err) {
		rec, ok := err.(*echo.HTTPError)
		if ok {
			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	}
}

func TestCreateRequestShouldReturnBadRequestWhenRequestCountExceedMaxRequestCount(t *testing.T) {
	// Setup
	key := uuid.New().String()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/r/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("uuid")
	c.SetParamValues(key)

	rd := redisClient()
	defer teardown()

	conf := &config.App{MaxRequestCount: 2}
	request := Request{Body: "request-1"}
	request2 := Request{Body: "request-2"}

	var rl []Request
	requests := append(rl, request, request2)
	requestBytes, _ := json.Marshal(requests)
	rd.Set(key, requestBytes, time.Minute*1)

	// Assertions
	err := CreateRequest(conf, rd)(c)
	if assert.NotNil(t, err) {
		rec, ok := err.(*echo.HTTPError)
		if ok {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	}
}

func TestCreateRequestShouldCreateRequest(t *testing.T) {
	// Setup
	key := uuid.New().String()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/r/"+key+"/", strings.NewReader("ok"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("User-Agent", "Go-Agent")
	req.Host = "example.com"

	cookie := new(http.Cookie)
	cookie.Name = "username"
	cookie.Value = "jon"
	req.AddCookie(cookie)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("uuid")
	c.SetParamValues(key)

	rd := redisClient()
	defer teardown()

	conf := &config.App{MaxRequestCount: 50}

	rd.Set(key, nil, time.Minute*1)

	// Assertions
	if assert.NoError(t, CreateRequest(conf, rd)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "ok", rec.Body.String())

		val, err := rd.Get(key).Result()
		if err != nil {
			assert.Fail(t, "request not found")
		}

		result := getRequest(val)
		assert.Equal(t, req.Method, result.Method)
		assert.Equal(t, req.Host, result.Host)
		assert.Equal(t, req.UserAgent(), result.UserAgent)
		assert.Equal(t, c.RealIP(), result.Ip)
		assert.Equal(t, echo.MIMEApplicationJSON, result.ContentType)
		assert.Equal(t, convertHeaderToMap(req.Header), result.Headers)
		assert.Equal(t, convertCookieToMap(req), result.Cookies)
	}
}

func TestCreateRequestShouldRemoveForbiddenHeaderFromHeaders(t *testing.T) {
	// Setup
	key := uuid.New().String()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/r/"+key+"/", strings.NewReader("ok"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("X-Forwarded-For", "127.0.0.1")
	req.RemoteAddr = "192.168.1.1"
	req.Host = "example.com"

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("uuid")
	c.SetParamValues(key)

	rd := redisClient()
	defer teardown()

	conf := &config.App{
		ForbiddenHeaders: []string{"X-Forwarded-For", "X-Forwarded-Port", "X-Forwarded-Proto", "X-Request-Start"},
		MaxRequestCount:  50,
	}

	rd.Set(key, nil, time.Minute*1)

	// Assertions
	if assert.NoError(t, CreateRequest(conf, rd)(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "ok", rec.Body.String())

		val, err := rd.Get(key).Result()
		if err != nil {
			assert.Fail(t, "request not found")
		}

		result := getRequest(val)
		assert.Equal(t, map[string]string{"Content-Type": "application/json"}, result.Headers)
	}
}

func getRequest(value string) Request {
	var rl []Request
	err := json.Unmarshal([]byte(value), &rl)
	if err != nil {
		fmt.Println("Unmarshal error")
	}
	return rl[0]
}

func convertHeaderToMap(header http.Header) map[string]string {
	var headerMap = make(map[string]string)
	for key, value := range header {
		headerMap[key] = strings.Join(value, ",")
	}
	return headerMap
}

func convertCookieToMap(req *http.Request) map[string]string {
	var cookieMap = make(map[string]string)
	for _, c := range req.Cookies() {
		cookieMap[c.Name] = c.Value
	}
	return cookieMap
}
