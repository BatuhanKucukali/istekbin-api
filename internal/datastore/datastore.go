package datastore

import (
	"fmt"
	"github.com/batuhankucukali/istekbin/internal/config"
	"github.com/go-redis/redis/v7"
	"log"
)

func InitRedis(conf config.Redis) *redis.Client {
	rd := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       conf.DB,
	})

	if err := rd.Ping().Err(); err != nil {
		log.Fatal("Redis connection error.", err)
	}

	return rd
}
