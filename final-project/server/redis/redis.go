package redis

import "github.com/go-redis/redis/v8"

var Client = NewClient()

func NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
}
