package cache

import (
	"context"
	"log"

	"time"

	"trithong.com/task-golang/internal/config"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func ConnectRedis(cfg config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisAddr,
		Password:     cfg.RedisPassword,
		DB:           cfg.RedisDB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	if err := client.Ping(Ctx).Err(); err != nil {
		log.Fatal("Failed to connect Redis: ", err)
	}

	log.Println("Connected to Redis successfully")

	return client
}
