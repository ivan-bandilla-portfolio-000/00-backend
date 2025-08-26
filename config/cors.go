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
		AllowedMethods:         []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedOrigins:         []string{"http://localhost:5173", "http://localhost:3000"},
		AllowedOriginsPatterns: []string{},
		AllowedHeaders:         []string{"*"},
		ExposedHeaders:         []string{},
		MaxAge:                 3600,
		SupportsCredentials:    false,
	}
}
