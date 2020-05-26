package router

import (
	"github.com/batuhankucukali/istekbin/internal/config"
	"github.com/labstack/echo/v4"
	_ "github.com/labstack/echo/v4"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/redis"
	"log"
	_ "log"
	_ "net"
	"net/http"
	"strconv"
)

var (
	ipRateLimiter *limiter.Limiter
	store         limiter.Store
)

func RateLimit(conf *config.Rate, rd redis.Client) echo.MiddlewareFunc {
	rate := limiter.Rate{
		Period: conf.Period,
		Limit:  conf.Limit,
	}
	store, err := redis.NewStore(rd)
	if err != nil {
		panic(err)
	}
	ipRateLimiter = limiter.New(store, rate)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			ip := c.RealIP()
			limiterCtx, err := ipRateLimiter.Get(c.Request().Context(), ip)
			if err != nil {
				log.Printf("IPRateLimit - ipRateLimiter.Get - err: %v, %s on %s", err, ip, c.Request().URL)
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": err,
				})
			}

			h := c.Response().Header()
			h.Set("X-RateLimit-Limit", strconv.FormatInt(limiterCtx.Limit, 10))
			h.Set("X-RateLimit-Remaining", strconv.FormatInt(limiterCtx.Remaining, 10))
			h.Set("X-RateLimit-Reset", strconv.FormatInt(limiterCtx.Reset, 10))

			if limiterCtx.Reached {
				log.Printf("Too Many Requests from %s on %s", ip, c.Request().URL)
				return c.JSON(http.StatusTooManyRequests, echo.Map{
					"message": "Too Many Requests on " + c.Request().URL.String(),
				})
			}

			return next(c)
		}
	}
}
