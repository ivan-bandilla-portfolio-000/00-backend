package services

import (
	"context"
	"time"

	"portfolio-backend/app/middlewares"

	"github.com/redis/go-redis/v9"
)

type RedisRateLimiter struct {
	Client *redis.Client
}

func (r *RedisRateLimiter) Allow(key string, limit int, ttl time.Duration) (bool, int, error) {
	ctx := context.Background()
	count, err := r.Client.Incr(ctx, key).Result()
	if err != nil {
		return false, int(ttl.Seconds()), err
	}
	if count == 1 {
		r.Client.Expire(ctx, key, ttl)
	}
	if int(count) > limit {
		return false, int(ttl.Seconds()), nil
	}
	return true, 0, nil
}

// Ensure RedisRateLimiter implements RateLimiterService
var _ middlewares.RateLimiterService = (*RedisRateLimiter)(nil)
