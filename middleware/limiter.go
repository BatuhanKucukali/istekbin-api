package middleware

import (
	"fmt"
	"github.com/batuhankucukali/istekbin/config"
	"github.com/labstack/echo/v4"
	_ "github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
	_ "log"
	_ "net"
	"sync"
	"time"
)

type visitor struct {
	limiter  *rate.Limiter
	rateConfig  config.Rate
	lastSeen time.Time
	exceedTime time.Time
	exceeded bool
	key string
}

var visitors = make(map[string]*visitor)
var mu sync.Mutex

func (v *visitor) Allow() bool {
	if v.exceeded == true {
		if time.Since(v.exceedTime) > v.rateConfig.Every {
			delete(visitors, v.key)
			return true
		}
		return false
	}
	if v.limiter.Allow() == false {
		v.exceeded = true
		v.exceedTime = time.Now()
		return false
	}
	return true
}

func cleanupVisitors(rc config.Rate) {
	for {
		time.Sleep(time.Minute * 1)

		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 2*rc.Every {
				fmt.Printf("cleared %s", ip)
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

func getVisitor(uuid string, rc config.Rate) *visitor {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[uuid]
	if !exists {
		limiter := rate.NewLimiter(rate.Every(rc.Every), rc.Limit)
		v = &visitor{limiter, rc, time.Now(), time.Now().Add(5* rc.Every), false, uuid}
		visitors[uuid] = v
	}

	v.lastSeen = time.Now()
	return v
}

func RateLimit(rc config.Rate) echo.MiddlewareFunc {
	go cleanupVisitors(rc)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			u := c.Param("uuid")
			v := getVisitor(u, rc)
			if v.Allow() == false {
				return echo.ErrTooManyRequests
			}

			next(c)
			return
		}
	}
}
