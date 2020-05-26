package api

import (
	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v7"
)

var redisServer *miniredis.Miniredis

func mockRedis() *miniredis.Miniredis {
	m, err := miniredis.Run()

	if err != nil {
		panic(err)
	}

	return m
}

func redisClient() *redis.Client {
	redisServer = mockRedis()
	return redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})
}

func teardown() {
	redisServer.Close()
}
