package jobs

import (
	"encoding/json"
	"log"
	"time"

	"trithong.com/task-golang/internal/cache"

	"github.com/redis/go-redis/v9"
)

func StartNotificationWorker(redisClient *redis.Client) {
	log.Println("Notification worker started")

	for {
		result, err := redisClient.BLPop(
			cache.Ctx,
			0,
			NotificationQueueKey,
		).Result()

		if err != nil {
			log.Println("Failed to pop notification job:", err)
			continue
		}

		if len(result) < 2 {
			continue
		}

		var job NotificationJob

		err = json.Unmarshal([]byte(result[1]), &job)

		if err != nil {
			log.Println("Failed to parse notification job:", err)
			continue
		}

		processNotificationJob(redisClient, job)
	}
}

func processNotificationJob(redisClient *redis.Client, job NotificationJob) {
	log.Println("Processing notification job...")
	log.Println("Job type:", job.Type)
	log.Println("Task ID:", job.TaskID)
	log.Println("Assignee ID:", job.AssigneeID)

	err := sendNotification(job)

	if err != nil {
		if job.RetryCount < job.MaxRetry {
			job.RetryCount++

			log.Printf("Notification failed, retry %d/%d\n", job.RetryCount, job.MaxRetry)

			err = PushNotificationJob(redisClient, job)

			if err != nil {
				log.Println("Failed to push retry notification job:", err)
			}

			return
		}

		log.Println("Notification job failed permanently")
		return
	}

	log.Println("Notification sent successfully")
}

func sendNotification(job NotificationJob) error {
	// Giả lập delay như đang gửi thật
	time.Sleep(1 * time.Second)

	// Giả lập gửi notification — log ra console
	log.Printf("[NOTIFICATION] Type: %s | Task ID: %d | Assignee ID: %d\n",
		job.Type,
		job.TaskID,
		job.AssigneeID,
	)

	return nil
}