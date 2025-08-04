package providers

import (
	"net/http"
	"portfolio-backend/config"
	"strings"
)

type CorsProvider struct {
	config *config.CORSConfig
}

func NewCorsProvider(config *config.CORSConfig) *CorsProvider {
	return &CorsProvider{config: config}
}

func (cp *CorsProvider) Handler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", strings.Join(cp.config.AllowedOrigins, ","))
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(cp.config.AllowedMethods, ","))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(cp.config.AllowedHeaders, ","))
		if len(cp.config.ExposedHeaders) > 0 {
			w.Header().Set("Access-Control-Expose-Headers", strings.Join(cp.config.ExposedHeaders, ","))
		}
		if cp.config.SupportsCredentials {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		if cp.config.MaxAge > 0 {
			w.Header().Set("Access-Control-Max-Age", string(cp.config.MaxAge))
		}
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}
