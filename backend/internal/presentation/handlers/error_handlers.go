package handlers

import (
	"net/http"

	"github.com/k6zma/DockerMonitoringApp/backend/pkg/utils"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := utils.GetClientIP(r)
	utils.LoggerInstance.Warnf("REQUESTS: - %s - %s - %s - %d - %s", clientIP, r.Method, r.URL.Path, http.StatusNotFound, "0ms")
	http.Error(w, "404 page not found", http.StatusNotFound)
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	clientIP := utils.GetClientIP(r)
	utils.LoggerInstance.Warnf("REQUESTS: - %s - %s - %s - %d - %s", clientIP, r.Method, r.URL.Path, http.StatusMethodNotAllowed, "0ms")
	http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
}
