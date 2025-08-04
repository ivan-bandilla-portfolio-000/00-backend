package config

type RateLimitingConfig struct {
	RateLimit  int // requests per second
	BurstLimit int // max burst size
}

func GlobalLoadRateLimitingConfig() RateLimitingConfig {
	return RateLimitingConfig{
		RateLimit:  5,
		BurstLimit: 10,
	}
}
