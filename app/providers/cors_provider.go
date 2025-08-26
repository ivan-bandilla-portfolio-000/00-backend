package providers

import (
	"net/http"
	"strconv"
	"strings"

	"portfolio-backend/config"
)

type CorsProvider struct {
	config *config.CORSConfig
}

func NewCorsProvider(cfg *config.CORSConfig) *CorsProvider {
	return &CorsProvider{config: cfg}
}

func (cp *CorsProvider) allowedOrigin(origin string) bool {
	if origin == "" {
		return false
	}
	for _, o := range cp.config.AllowedOrigins {
		if strings.EqualFold(o, origin) {
			return true
		}
	}
	for _, p := range cp.config.AllowedOriginsPatterns {
		if p != "" && strings.Contains(origin, p) {
			return true
		}
	}
	return false
}

func (cp *CorsProvider) Handler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		// not a cross-origin request -> continue
		if origin == "" {
			next(w, r)
			return
		}

		if !cp.allowedOrigin(origin) {
			// origin not allowed: do not set CORS headers (browser will block)
			next(w, r)
			return
		}

		// Echo the origin (required when not using "*")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Add("Vary", "Origin")

		if cp.config.SupportsCredentials {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// Methods
		if len(cp.config.AllowedMethods) > 0 {
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(cp.config.AllowedMethods, ", "))
		} else {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		}

		// Headers
		if len(cp.config.AllowedHeaders) == 1 && cp.config.AllowedHeaders[0] == "*" {
			reqHeaders := r.Header.Get("Access-Control-Request-Headers")
			if reqHeaders == "" {
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-Timezone")
			} else {
				w.Header().Set("Access-Control-Allow-Headers", reqHeaders)
			}
		} else if len(cp.config.AllowedHeaders) > 0 {
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(cp.config.AllowedHeaders, ", "))
		}

		if len(cp.config.ExposedHeaders) > 0 {
			w.Header().Set("Access-Control-Expose-Headers", strings.Join(cp.config.ExposedHeaders, ", "))
		}

		if cp.config.MaxAge > 0 {
			w.Header().Set("Access-Control-Max-Age", strconv.Itoa(cp.config.MaxAge))
		}

		// Preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next(w, r)
	}
}
