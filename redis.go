package pkg

import (
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
	/* switch value.(type) {
	case string:
		value = value
		break
	default:
		// case array:
		value, err := json.Marshal(value)

		if err != nil {
			panic(err)
		}
		break
	} */
	// value, err := json.Marshal(value)

	// if err != nil {
	// 	panic(err)
	// }
	err := Client.Set(key, value, time.Duration(ttl)*time.Second).Err()
	if err != nil {
		panic(err)
	}
}

func Get(key string) string {
	val, err := Client.Get(key).Result()
	if err != nil {
		panic(err)
	}

	return val
}
