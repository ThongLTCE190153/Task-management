package main

import (
	"log"

	"trithong.com/task-golang/internal/cache"
	"trithong.com/task-golang/internal/config"
	"trithong.com/task-golang/internal/jobs"
)

func main() {
	cfg := config.LoadConfig()

	if err := cfg.Validate(); err != nil {
		log.Fatal(err)
	}

	redisClient := cache.ConnectRedis(cfg)
	defer redisClient.Close()

	log.Println("Worker started, waiting for jobs...")

	jobs.StartNotificationWorker(redisClient)
}