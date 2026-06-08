package jobs

import (
	"encoding/json"

	"trithong.com/task-golang/internal/cache"

	"github.com/redis/go-redis/v9"
)

const NotificationQueueKey = "notification_jobs"

func PushNotificationJob(redisClient *redis.Client, job NotificationJob) error {
	data, err := json.Marshal(job)

	if err != nil {
		return err
	}

	return redisClient.RPush(cache.Ctx, NotificationQueueKey, data).Err()
}
