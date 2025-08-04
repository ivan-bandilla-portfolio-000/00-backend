package middlewares

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// RateLimiterService is the interface for rate limiting backends
type RateLimiterService interface {
	Allow(key string, limit int, ttl time.Duration) (allowed bool, retryAfter int, err error)
}

func RateLimitMiddlewareWithKey(
	service RateLimiterService,
	baseKey string,
	rps, burst int,
	ttl time.Duration,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if os.Getenv("APP_ENV") == "local" {
				next.ServeHTTP(w, r)
				return
			}

			ip := r.Header.Get("X-Forwarded-For")
			if ip == "" {
				ip = r.RemoteAddr
			}
			key := "ratelimit:" + baseKey + ":" + ip
			limit := rps * burst

			allowed, retryAfter, err := service.Allow(key, limit, ttl)
			if err != nil {
				log.Printf("Rate limiter error for key %s: %v", key, err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			if !allowed {
				w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RateLimitMiddlewareWithKeyDefault(
	service RateLimiterService,
	baseKey string,
	rps, burst int,
	ttl time.Duration,
) func(http.Handler) http.Handler {
	if rps <= 0 {
		rps = 1
	}
	if burst <= 0 {
		burst = 1
	}
	if ttl <= 0 {
		ttl = time.Minute
	}
	return RateLimitMiddlewareWithKey(service, baseKey, rps, burst, ttl)
}
