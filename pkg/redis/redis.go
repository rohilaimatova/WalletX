package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("could not connect to Redis: %v", err)
	}

	log.Println("Connected to Redis")

	return rdb
}
