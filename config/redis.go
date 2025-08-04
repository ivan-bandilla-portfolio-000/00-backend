package config

import (
	"os"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient returns a new Redis client using the REDIS_URL environment variable
func NewRedisClient() *redis.Client {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic("Invalid REDIS_URL: " + err.Error())
	}
	return redis.NewClient(opt)
}
