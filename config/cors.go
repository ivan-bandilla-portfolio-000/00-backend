package config

import (
	"strings"

	"portfolio-backend/utils"
)

type CORSConfig struct {
	Paths                  []string
	AllowedMethods         []string
	AllowedOrigins         []string
	AllowedOriginsPatterns []string
	AllowedHeaders         []string
	ExposedHeaders         []string
	MaxAge                 int
	SupportsCredentials    bool
}

func LoadCORSConfig() *CORSConfig {
	rawOrigins := utils.GetEnvOrDefault("ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:3000")
	origins := []string{}
	for _, o := range strings.Split(rawOrigins, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			origins = append(origins, o)
		}
	}
	return &CORSConfig{
		Paths:                  []string{"api/*", "sanctum/csrf-cookie"},
		AllowedMethods:         []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedOrigins:         origins,
		AllowedOriginsPatterns: []string{},
		AllowedHeaders:         []string{"*"},
		ExposedHeaders:         []string{},
		MaxAge:                 3600,
		SupportsCredentials:    false,
	}
}
