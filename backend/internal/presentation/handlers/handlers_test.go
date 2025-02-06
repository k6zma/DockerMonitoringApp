package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	adto "github.com/k6zma/DockerMonitoringApp/backend/internal/application/dto"
	pdto "github.com/k6zma/DockerMonitoringApp/backend/internal/presentation/dto"
	"github.com/k6zma/DockerMonitoringApp/backend/internal/presentation/handlers"
	"github.com/k6zma/DockerMonitoringApp/backend/mocks"
)

const (
	ipAddress = "192.168.1.101"
	pingTime  = 15.5
)

func TestGetContainerStatuses_ReturnsDataSuccessfully(t *testing.T) {
	mockUseCase := new(mocks.ContainerStatusUseCaseInterface)
	mockLogger := new(mocks.LoggerInterface)

	handler := handlers.NewContainerStatusHandler(mockUseCase, mockLogger)

	expectedStatuses := []*adto.ContainerStatusDTO{
		{IPAddress: ipAddress, PingTime: pingTime, LastSuccessfulPing: time.Now()},
	}

	mockUseCase.On("FindContainerStatuses", mock.Anything).Return(expectedStatuses, nil)
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()

	req := httptest.NewRequest(http.MethodGet, "/container_status", http.NoBody)
	rec := httptest.NewRecorder()

	handler.GetFilteredContainerStatuses(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response []pdto.GetContainerStatusResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, ipAddress, response[0].IPAddress)

	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetContainerStatuses_WithQueryParams_ReturnsFilteredData(t *testing.T) {
	mockUseCase := new(mocks.ContainerStatusUseCaseInterface)
	mockLogger := new(mocks.LoggerInterface)

	handler := handlers.NewContainerStatusHandler(mockUseCase, mockLogger)

	expectedStatuses := []*adto.ContainerStatusDTO{
		{IPAddress: ipAddress, PingTime: pingTime, LastSuccessfulPing: time.Now()},
	}

	mockUseCase.On("FindContainerStatuses", mock.Anything).Return(expectedStatuses, nil)
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()

	req := httptest.NewRequest(http.MethodGet,
		"/container_status?ip=192.168.1.101&id=123&ping_time_min=10&ping_time_max=30&created_at_gte=2023-01-01T00:00:00Z"+
			"&created_at_lte=2023-12-31T23:59:59Z&updated_at_gte=2023-06-01T12:00:00Z&updated_at_lte=2023-12-31T23:59:59Z&limit=5",
		http.NoBody)
	rec := httptest.NewRecorder()

	handler.GetFilteredContainerStatuses(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response []pdto.GetContainerStatusResponse
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, ipAddress, response[0].IPAddress)

	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetContainerStatuses_ErrorFromUseCase_ReturnsInternalServerError(t *testing.T) {
	mockUseCase := new(mocks.ContainerStatusUseCaseInterface)
	mockLogger := new(mocks.LoggerInterface)

	handler := handlers.NewContainerStatusHandler(mockUseCase, mockLogger)

	mockUseCase.On("FindContainerStatuses", mock.Anything).Return(nil, fmt.Errorf("database error"))
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockLogger.On("Errorf", mock.Anything, mock.Anything).Return()

	req := httptest.NewRequest(http.MethodGet, "/container_status", http.NoBody)
	rec := httptest.NewRecorder()

	handler.GetFilteredContainerStatuses(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGetContainerStatuses_InvalidID_ReturnsInternalServerError(t *testing.T) {
	mockUseCase := new(mocks.ContainerStatusUseCaseInterface)
	mockLogger := new(mocks.LoggerInterface)

	handler := handlers.NewContainerStatusHandler(mockUseCase, mockLogger)

	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockLogger.On("Errorf", mock.Anything, mock.Anything).Return()

	req := httptest.NewRequest(http.MethodGet, "/container_status?id=invalid", http.NoBody)
	rec := httptest.NewRecorder()

	handler.GetFilteredContainerStatuses(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockLogger.AssertExpectations(t)
}

func TestGetContainerStatuses_InvalidPingTimeMin_ReturnsInternalServerError(t *testing.T) {
	mockUseCase := new(mocks.ContainerStatusUseCaseInterface)
	mockLogger := new(mocks.LoggerInterface)

	handler := handlers.NewContainerStatusHandler(mockUseCase, mockLogger)

	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockLogger.On("Errorf", mock.Anything, mock.Anything).Return()

	req := httptest.NewRequest(http.MethodGet, "/container_status?ping_time_min=invalid", http.NoBody)
	rec := httptest.NewRecorder()

	handler.GetFilteredContainerStatuses(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockLogger.AssertExpectations(t)
}

func TestGetContainerStatuses_InvalidPingTimeMax_ReturnsInternalServerError(t *testing.T) {
	mockUseCase := new(mocks.ContainerStatusUseCaseInterface)
	mockLogger := new(mocks.LoggerInterface)

	handler := handlers.NewContainerStatusHandler(mockUseCase, mockLogger)

	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockLogger.On("Errorf", mock.Anything, mock.Anything).Return()

	req := httptest.NewRequest(http.MethodGet, "/container_status?ping_time_max=invalid", http.NoBody)
	rec := httptest.NewRecorder()

	handler.GetFilteredContainerStatuses(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockLogger.AssertExpectations(t)
}

func TestGetContainerStatuses_InvalidCreatedAtGte_ReturnsInternalServerError(t *testing.T) {
	mockUseCase := new(mocks.ContainerStatusUseCaseInterface)
	mockLogger := new(mocks.LoggerInterface)

	handler := handlers.NewContainerStatusHandler(mockUseCase, mockLogger)

	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockLogger.On("Errorf", mock.Anything, mock.Anything).Return()

	req := httptest.NewRequest(http.MethodGet, "/container_status?created_at_gte=not-a-date", http.NoBody)
	rec := httptest.NewRecorder()

	handler.GetFilteredContainerStatuses(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockLogger.AssertExpectations(t)
}

func TestGetContainerStatuses_InvalidCreatedAtLte_ReturnsInternalServerError(t *testing.T) {
	mockUseCase := new(mocks.ContainerStatusUseCaseInterface)
	mockLogger := new(mocks.LoggerInterface)

	handler := handlers.NewContainerStatusHandler(mockUseCase, mockLogger)

	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockLogger.On("Errorf", mock.Anything, mock.Anything).Return()

	req := httptest.NewRequest(http.MethodGet, "/container_status?created_at_lte=not-a-date", http.NoBody)
	rec := httptest.NewRecorder()

	handler.GetFilteredContainerStatuses(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockLogger.AssertExpectations(t)
}

func TestGetContainerStatuses_InvalidUpdatedAtGte_ReturnsInternalServerError(t *testing.T) {
	mockUseCase := new(mocks.ContainerStatusUseCaseInterface)
	mockLogger := new(mocks.LoggerInterface)

	handler := handlers.NewContainerStatusHandler(mockUseCase, mockLogger)

	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockLogger.On("Errorf", mock.Anything, mock.Anything).Return()

	req := httptest.NewRequest(http.MethodGet, "/container_status?updated_at_gte=not-a-date", http.NoBody)
	rec := httptest.NewRecorder()

	handler.GetFilteredContainerStatuses(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockLogger.AssertExpectations(t)
}

func TestGetContainerStatuses_InvalidUpdatedAtLte_ReturnsInternalServerError(t *testing.T) {
	mockUseCase := new(mocks.ContainerStatusUseCaseInterface)
	mockLogger := new(mocks.LoggerInterface)

	handler := handlers.NewContainerStatusHandler(mockUseCase, mockLogger)

	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockLogger.On("Errorf", mock.Anything, mock.Anything).Return()

	req := httptest.NewRequest(http.MethodGet, "/container_status?updated_at_lte=not-a-date", http.NoBody)
	rec := httptest.NewRecorder()

	handler.GetFilteredContainerStatuses(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockLogger.AssertExpectations(t)
}

func TestGetContainerStatuses_InvalidLimit_ReturnsInternalServerError(t *testing.T) {
	mockUseCase := new(mocks.ContainerStatusUseCaseInterface)
	mockLogger := new(mocks.LoggerInterface)

	handler := handlers.NewContainerStatusHandler(mockUseCase, mockLogger)

	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockLogger.On("Errorf", mock.Anything, mock.Anything).Return()

	req := httptest.NewRequest(http.MethodGet, "/container_status?limit=not-a-number", http.NoBody)
	rec := httptest.NewRecorder()

	handler.GetFilteredContainerStatuses(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockLogger.AssertExpectations(t)
}

func TestCreateContainerStatus_SuccessfullyCreatesContainer(t *testing.T) {
	mockUseCase := new(mocks.ContainerStatusUseCaseInterface)
	mockLogger := new(mocks.LoggerInterface)

	handler := handlers.NewContainerStatusHandler(mockUseCase, mockLogger)

	timeData := time.Now()

	requestBody := pdto.CreateContainerStatusRequest{
		IPAddress:          ipAddress,
		PingTime:           pingTime,
		LastSuccessfulPing: timeData,
	}
	jsonBody, _ := json.Marshal(requestBody)

	mockUseCase.On("CreateContainerStatus", mock.Anything).Return(&adto.ContainerStatusDTO{
		IPAddress:          ipAddress,
		PingTime:           pingTime,
		LastSuccessfulPing: timeData,
	}, nil)
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()

	req := httptest.NewRequest(http.MethodPost, "/container_status", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.CreateContainerStatus(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code, "Response Body: %s", rec.Body.String())
	assert.Equal(t, http.StatusCreated, rec.Code)
	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestCreateContainerStatus_InvalidJSON_ReturnsBadRequest(t *testing.T) {
	mockUseCase := new(mocks.ContainerStatusUseCaseInterface)
	mockLogger := new(mocks.LoggerInterface)

	handler := handlers.NewContainerStatusHandler(mockUseCase, mockLogger)

	invalidJSON := []byte(`{"IPAddress": ` + ipAddress + `}`)

	mockLogger.On("Debugf", mock.Anything).Return()
	mockLogger.On("Errorf", mock.Anything, mock.Anything).Return()

	req := httptest.NewRequest(http.MethodPost, "/container_status", bytes.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.CreateContainerStatus(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockLogger.AssertExpectations(t)
}

func TestCreateContainerStatus_InvalidJSONDecode_ReturnsBadRequest(t *testing.T) {
	mockUseCase := new(mocks.ContainerStatusUseCaseInterface)
	mockLogger := new(mocks.LoggerInterface)

	handler := handlers.NewContainerStatusHandler(mockUseCase, mockLogger)

	invalidJSON := []byte(`{"IPAddress":}`)

	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockLogger.On("Errorf", mock.Anything, mock.Anything).Return()

	req := httptest.NewRequest(http.MethodPost, "/container_status", bytes.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.CreateContainerStatus(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockLogger.AssertExpectations(t)
}

func TestCreateContainerStatus_ErrorFromUseCase_ReturnsInternalServerError(t *testing.T) {
	mockUseCase := new(mocks.ContainerStatusUseCaseInterface)
	mockLogger := new(mocks.LoggerInterface)

	handler := handlers.NewContainerStatusHandler(mockUseCase, mockLogger)

	requestBody := pdto.CreateContainerStatusRequest{
		IPAddress:          ipAddress,
		PingTime:           pingTime,
		LastSuccessfulPing: time.Now(),
	}
	jsonBody, _ := json.Marshal(requestBody)

	mockUseCase.On("CreateContainerStatus", mock.Anything).Return(nil, fmt.Errorf("database error"))
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockLogger.On("Errorf", mock.Anything, mock.Anything).Return()

	req := httptest.NewRequest(http.MethodPost, "/container_status", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.CreateContainerStatus(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestUpdateContainerStatus_SuccessfullyUpdatesContainer(t *testing.T) {
	mockUseCase := new(mocks.ContainerStatusUseCaseInterface)
	mockLogger := new(mocks.LoggerInterface)

	handler := handlers.NewContainerStatusHandler(mockUseCase, mockLogger)

	requestBody := pdto.UpdateContainerStatusRequest{
		PingTime: pingTime,
	}
	jsonBody, _ := json.Marshal(requestBody)

	mockUseCase.On("UpdateContainerStatus", ipAddress, mock.Anything).Return(nil)
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()

	req := httptest.NewRequest(http.MethodPatch, "/container_status/"+ipAddress, bytes.NewReader(jsonBody))
	req = mux.SetURLVars(req, map[string]string{"ip": ipAddress})
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.UpdateContainerStatus(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestUpdateContainerStatus_ErrorFromUseCase_ReturnsInternalServerError(t *testing.T) {
	mockUseCase := new(mocks.ContainerStatusUseCaseInterface)
	mockLogger := new(mocks.LoggerInterface)

	handler := handlers.NewContainerStatusHandler(mockUseCase, mockLogger)

	requestBody := pdto.UpdateContainerStatusRequest{
		PingTime: pingTime,
	}
	jsonBody, _ := json.Marshal(requestBody)

	mockUseCase.On("UpdateContainerStatus", ipAddress, mock.Anything).Return(fmt.Errorf("update failed"))
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockLogger.On("Errorf", mock.Anything, mock.Anything, mock.Anything).Return()

	req := httptest.NewRequest(http.MethodPatch, "/container_status/"+ipAddress, bytes.NewReader(jsonBody))
	req = mux.SetURLVars(req, map[string]string{"ip": ipAddress})
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.UpdateContainerStatus(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestUpdateContainerStatus_InvalidJSON_ReturnsBadRequest(t *testing.T) {
	mockUseCase := new(mocks.ContainerStatusUseCaseInterface)
	mockLogger := new(mocks.LoggerInterface)

	handler := handlers.NewContainerStatusHandler(mockUseCase, mockLogger)

	invalidJSON := []byte(`{"PingTime":}`)

	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockLogger.On("Errorf", mock.Anything, mock.Anything, mock.Anything).Return()

	req := httptest.NewRequest(http.MethodPatch, "/container_status/"+ipAddress, bytes.NewReader(invalidJSON))
	req = mux.SetURLVars(req, map[string]string{"ip": ipAddress})
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.UpdateContainerStatus(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestUpdateContainerStatus_NoFieldsProvided_ReturnsBadRequest(t *testing.T) {
	mockUseCase := new(mocks.ContainerStatusUseCaseInterface)
	mockLogger := new(mocks.LoggerInterface)

	handler := handlers.NewContainerStatusHandler(mockUseCase, mockLogger)

	emptyRequest := pdto.UpdateContainerStatusRequest{}
	jsonBody, _ := json.Marshal(emptyRequest)

	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockLogger.On("Errorf", mock.Anything, mock.Anything, mock.Anything).Return()

	req := httptest.NewRequest(http.MethodPatch, "/container_status/"+ipAddress, bytes.NewReader(jsonBody))
	req = mux.SetURLVars(req, map[string]string{"ip": ipAddress})
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.UpdateContainerStatus(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockLogger.AssertExpectations(t)
}

func TestDeleteContainerStatus_SuccessfullyDeletesContainer(t *testing.T) {
	mockUseCase := new(mocks.ContainerStatusUseCaseInterface)
	mockLogger := new(mocks.LoggerInterface)

	handler := handlers.NewContainerStatusHandler(mockUseCase, mockLogger)

	mockUseCase.On("DeleteContainerStatusByIP", ipAddress).Return(nil)
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()

	req := httptest.NewRequest(http.MethodDelete, "/container_status/"+ipAddress, http.NoBody)
	req = mux.SetURLVars(req, map[string]string{"ip": ipAddress})
	rec := httptest.NewRecorder()

	handler.DeleteContainerStatus(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestDeleteContainerStatus_ErrorFromUseCase_ReturnsInternalServerError(t *testing.T) {
	mockUseCase := new(mocks.ContainerStatusUseCaseInterface)
	mockLogger := new(mocks.LoggerInterface)

	handler := handlers.NewContainerStatusHandler(mockUseCase, mockLogger)

	mockUseCase.On("DeleteContainerStatusByIP", ipAddress).Return(fmt.Errorf("delete failed"))
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockLogger.On("Errorf", mock.Anything, mock.Anything, mock.Anything).Return()

	req := httptest.NewRequest(http.MethodDelete, "/container_status/"+ipAddress, http.NoBody)
	req = mux.SetURLVars(req, map[string]string{"ip": ipAddress})
	rec := httptest.NewRecorder()

	handler.DeleteContainerStatus(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestDeleteContainerStatus_ContainerNotFound_ReturnsNotFound(t *testing.T) {
	mockUseCase := new(mocks.ContainerStatusUseCaseInterface)
	mockLogger := new(mocks.LoggerInterface)

	handler := handlers.NewContainerStatusHandler(mockUseCase, mockLogger)

	expectedErr := fmt.Errorf("container status with IP %s not found", ipAddress)
	mockUseCase.On("DeleteContainerStatusByIP", ipAddress).Return(expectedErr)
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockLogger.On("Warnf", mock.Anything, mock.Anything).Return()

	req := httptest.NewRequest(http.MethodDelete, "/container_status/"+ipAddress, http.NoBody)
	req = mux.SetURLVars(req, map[string]string{"ip": ipAddress})
	rec := httptest.NewRecorder()

	handler.DeleteContainerStatus(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	mockUseCase.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
