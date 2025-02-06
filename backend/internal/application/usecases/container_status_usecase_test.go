package usecases_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/k6zma/DockerMonitoringApp/backend/internal/application/dto"
	"github.com/k6zma/DockerMonitoringApp/backend/internal/application/usecases"
	"github.com/k6zma/DockerMonitoringApp/backend/internal/domain"
	"github.com/k6zma/DockerMonitoringApp/backend/mocks"
)

const (
	testContainerID     = 1
	testContainerIP     = "192.168.1.101"
	testPingTimeDefault = 10
	testPingTimeUpdated = 20
)

func TestFindContainerStatuses_Success(t *testing.T) {
	mockRepo := new(mocks.ContainerStatusRepository)
	mockLogger := new(mocks.LoggerInterface)

	useCase := usecases.NewContainerStatusUseCase(mockRepo, mockLogger)

	mockFilter := &dto.ContainerStatusFilter{}
	mockResult := []*domain.ContainerStatus{
		{ID: testContainerID, IPAddress: testContainerIP, PingTime: testPingTimeDefault},
	}

	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockRepo.On("Find", mockFilter).Return(mockResult, nil)

	result, err := useCase.FindContainerStatuses(mockFilter)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, testContainerIP, result[0].IPAddress)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestFindContainerStatuses_Error(t *testing.T) {
	mockRepo := new(mocks.ContainerStatusRepository)
	mockLogger := new(mocks.LoggerInterface)

	useCase := usecases.NewContainerStatusUseCase(mockRepo, mockLogger)

	mockFilter := &dto.ContainerStatusFilter{}

	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockRepo.On("Find", mockFilter).Return(nil, fmt.Errorf("database error"))
	mockLogger.On("Errorf", mock.Anything, mock.Anything, mock.Anything).Return()

	result, err := useCase.FindContainerStatuses(mockFilter)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestCreateContainerStatus_Success(t *testing.T) {
	mockRepo := new(mocks.ContainerStatusRepository)
	mockLogger := new(mocks.LoggerInterface)

	useCase := usecases.NewContainerStatusUseCase(mockRepo, mockLogger)

	mockDTO := &dto.ContainerStatusDTO{
		IPAddress: testContainerIP,
		PingTime:  testPingTimeDefault,
	}

	mockLogger.On("Debugf", mock.Anything, mock.Anything, mock.Anything).Return()
	mockRepo.On("Create", mock.Anything).Return(nil)

	result, err := useCase.CreateContainerStatus(mockDTO)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testContainerIP, result.IPAddress)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestCreateContainerStatus_Error(t *testing.T) {
	mockRepo := new(mocks.ContainerStatusRepository)
	mockLogger := new(mocks.LoggerInterface)

	useCase := usecases.NewContainerStatusUseCase(mockRepo, mockLogger)

	mockDTO := &dto.ContainerStatusDTO{
		IPAddress: testContainerIP,
		PingTime:  testPingTimeDefault,
	}

	mockLogger.On("Debugf", mock.Anything, mock.Anything, mock.Anything).Return()
	mockRepo.On("Create", mock.Anything).Return(fmt.Errorf("failed to insert"))
	mockLogger.On("Errorf", mock.Anything, mock.Anything, mock.Anything).Return()

	result, err := useCase.CreateContainerStatus(mockDTO)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestUpdateContainerStatus_Success(t *testing.T) {
	mockRepo := new(mocks.ContainerStatusRepository)
	mockLogger := new(mocks.LoggerInterface)

	useCase := usecases.NewContainerStatusUseCase(mockRepo, mockLogger)

	mockIP := testContainerIP
	mockDTO := &dto.ContainerStatusDTO{
		PingTime:           testPingTimeUpdated,
		LastSuccessfulPing: time.Now(),
	}
	existingStatus := []*domain.ContainerStatus{
		{ID: testContainerID, IPAddress: mockIP, PingTime: testPingTimeDefault},
	}

	mockLogger.On("Debugf", mock.Anything, mock.Anything, mock.Anything).Return()
	mockRepo.On("Find", &dto.ContainerStatusFilter{IPAddress: &mockIP}).Return(existingStatus, nil)
	mockRepo.On("Update", mock.Anything).Return(nil)

	err := useCase.UpdateContainerStatus(mockIP, mockDTO)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestUpdateContainerStatus_ErrorFetching(t *testing.T) {
	mockRepo := new(mocks.ContainerStatusRepository)
	mockLogger := new(mocks.LoggerInterface)

	useCase := usecases.NewContainerStatusUseCase(mockRepo, mockLogger)

	mockIP := testContainerIP
	mockDTO := &dto.ContainerStatusDTO{PingTime: testPingTimeUpdated}

	mockLogger.On("Debugf", mock.Anything, mock.Anything, mock.Anything).Return()
	mockRepo.On("Find", &dto.ContainerStatusFilter{IPAddress: &mockIP}).Return(nil, fmt.Errorf("database error"))
	mockLogger.On("Errorf", mock.Anything, mock.Anything, mock.Anything).Return()

	err := useCase.UpdateContainerStatus(mockIP, mockDTO)

	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestDeleteContainerStatusByIP_ErrorFetching(t *testing.T) {
	mockRepo := new(mocks.ContainerStatusRepository)
	mockLogger := new(mocks.LoggerInterface)

	useCase := usecases.NewContainerStatusUseCase(mockRepo, mockLogger)

	mockIP := testContainerIP

	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockRepo.On("Find", &dto.ContainerStatusFilter{IPAddress: &mockIP}).Return(nil, fmt.Errorf("database error"))
	mockLogger.On("Errorf", mock.Anything, mock.Anything, mock.Anything).Return()

	err := useCase.DeleteContainerStatusByIP(mockIP)

	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestDeleteContainerStatusByIP_NotFound(t *testing.T) {
	mockRepo := new(mocks.ContainerStatusRepository)
	mockLogger := new(mocks.LoggerInterface)

	useCase := usecases.NewContainerStatusUseCase(mockRepo, mockLogger)

	mockIP := testContainerIP

	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockRepo.On("Find", &dto.ContainerStatusFilter{IPAddress: &mockIP}).Return([]*domain.ContainerStatus{}, nil)
	mockLogger.On("Warnf", mock.Anything, mock.Anything, mock.Anything).Return()

	err := useCase.DeleteContainerStatusByIP(mockIP)

	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestDeleteContainerStatusByIP_ErrorDeleting(t *testing.T) {
	mockRepo := new(mocks.ContainerStatusRepository)
	mockLogger := new(mocks.LoggerInterface)

	useCase := usecases.NewContainerStatusUseCase(mockRepo, mockLogger)

	mockIP := testContainerIP
	existingStatus := []*domain.ContainerStatus{
		{ID: testContainerID, IPAddress: testContainerIP, PingTime: testPingTimeDefault},
	}

	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockRepo.On("Find", &dto.ContainerStatusFilter{IPAddress: &mockIP}).Return(existingStatus, nil)
	mockRepo.On("DeleteByIP", mockIP).Return(fmt.Errorf("delete failed"))
	mockLogger.On("Errorf", mock.Anything, mock.Anything, mock.Anything).Return()

	err := useCase.DeleteContainerStatusByIP(mockIP)

	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestDeleteContainerStatusByIP_Success(t *testing.T) {
	mockRepo := new(mocks.ContainerStatusRepository)
	mockLogger := new(mocks.LoggerInterface)

	useCase := usecases.NewContainerStatusUseCase(mockRepo, mockLogger)

	mockIP := testContainerIP
	existingStatus := []*domain.ContainerStatus{
		{ID: testContainerID, IPAddress: testContainerIP, PingTime: testPingTimeDefault},
	}

	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockRepo.On("Find", &dto.ContainerStatusFilter{IPAddress: &mockIP}).Return(existingStatus, nil)
	mockRepo.On("DeleteByIP", mockIP).Return(nil)
	mockLogger.On("Debugf", "USECASES: successfully deleted container status for IP: %s", mockIP).Return()

	err := useCase.DeleteContainerStatusByIP(mockIP)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestUpdateContainerStatus_NotFound(t *testing.T) {
	mockRepo := new(mocks.ContainerStatusRepository)
	mockLogger := new(mocks.LoggerInterface)

	useCase := usecases.NewContainerStatusUseCase(mockRepo, mockLogger)

	mockIP := testContainerIP
	mockDTO := &dto.ContainerStatusDTO{PingTime: testPingTimeUpdated}

	mockLogger.On("Debugf", "USECASES: updating container status for IP: %s with data: %+v", mockIP, mock.Anything).Return()
	mockRepo.On("Find", &dto.ContainerStatusFilter{IPAddress: &mockIP}).Return([]*domain.ContainerStatus{}, nil)
	mockLogger.On("Errorf", "USECASES: error fetching container status with IP %s not found", mockIP).Return()

	err := useCase.UpdateContainerStatus(mockIP, mockDTO)

	assert.Error(t, err)
	assert.Equal(t, fmt.Errorf("container status with IP %s not found", mockIP), err)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestUpdateContainerStatus_UpdateError(t *testing.T) {
	mockRepo := new(mocks.ContainerStatusRepository)
	mockLogger := new(mocks.LoggerInterface)

	useCase := usecases.NewContainerStatusUseCase(mockRepo, mockLogger)

	mockIP := testContainerIP
	mockDTO := &dto.ContainerStatusDTO{PingTime: testPingTimeUpdated}
	existingStatus := []*domain.ContainerStatus{
		{ID: testContainerID, IPAddress: mockIP, PingTime: testPingTimeDefault},
	}

	mockLogger.On("Debugf", "USECASES: updating container status for IP: %s with data: %+v", mockIP, mock.Anything).Return()
	mockRepo.On("Find", &dto.ContainerStatusFilter{IPAddress: &mockIP}).Return(existingStatus, nil)
	mockRepo.On("Update", mock.Anything).Return(fmt.Errorf("update failed"))
	mockLogger.On("Errorf", "USECASES: failed to update container status for IP %s: %v", mockIP, mock.Anything).Return()

	err := useCase.UpdateContainerStatus(mockIP, mockDTO)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "failed to update container status: update failed")

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
