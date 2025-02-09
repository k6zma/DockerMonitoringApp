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
	testContainerIDStr  = "container1234567890abcdef"
	testContainerIP     = "192.168.1.101"
	testPingTimeDefault = 10.0
	testPingTimeUpdated = 20.0
)

func TestFindContainerStatuses_Success(t *testing.T) {
	mockRepo := new(mocks.ContainerStatusRepository)
	mockLogger := new(mocks.LoggerInterface)

	useCase := usecases.NewContainerStatusUseCase(mockRepo, mockLogger)

	mockFilter := &dto.ContainerStatusFilter{
		ContainerID: new(string),
	}
	*mockFilter.ContainerID = testContainerIDStr

	mockResult := []*domain.ContainerStatus{
		{
			ContainerID: testContainerIDStr,
			IPAddress:   testContainerIP,
			PingTime:    testPingTimeDefault,
		},
	}

	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockRepo.On("Find", mockFilter).Return(mockResult, nil)

	result, err := useCase.FindContainerStatuses(mockFilter)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, testContainerIP, result[0].IPAddress)
	assert.Equal(t, testContainerIDStr, result[0].ContainerID)

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
		ContainerID: testContainerIDStr,
		IPAddress:   testContainerIP,
		PingTime:    testPingTimeDefault,
	}

	mockLogger.On("Debugf", mock.Anything, mock.Anything, mock.Anything).Return()
	mockRepo.On("Create", mock.Anything).Return(nil)

	result, err := useCase.CreateContainerStatus(mockDTO)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testContainerIP, result.IPAddress)
	assert.Equal(t, testContainerIDStr, result.ContainerID)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestCreateContainerStatus_Error(t *testing.T) {
	mockRepo := new(mocks.ContainerStatusRepository)
	mockLogger := new(mocks.LoggerInterface)

	useCase := usecases.NewContainerStatusUseCase(mockRepo, mockLogger)

	mockDTO := &dto.ContainerStatusDTO{
		ContainerID: testContainerIDStr,
		IPAddress:   testContainerIP,
		PingTime:    testPingTimeDefault,
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

	mockContainerID := testContainerIDStr
	mockDTO := &dto.ContainerStatusDTO{
		PingTime:           testPingTimeUpdated,
		LastSuccessfulPing: time.Now(),
	}
	existingStatus := []*domain.ContainerStatus{
		{
			ContainerID: mockContainerID,
			IPAddress:   testContainerIP,
			PingTime:    testPingTimeDefault,
		},
	}

	mockLogger.On("Debugf", mock.Anything, mock.Anything, mock.Anything).Return()
	mockRepo.On("Find", &dto.ContainerStatusFilter{ContainerID: &mockContainerID}).Return(existingStatus, nil)
	mockRepo.On("Update", mock.Anything).Return(nil)

	err := useCase.UpdateContainerStatus(mockContainerID, mockDTO)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestUpdateContainerStatus_ErrorFetching(t *testing.T) {
	mockRepo := new(mocks.ContainerStatusRepository)
	mockLogger := new(mocks.LoggerInterface)

	useCase := usecases.NewContainerStatusUseCase(mockRepo, mockLogger)

	mockContainerID := testContainerIDStr
	mockDTO := &dto.ContainerStatusDTO{PingTime: testPingTimeUpdated}

	mockLogger.On("Debugf", mock.Anything, mock.Anything, mock.Anything).Return()
	mockRepo.On("Find", &dto.ContainerStatusFilter{ContainerID: &mockContainerID}).
		Return(nil, fmt.Errorf("database error"))
	mockLogger.On("Errorf", mock.Anything, mock.Anything, mock.Anything).Return()

	err := useCase.UpdateContainerStatus(mockContainerID, mockDTO)

	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestUpdateContainerStatus_NotFound(t *testing.T) {
	mockRepo := new(mocks.ContainerStatusRepository)
	mockLogger := new(mocks.LoggerInterface)

	useCase := usecases.NewContainerStatusUseCase(mockRepo, mockLogger)

	mockContainerID := testContainerIDStr
	mockDTO := &dto.ContainerStatusDTO{PingTime: testPingTimeUpdated}

	mockLogger.On("Debugf", "USECASES: updating container status for container ID: %s with data: %+v", mockContainerID, mock.Anything).
		Return()
	mockRepo.On("Find", &dto.ContainerStatusFilter{ContainerID: &mockContainerID}).
		Return([]*domain.ContainerStatus{}, nil)
	mockLogger.On("Errorf", "USECASES: error fetching container status with container ID %s not found", mockContainerID).
		Return()

	err := useCase.UpdateContainerStatus(mockContainerID, mockDTO)

	assert.Error(t, err)
	assert.Equal(t, fmt.Errorf("container status with container ID %s not found", mockContainerID), err)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestUpdateContainerStatus_UpdateError(t *testing.T) {
	mockRepo := new(mocks.ContainerStatusRepository)
	mockLogger := new(mocks.LoggerInterface)

	useCase := usecases.NewContainerStatusUseCase(mockRepo, mockLogger)

	mockContainerID := testContainerIDStr
	mockDTO := &dto.ContainerStatusDTO{PingTime: testPingTimeUpdated}
	existingStatus := []*domain.ContainerStatus{
		{
			ContainerID: mockContainerID,
			IPAddress:   testContainerIP,
			PingTime:    testPingTimeDefault,
		},
	}

	mockLogger.On("Debugf", "USECASES: updating container status for container ID: %s with data: %+v", mockContainerID, mock.Anything).
		Return()
	mockRepo.On("Find", &dto.ContainerStatusFilter{ContainerID: &mockContainerID}).Return(existingStatus, nil)
	mockRepo.On("Update", mock.Anything).Return(fmt.Errorf("update failed"))
	mockLogger.On("Errorf", "USECASES: failed to update container status for container ID %s: %v", mockContainerID, mock.Anything).
		Return()

	err := useCase.UpdateContainerStatus(mockContainerID, mockDTO)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "failed to update container status: update failed")

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestDeleteContainerStatusByContainerID_ErrorFetching(t *testing.T) {
	mockRepo := new(mocks.ContainerStatusRepository)
	mockLogger := new(mocks.LoggerInterface)

	useCase := usecases.NewContainerStatusUseCase(mockRepo, mockLogger)

	mockContainerID := testContainerIDStr

	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockRepo.On("Find", &dto.ContainerStatusFilter{ContainerID: &mockContainerID}).
		Return(nil, fmt.Errorf("database error"))
	mockLogger.On("Errorf", mock.Anything, mock.Anything, mock.Anything).Return()

	err := useCase.DeleteContainerStatusByContainerID(mockContainerID)

	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestDeleteContainerStatusByContainerID_NotFound(t *testing.T) {
	mockRepo := new(mocks.ContainerStatusRepository)
	mockLogger := new(mocks.LoggerInterface)

	useCase := usecases.NewContainerStatusUseCase(mockRepo, mockLogger)

	mockContainerID := testContainerIDStr

	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockRepo.On("Find", &dto.ContainerStatusFilter{ContainerID: &mockContainerID}).
		Return([]*domain.ContainerStatus{}, nil)
	mockLogger.On("Warnf", mock.Anything, mock.Anything, mock.Anything).Return()

	err := useCase.DeleteContainerStatusByContainerID(mockContainerID)

	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestDeleteContainerStatusByContainerID_ErrorDeleting(t *testing.T) {
	mockRepo := new(mocks.ContainerStatusRepository)
	mockLogger := new(mocks.LoggerInterface)

	useCase := usecases.NewContainerStatusUseCase(mockRepo, mockLogger)

	mockContainerID := testContainerIDStr
	existingStatus := []*domain.ContainerStatus{
		{
			ContainerID: testContainerIDStr,
			IPAddress:   testContainerIP,
			PingTime:    testPingTimeDefault,
		},
	}

	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockRepo.On("Find", &dto.ContainerStatusFilter{ContainerID: &mockContainerID}).Return(existingStatus, nil)
	mockRepo.On("DeleteByContainerID", mockContainerID).Return(fmt.Errorf("delete failed"))
	mockLogger.On("Errorf", mock.Anything, mock.Anything, mock.Anything).Return()

	err := useCase.DeleteContainerStatusByContainerID(mockContainerID)

	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestDeleteContainerStatusByContainerID_Success(t *testing.T) {
	mockRepo := new(mocks.ContainerStatusRepository)
	mockLogger := new(mocks.LoggerInterface)

	useCase := usecases.NewContainerStatusUseCase(mockRepo, mockLogger)

	mockContainerID := testContainerIDStr
	existingStatus := []*domain.ContainerStatus{
		{
			ContainerID: testContainerIDStr,
			IPAddress:   testContainerIP,
			PingTime:    testPingTimeDefault,
		},
	}

	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
	mockRepo.On("Find", &dto.ContainerStatusFilter{ContainerID: &mockContainerID}).Return(existingStatus, nil)
	mockRepo.On("DeleteByContainerID", mockContainerID).Return(nil)
	mockLogger.On("Debugf", "USECASES: successfully deleted container status for container_id: %s", mockContainerID).
		Return()

	err := useCase.DeleteContainerStatusByContainerID(mockContainerID)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}
