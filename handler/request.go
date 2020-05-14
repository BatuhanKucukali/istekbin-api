package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	Method      string            `json:"method"`
	Host        string            `json:"host"`
	UserAgent   string            `json:"user_agent"`
	Ip          string            `json:"ip"`
	Uri         string            `json:"uri"`
	ContentType string            `json:"content_type"`
	Header      map[string]string `json:"header"`
	Body        string            `json:"body"`
	CreatedAt   time.Time         `json:"created_at"`
}

const ForbiddenContentType = "multipart/form-data"

func RequestHandler(rd *redis.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		u, err := uuid.Parse(c.Param("uuid"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid path.")
		}

		_, err = rd.Get(u.String()).Result()
		if err == redis.Nil {
			return echo.NewHTTPError(http.StatusBadRequest, "request not found.")
		}

		req := c.Request()
		contentType := req.Header.Get("Content-Type")

		if isSupportedContentType(contentType) { // TODO it must be support all content type
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("%s not supported yet.", ForbiddenContentType))
		}

		r := new(Request)
		r.Method = req.Method
		r.Host = req.Host
		r.Uri = req.RequestURI
		r.ContentType = contentType
		r.UserAgent = req.UserAgent()
		r.Ip = req.RemoteAddr
		r.CreatedAt = time.Now()
		r.Header = getHeader(req.Header)

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "body can not reading.")
		}
		r.Body = string(body)

		val, err := rd.Get(u.String()).Result()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "request not found.")
		}

		var rl []Request

		if len(val) > 0 {
			err := json.Unmarshal([]byte(val), &rl)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "request can not deserialized.")
			}
		}

		requests := append(rl, *r)

		v, err := json.Marshal(requests)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "request can not serialized.")
		}
		rd.Set(u.String(), v, time.Hour*24)

		return c.String(http.StatusOK, "ok")
	}
}

func getHeader(header http.Header) map[string]string {
	var headerMap = make(map[string]string)
	for key, value := range header {
		headerMap[key] = strings.Join(value, ",")
	}
	return headerMap
}

func isSupportedContentType(contentType string) bool {
	return strings.Contains(contentType, ForbiddenContentType)
}
