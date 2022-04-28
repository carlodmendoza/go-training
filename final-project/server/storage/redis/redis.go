package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisDB struct {
	Client *redis.Client
}

func Initialize(host, port string) *RedisDB {
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", host, port),
		Password: "",
		DB:       0,
	})

	client.SAdd(ctx, "users", "cmendoza")
	client.SetNX(ctx, "nextUserID", 1, 0)
	client.SetNX(ctx, "nextTransactionID", 1, 0)

	return &RedisDB{Client: client}
}

func (rdb *RedisDB) Shutdown() error {
	return rdb.Client.Close()
}
