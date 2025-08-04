package redis

import (
	"context"
	"time"

	"portfolio-backend/app/middlewares"
	"portfolio-backend/config"

	"github.com/redis/go-redis/v9"
)

type UpstashService struct {
	Client *redis.Client
}

func NewUpstashService() *UpstashService {
	client := config.NewRedisClient()
	return &UpstashService{Client: client}
}

// Allow implements the RateLimiterService interface
func (u *UpstashService) Allow(key string, limit int, ttl time.Duration) (bool, int, error) {
	ctx := context.Background()
	count, err := u.Client.Incr(ctx, key).Result()
	if err != nil {
		return false, int(ttl.Seconds()), err
	}
	if count == 1 {
		u.Client.Expire(ctx, key, ttl)
	}
	if int(count) > limit {
		return false, int(ttl.Seconds()), nil
	}
	return true, 0, nil
}

// Ensure UpstashService implements RateLimiterService
var _ middlewares.RateLimiterService = (*UpstashService)(nil)
