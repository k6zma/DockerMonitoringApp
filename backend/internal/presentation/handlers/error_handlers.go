package handlers

import (
	"net/http"

	"github.com/k6zma/DockerMonitoringApp/backend/pkg/utils"
)

type ErrorHandlers struct {
	logger utils.LoggerInterface
}

func NewErrorHandlers(logger utils.LoggerInterface) *ErrorHandlers {
	return &ErrorHandlers{
		logger: logger,
	}
}

func (e *ErrorHandlers) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := utils.GetClientIP(r)
	e.logger.Warnf(
		"REQUESTS: - %s - %s - %s - %d - %s",
		clientIP,
		r.Method,
		r.URL.Path,
		http.StatusNotFound,
		"0ms",
	)
	http.Error(w, "404 page not found", http.StatusNotFound)
}

func (e *ErrorHandlers) MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := utils.GetClientIP(r)
	e.logger.Warnf(
		"REQUESTS: - %s - %s - %s - %d - %s",
		clientIP,
		r.Method,
		r.URL.Path,
		http.StatusMethodNotAllowed,
		"0ms",
	)
	http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
}
