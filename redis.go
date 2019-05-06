package pkg

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/gookit/config"
)

type Redis struct {
}

var Client *redis.Client

func init() {
	db, _ := strconv.Atoi(config.GetEnv("REDIS_SELECT", "1"))
	Client = redis.NewClient(&redis.Options{
		Addr:     config.GetEnv("REDIS_HOST", "127.0.0.1") + ":" + config.GetEnv("REDIS_PORT", "3306"),
		Password: config.GetEnv("REDIS_PASSWORD", ""),
		DB:       db,
		PoolSize: 100,
	})
}

func Set(key string, value interface{}, ttl int) {
	val, err := json.Marshal(value)

	if err != nil {
		panic(err)
	}
	err = Client.Set(key, string(val), time.Duration(ttl)*time.Second).Err()
	if err != nil {
		panic(err)
	}
}

func Get(key string, structs interface{}) interface{} {
	value, err := Client.Get(key).Result()
	if err != nil {
		panic(err)
	}

	str := []byte(value)
	err = json.Unmarshal(str, &structs)
	if err != nil {
		panic(err)
	}
	return structs
}
