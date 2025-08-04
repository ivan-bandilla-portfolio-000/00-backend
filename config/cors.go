package config

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
	return &CORSConfig{
		Paths:                  []string{"api/*", "sanctum/csrf-cookie"},
		AllowedMethods:         []string{"*"},
		AllowedOrigins:         []string{"http://localhost:5173"},
		AllowedOriginsPatterns: []string{},
		AllowedHeaders:         []string{"Content-Type", "X-Timezone"},
		ExposedHeaders:         []string{},
		MaxAge:                 0,
		SupportsCredentials:    false,
	}
}
