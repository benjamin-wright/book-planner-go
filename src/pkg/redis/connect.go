package redis

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

func Connect(config ConnectConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("https://%s:%d", config.Host, config.Port),
		Password: "",
		DB:       config.Database,
	})

	return rdb, nil
}
