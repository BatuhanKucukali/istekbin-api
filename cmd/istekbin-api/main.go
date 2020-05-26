package main

import (
	"fmt"
	_ "github.com/batuhankucukali/istekbin-api/docs"
	"github.com/batuhankucukali/istekbin-api/internal/config"
	"github.com/batuhankucukali/istekbin-api/internal/datastore"
	"github.com/batuhankucukali/istekbin-api/internal/router"
)

// @title Istekbin API
// @description Istekbin is a free service that allows you to collect http request.

// @contact.name API Support
// @contact.url https://github.com/BatuhanKucukali/istekbin-api/issues/new

// @license.name Apache 2.0
// @license.url https://github.com/BatuhanKucukali/istekbin-api/blob/master/LICENSE

// @host api.istekbin.com
// @BasePath /
// @schemes https
func main() {
	conf := config.InitConfig()
	rd := datastore.InitRedis(conf.RedisConfig)
	e := router.Init(conf, rd)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", conf.AppConfig.Port)))
}
