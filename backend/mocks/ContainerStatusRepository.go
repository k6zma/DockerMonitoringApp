// Code generated by mockery v2.47.0. DO NOT EDIT.

package mocks

import (
	dto "github.com/k6zma/DockerMonitoringApp/backend/internal/application/dto"
	domain "github.com/k6zma/DockerMonitoringApp/backend/internal/domain"

	mock "github.com/stretchr/testify/mock"
)

// ContainerStatusRepository is an autogenerated mock type for the ContainerStatusRepository type
type ContainerStatusRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: status
func (_m *ContainerStatusRepository) Create(status *domain.ContainerStatus) error {
	ret := _m.Called(status)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.ContainerStatus) error); ok {
		r0 = rf(status)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteByContainerID provides a mock function with given fields: containerID
func (_m *ContainerStatusRepository) DeleteByContainerID(containerID string) error {
	ret := _m.Called(containerID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteByContainerID")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(containerID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Find provides a mock function with given fields: filter
func (_m *ContainerStatusRepository) Find(filter *dto.ContainerStatusFilter) ([]*domain.ContainerStatus, error) {
	ret := _m.Called(filter)

	if len(ret) == 0 {
		panic("no return value specified for Find")
	}

	var r0 []*domain.ContainerStatus
	var r1 error
	if rf, ok := ret.Get(0).(func(*dto.ContainerStatusFilter) ([]*domain.ContainerStatus, error)); ok {
		return rf(filter)
	}
	if rf, ok := ret.Get(0).(func(*dto.ContainerStatusFilter) []*domain.ContainerStatus); ok {
		r0 = rf(filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.ContainerStatus)
		}
	}

	if rf, ok := ret.Get(1).(func(*dto.ContainerStatusFilter) error); ok {
		r1 = rf(filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: status
func (_m *ContainerStatusRepository) Update(status *domain.ContainerStatus) error {
	ret := _m.Called(status)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.ContainerStatus) error); ok {
		r0 = rf(status)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewContainerStatusRepository creates a new instance of ContainerStatusRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewContainerStatusRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ContainerStatusRepository {
	mock := &ContainerStatusRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
