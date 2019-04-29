package model

import (
	"github.com/go-redis/redis"
	"github.com/gookit/config"
	"strconv"
)

type Redis struct {
}

var Client *redis.Client

func init() {
	db, _ := strconv.Atoi(config.GetEnv("REDIS_SELECT"))
	Client = redis.NewClient(&redis.Options{
		Addr:     config.GetEnv("REDIS_HOST") + ":" + config.GetEnv("REDIS_PORT"),
		Password: config.GetEnv("REDIS_PASSWORD"),
		DB:       db,
		PoolSize: 100,
	})
}
