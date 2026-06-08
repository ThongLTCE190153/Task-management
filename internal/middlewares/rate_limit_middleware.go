package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"trithong.com/task-golang/internal/responses"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimitMiddleware(redisClient *redis.Client, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {

		userID := c.GetInt("user_id")

		key := fmt.Sprintf("rate_limit:user:%d", userID)

		ctx := context.Background()

		count, err := redisClient.Get(ctx, key).Int()

		if err != nil && err != redis.Nil {
			c.Next()
			return
		}

		if count >= limit {
			responses.Error(
				c,
				http.StatusTooManyRequests,
				fmt.Sprintf("Too many requests, please try again after %s", window),
			)
			c.Abort()
			return
		}

		pipe := redisClient.Pipeline()
		pipe.Incr(ctx, key)

		if count == 0 {
			pipe.Expire(ctx, key, window)
		}

		pipe.Exec(ctx)

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", limit-count-1))

		c.Next()
	}
}