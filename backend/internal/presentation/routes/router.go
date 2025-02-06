package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/k6zma/DockerMonitoringApp/backend/internal/infrastructure/config"
	"github.com/k6zma/DockerMonitoringApp/backend/internal/presentation/handlers"
	"github.com/k6zma/DockerMonitoringApp/backend/internal/presentation/middlewares"
	"github.com/k6zma/DockerMonitoringApp/backend/pkg/utils"
)

func InitRoutes(cfg *config.Config, errHandlers *handlers.ErrorHandlers, conHandlers *handlers.ContainerStatusHandler, logger utils.LoggerInterface) *mux.Router {
	router := mux.NewRouter()

	router.Use(middlewares.LoggingMiddleware(logger))

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	router.NotFoundHandler = http.HandlerFunc(errHandlers.NotFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(errHandlers.MethodNotAllowedHandler)

	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	apiRouter.Use(middlewares.AuthMiddleware(cfg, logger))

	apiRouter.HandleFunc("/container_status", conHandlers.GetFilteredContainerStatuses).Methods("GET")
	apiRouter.HandleFunc("/container_status", conHandlers.CreateContainerStatus).Methods("POST")
	apiRouter.HandleFunc("/container_status/{ip}", conHandlers.UpdateContainerStatus).Methods("PATCH")
	apiRouter.HandleFunc("/container_status/{ip}", conHandlers.DeleteContainerStatus).Methods("DELETE")

	return router
}
