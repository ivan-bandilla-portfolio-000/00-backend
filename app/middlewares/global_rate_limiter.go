package middlewares

import (
	"net/http"
	"portfolio-backend/config"
	"sync"

	"golang.org/x/time/rate"
)

var (
	limiter *rate.Limiter
	once    sync.Once
)

func getLimiter() *rate.Limiter {
	once.Do(func() {
		cfg := config.GlobalLoadRateLimitingConfig()
		limiter = rate.NewLimiter(rate.Limit(cfg.RateLimit), cfg.BurstLimit)
	})
	return limiter
}

func GlobalRateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !getLimiter().Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
