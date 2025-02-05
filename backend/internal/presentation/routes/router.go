package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/k6zma/DockerMonitoringApp/backend/internal/infrastructure/config"
	"github.com/k6zma/DockerMonitoringApp/backend/internal/presentation/handlers"
	"github.com/k6zma/DockerMonitoringApp/backend/internal/presentation/middlewares"
)

func InitRoutes(cfg *config.Config, handler *handlers.ContainerStatusHandler) *mux.Router {
	router := mux.NewRouter()

	router.Use(middlewares.LoggingMiddleware)

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	router.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(handlers.MethodNotAllowedHandler)

	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	apiRouter.Use(middlewares.AuthMiddleware(cfg))

	apiRouter.HandleFunc("/container_status", handler.GetFilteredContainerStatuses).Methods("GET")
	apiRouter.HandleFunc("/container_status", handler.CreateContainerStatus).Methods("POST")
	apiRouter.HandleFunc("/container_status/{ip}", handler.UpdateContainerStatus).Methods("PATCH")
	apiRouter.HandleFunc("/container_status/{ip}", handler.DeleteContainerStatus).Methods("DELETE")

	return router
}
