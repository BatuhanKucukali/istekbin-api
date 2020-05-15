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

func RequestHandler(rd *redis.Client) func(c echo.Context) error {
	return func(c echo.Context) error {
		u, err := uuid.Parse(c.Param("uuid"))
		if err != nil {
			return echo.ErrNotFound
		}

		if err := rd.Get(u.String()).Err(); err != redis.Nil {
			return echo.ErrNotFound
		}

		req := c.Request()
		contentType := req.Header.Get("Content-Type")

		if isSupportedContentType(contentType) { // TODO it must be support all content type
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("%s not supported.", echo.MIMEMultipartForm))
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
			if err := json.Unmarshal([]byte(val), &rl); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "request can not deserialized.")
			}
		}

		requests := append(rl, *r)

		if err := set(rd, u.String(), requests); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "request can not serialized.")
		}

		return c.String(http.StatusOK, "ok")
	}
}

func set(rd *redis.Client, key string, value interface{}) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	rd.Set(key, p, time.Hour*24)
	return nil
}

func getHeader(header http.Header) map[string]string {
	var headerMap = make(map[string]string)
	for key, value := range header {
		headerMap[key] = strings.Join(value, ",")
	}
	return headerMap
}

func isSupportedContentType(contentType string) bool {
	return strings.Contains(contentType, echo.MIMEMultipartForm)
}
