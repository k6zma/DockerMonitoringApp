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

func InitRoutes(
	cfg *config.Config,
	errHandler *handlers.ErrorHandlers,
	conHandler *handlers.ContainerStatusHandler,
	logger utils.LoggerInterface,
) *mux.Router {
	router := mux.NewRouter()

	router.Use(middlewares.LoggingMiddleware(logger))

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	router.NotFoundHandler = http.HandlerFunc(errHandler.NotFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(errHandler.MethodNotAllowedHandler)

	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	apiRouter.Use(middlewares.AuthMiddleware(cfg, logger))

	apiRouter.HandleFunc("/container_status", conHandler.GetFilteredContainerStatuses).
		Methods("GET")
	apiRouter.HandleFunc("/container_status", conHandler.CreateContainerStatus).Methods("POST")
	apiRouter.HandleFunc("/container_status/{ip}", conHandler.UpdateContainerStatus).
		Methods("PATCH")
	apiRouter.HandleFunc("/container_status/{ip}", conHandler.DeleteContainerStatus).
		Methods("DELETE")

	return router
}
