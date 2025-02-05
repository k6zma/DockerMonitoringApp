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
	useCase  *usecases.ContainerStatusUseCase
	validate *validator.Validate
}

func NewContainerStatusHandler(useCase *usecases.ContainerStatusUseCase) *ContainerStatusHandler {
	return &ContainerStatusHandler{
		useCase:  useCase,
		validate: validator.New(),
	}
}

// GetFilteredContainerStatuses godoc
// @Summary Retrieve a list of containers
// @Description Returns a list of containers with optional filtering by various parameters
// @Tags Containers
// @Accept json
// @Produce json
// @Param ip query string false "Filter by IP"
// @Param id query int false "Filter by ID"
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
func (h *ContainerStatusHandler) GetFilteredContainerStatuses(w http.ResponseWriter, r *http.Request) {
	utils.LoggerInstance.Debugf("HANDLERS: received GetFilteredContainerStatuses request with query: %s", r.URL.RawQuery)

	queryParams := r.URL.Query()
	filter := adto.ContainerStatusFilter{}

	if ip := queryParams.Get("ip"); ip != "" {
		filter.IPAddress = &ip
	}

	if idStr := queryParams.Get("id"); idStr != "" {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err == nil {
			filter.ID = &id
		} else {
			utils.LoggerInstance.Errorf("HANDLERS: error parsing id param: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if pingMinStr := queryParams.Get("ping_time_min"); pingMinStr != "" {
		pingMin, err := strconv.ParseFloat(pingMinStr, 64)
		if err == nil {
			filter.PingTimeMin = &pingMin
		} else {
			utils.LoggerInstance.Errorf("HANDLERS: error parsing ping_time_min param: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if pingMaxStr := queryParams.Get("ping_time_max"); pingMaxStr != "" {
		pingMax, err := strconv.ParseFloat(pingMaxStr, 64)
		if err == nil {
			filter.PingTimeMax = &pingMax
		} else {
			utils.LoggerInstance.Errorf("HANDLERS: error parsing ping_time_max param: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if createdAtGteStr := queryParams.Get("created_at_gte"); createdAtGteStr != "" {
		createdAtGte, err := time.Parse(time.RFC3339, createdAtGteStr)
		if err == nil {
			filter.CreatedAtGte = &createdAtGte
		} else {
			utils.LoggerInstance.Errorf("HANDLERS: error parsing created_at_gte param: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if createdAtLteStr := queryParams.Get("created_at_lte"); createdAtLteStr != "" {
		createdAtLte, err := time.Parse(time.RFC3339, createdAtLteStr)
		if err == nil {
			filter.CreatedAtLte = &createdAtLte
		} else {
			utils.LoggerInstance.Errorf("HANDLERS: error parsing created_at_lte param: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if updatedAtGteStr := queryParams.Get("updated_at_gte"); updatedAtGteStr != "" {
		updatedAtGte, err := time.Parse(time.RFC3339, updatedAtGteStr)
		if err == nil {
			filter.UpdatedAtGte = &updatedAtGte
		} else {
			utils.LoggerInstance.Errorf("HANDLERS: error parsing updated_at_gte param: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if updatedAtLteStr := queryParams.Get("updated_at_lte"); updatedAtLteStr != "" {
		updatedAtLte, err := time.Parse(time.RFC3339, updatedAtLteStr)
		if err == nil {
			filter.UpdatedAtLte = &updatedAtLte
		} else {
			utils.LoggerInstance.Errorf("HANDLERS: error parsing updated_at_lte param: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if limitStr := queryParams.Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err == nil {
			filter.Limit = &limit
		} else {
			utils.LoggerInstance.Errorf("HANDLERS: error parsing limit param: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	statuses, err := h.useCase.FindContainerStatuses(&filter)
	if err != nil {
		utils.LoggerInstance.Errorf("HANDLERS: getFilteredContainerStatuses error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	utils.LoggerInstance.Debugf("HANDLERS: found %d container statuses", len(statuses))
	response := mapper.MapAppDTOsToResponse(statuses)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.LoggerInstance.Errorf("HANDLERS: error encoding response: %v", err)
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
	utils.LoggerInstance.Debugf("HANDLERS: received CreateContainerStatus request")

	var req pdto.CreateContainerStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.LoggerInstance.Errorf("HANDLERS: createContainerStatus decode error: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.LoggerInstance.Errorf("HANDLERS: createContainerStatus validation error: %v", err)
		http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	appDTO := mapper.MapCreateRequestToAppDTO(req)

	createdStatus, err := h.useCase.CreateContainerStatus(&appDTO)
	if err != nil {
		utils.LoggerInstance.Errorf("HANDLERS: createContainerStatus error: %v", err)
		http.Error(w, "Failed to create container status", http.StatusInternalServerError)
		return
	}

	utils.LoggerInstance.Debugf("HANDLERS: container status created with ID: %d", createdStatus.ID)

	response := mapper.MapAppDTOToResponse(*createdStatus)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.LoggerInstance.Errorf("HANDLERS: error encoding response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// UpdateContainerStatus godoc
// @Summary Update container by IP
// @Description Partially updates a container by its IP address
// @Tags Containers
// @Accept json
// @Produce json
// @Param ip path string true "Container IP address"
// @Param request body dto.UpdateContainerStatusRequest true "Fields to update"
// @Success 204
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Security ApiKeyAuth
// @Router /container_status/{ip} [patch].
func (h *ContainerStatusHandler) UpdateContainerStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ip := vars["ip"]

	utils.LoggerInstance.Debugf("HANDLERS: received UpdateContainerStatus request for IP: %s", ip)

	var req pdto.UpdateContainerStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.LoggerInstance.Errorf("HANDLERS: updateContainerStatus decode error for IP %s: %v", ip, err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)

		return
	}

	if req.PingTime == 0 && req.LastSuccessfulPing.IsZero() {
		utils.LoggerInstance.Errorf("HANDLERS: updateContainerStatus validation error for IP %s: No fields provided", ip)
		http.Error(w, "At least one field must be provided", http.StatusBadRequest)

		return
	}

	appDTO := mapper.MapUpdateRequestToAppDTO(req)

	err := h.useCase.UpdateContainerStatus(ip, &appDTO)
	if err != nil {
		utils.LoggerInstance.Errorf("HANDLERS: failed to update container status for IP %s: %v", ip, err)
		http.Error(w, "Failed to update container status", http.StatusInternalServerError)

		return
	}

	utils.LoggerInstance.Debugf("HANDLERS: successfully updated container status for IP: %s", ip)
	w.WriteHeader(http.StatusNoContent)
}

// DeleteContainerStatus godoc
// @Summary Delete container by IP
// @Description Deletes a container from the database
// @Tags Containers
// @Accept json
// @Produce json
// @Param ip path string true "Container IP address"
// @Success 204 "No Content"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Security ApiKeyAuth
// @Router /container_status/{ip} [delete].
func (h *ContainerStatusHandler) DeleteContainerStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ip := vars["ip"]

	utils.LoggerInstance.Debugf("HANDLERS: received DeleteContainerStatus request for IP: %s", ip)

	err := h.useCase.DeleteContainerStatusByIP(ip)
	if err != nil {
		if err.Error() == fmt.Sprintf("container status with IP %s not found", ip) {
			utils.LoggerInstance.Warnf("HANDLERS: container status with IP %s not found", ip)
			http.Error(w, "Container not found", http.StatusNotFound)
			return
		}

		utils.LoggerInstance.Errorf("HANDLERS: failed to delete container status for IP %s: %v", ip, err)
		http.Error(w, "Failed to delete container status", http.StatusInternalServerError)
		return
	}

	utils.LoggerInstance.Debugf("HANDLERS: successfully deleted container status for IP: %s", ip)
	w.WriteHeader(http.StatusNoContent)
}
