package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	adto "github.com/k6zma/DockerMonitoringApp/backend/internal/application/dto"
	"github.com/k6zma/DockerMonitoringApp/backend/internal/application/usecases"
	pdto "github.com/k6zma/DockerMonitoringApp/backend/internal/presentation/dto"
	"github.com/k6zma/DockerMonitoringApp/backend/internal/presentation/mapper"
	"github.com/k6zma/DockerMonitoringApp/backend/pkg/utils"
)

type ContainerStatusHandler struct {
	useCase  usecases.ContainerStatusUseCaseInterface
	validate *validator.Validate
	logger   utils.LoggerInterface
}

func NewContainerStatusHandler(
	useCase usecases.ContainerStatusUseCaseInterface,
	logger utils.LoggerInterface,
) *ContainerStatusHandler {
	return &ContainerStatusHandler{
		useCase:  useCase,
		validate: validator.New(),
		logger:   logger,
	}
}

// GetFilteredContainerStatuses godoc
// @Summary Retrieve a list of containers
// @Description Returns a list of containers with optional filtering by various parameters
// @Tags Containers
// @Accept json
// @Produce json
// @Param container_id query string false "Filter by container ID"
// @Param ip query string false "Filter by IP"
// @Param name query string false "Filter by name"
// @Param status query string false "Filter by status"
// @Param ping_time_min query number false "Filter by minimum ping time"
// @Param ping_time_max query number false "Filter by maximum ping time"
// @Param created_at_gte query string false "Filter by creation date (greater than or equal to), format: RFC3339"
// @Param created_at_lte query string false "Filter by creation date (less than or equal to), format: RFC3339"
// @Param updated_at_gte query string false "Filter by last update date (greater than or equal to), format: RFC3339"
// @Param updated_at_lte query string false "Filter by last update date (less than or equal to), format: RFC3339"
// @Param limit query int false "Limit the number of returned records"
// @Success 200 {array} dto.GetContainerStatusResponse
// @Failure 500 {string} string "Internal Server Error"
// @Security ApiKeyAuth
// @Router /container_status [get].
func (h *ContainerStatusHandler) GetFilteredContainerStatuses(
	w http.ResponseWriter,
	r *http.Request,
) {
	h.logger.Debugf("HANDLERS: received GetFilteredContainerStatuses request with query: %s", r.URL.RawQuery)

	queryParams := r.URL.Query()
	filter := adto.ContainerStatusFilter{}

	if containerID := queryParams.Get("container_id"); containerID != "" {
		filter.ContainerID = &containerID
	}

	if ip := queryParams.Get("ip"); ip != "" {
		filter.IPAddress = &ip
	}

	if name := queryParams.Get("name"); name != "" {
		filter.Name = &name
	}

	if status := queryParams.Get("status"); status != "" {
		filter.Status = &status
	}

	if pingMinStr := queryParams.Get("ping_time_min"); pingMinStr != "" {
		pingMin, err := strconv.ParseFloat(pingMinStr, 64)
		if err == nil {
			filter.PingTimeMin = &pingMin
		} else {
			h.logger.Errorf("HANDLERS: error parsing ping_time_min param: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if pingMaxStr := queryParams.Get("ping_time_max"); pingMaxStr != "" {
		pingMax, err := strconv.ParseFloat(pingMaxStr, 64)
		if err == nil {
			filter.PingTimeMax = &pingMax
		} else {
			h.logger.Errorf("HANDLERS: error parsing ping_time_max param: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if createdAtGteStr := queryParams.Get("created_at_gte"); createdAtGteStr != "" {
		createdAtGte, err := time.Parse(time.RFC3339, createdAtGteStr)
		if err == nil {
			filter.CreatedAtGte = &createdAtGte
		} else {
			h.logger.Errorf("HANDLERS: error parsing created_at_gte param: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if createdAtLteStr := queryParams.Get("created_at_lte"); createdAtLteStr != "" {
		createdAtLte, err := time.Parse(time.RFC3339, createdAtLteStr)
		if err == nil {
			filter.CreatedAtLte = &createdAtLte
		} else {
			h.logger.Errorf("HANDLERS: error parsing created_at_lte param: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if updatedAtGteStr := queryParams.Get("updated_at_gte"); updatedAtGteStr != "" {
		updatedAtGte, err := time.Parse(time.RFC3339, updatedAtGteStr)
		if err == nil {
			filter.UpdatedAtGte = &updatedAtGte
		} else {
			h.logger.Errorf("HANDLERS: error parsing updated_at_gte param: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if updatedAtLteStr := queryParams.Get("updated_at_lte"); updatedAtLteStr != "" {
		updatedAtLte, err := time.Parse(time.RFC3339, updatedAtLteStr)
		if err == nil {
			filter.UpdatedAtLte = &updatedAtLte
		} else {
			h.logger.Errorf("HANDLERS: error parsing updated_at_lte param: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if limitStr := queryParams.Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err == nil {
			filter.Limit = &limit
		} else {
			h.logger.Errorf("HANDLERS: error parsing limit param: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	statuses, err := h.useCase.FindContainerStatuses(&filter)
	if err != nil {
		h.logger.Errorf("HANDLERS: getFilteredContainerStatuses error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.logger.Debugf("HANDLERS: found %d container statuses", len(statuses))
	response := mapper.MapAppDTOsToResponse(statuses)

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Errorf("HANDLERS: error encoding response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// CreateContainerStatus godoc
// @Summary Create a new container
// @Description Adds a new container to the database
// @Tags Containers
// @Accept json
// @Produce json
// @Param request body dto.CreateContainerStatusRequest true "Container data"
// @Success 201 {object} dto.GetContainerStatusResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Security ApiKeyAuth
// @Router /container_status [post].
func (h *ContainerStatusHandler) CreateContainerStatus(w http.ResponseWriter, r *http.Request) {
	h.logger.Debugf("HANDLERS: received CreateContainerStatus request")

	var req pdto.CreateContainerStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("HANDLERS: createContainerStatus decode error: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		h.logger.Errorf("HANDLERS: createContainerStatus validation error: %v", err)
		http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	appDTO := mapper.MapCreateRequestToAppDTO(req)

	createdStatus, err := h.useCase.CreateContainerStatus(&appDTO)
	if err != nil {
		h.logger.Errorf("HANDLERS: createContainerStatus error: %v", err)
		http.Error(w, "Failed to create container status", http.StatusInternalServerError)
		return
	}

	h.logger.Debugf("HANDLERS: container status created with container_id: %s", createdStatus.ContainerID)

	response := mapper.MapAppDTOToResponse(*createdStatus)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Errorf("HANDLERS: error encoding response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// UpdateContainerStatus godoc
// @Summary Update container by container ID
// @Description Partially updates a container by its container ID
// @Tags Containers
// @Accept json
// @Produce json
// @Param container_id path string true "Container ID"
// @Param request body dto.UpdateContainerStatusRequest true "Fields to update"
// @Success 204
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Security ApiKeyAuth
// @Router /container_status/{container_id} [patch].
func (h *ContainerStatusHandler) UpdateContainerStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	containerID := vars["container_id"]

	h.logger.Debugf("HANDLERS: received UpdateContainerStatus request for container_id: %s", containerID)

	var req pdto.UpdateContainerStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("HANDLERS: updateContainerStatus decode error for container_id %s: %v", containerID, err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.PingTime == 0 && req.LastSuccessfulPing.IsZero() && req.Status == "" {
		h.logger.Errorf("HANDLERS: updateContainerStatus validation error for container_id %s: No fields provided", containerID)
		http.Error(w, "At least one field must be provided", http.StatusBadRequest)
		return
	}

	appDTO := mapper.MapUpdateRequestToAppDTO(req)

	err := h.useCase.UpdateContainerStatus(containerID, &appDTO)
	if err != nil {
		h.logger.Errorf("HANDLERS: failed to update container status for container_id %s: %v", containerID, err)
		http.Error(w, "Failed to update container status", http.StatusInternalServerError)
		return
	}

	h.logger.Debugf("HANDLERS: successfully updated container status for container_id: %s", containerID)
	w.WriteHeader(http.StatusNoContent)
}

// DeleteContainerStatus godoc
// @Summary Delete container by container ID
// @Description Deletes a container from the database
// @Tags Containers
// @Accept json
// @Produce json
// @Param container_id path string true "Container ID"
// @Success 204 "No Content"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Security ApiKeyAuth
// @Router /container_status/{container_id} [delete].
func (h *ContainerStatusHandler) DeleteContainerStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	containerID := vars["container_id"]

	h.logger.Debugf("HANDLERS: received DeleteContainerStatus request for container_id: %s", containerID)

	err := h.useCase.DeleteContainerStatusByContainerID(containerID)
	if err != nil {
		if err.Error() == fmt.Sprintf("container status with container_id %s not found", containerID) {
			h.logger.Warnf("HANDLERS: container status with container_id %s not found", containerID)
			http.Error(w, "Container not found", http.StatusNotFound)
			return
		}
		h.logger.Errorf("HANDLERS: failed to delete container status for container_id %s: %v", containerID, err)
		http.Error(w, "Failed to delete container status", http.StatusInternalServerError)
		return
	}

	h.logger.Debugf("HANDLERS: successfully deleted container status for container_id: %s", containerID)
	w.WriteHeader(http.StatusNoContent)
}
