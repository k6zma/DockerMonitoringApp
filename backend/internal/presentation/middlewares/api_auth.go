package middlewares

import (
	"net/http"

	"github.com/k6zma/DockerMonitoringApp/backend/internal/infrastructure/config"
	"github.com/k6zma/DockerMonitoringApp/backend/pkg/utils"
)

func AuthMiddleware(cfg *config.Config, logger utils.LoggerInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-Api-Key")
			if apiKey == "" || apiKey != cfg.AuthAPI.APIKey {
				logger.Warnf("MIDDLEWARE: unauthorized access attempt")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
