package mocks

import (
	"andikawhy/go-user-management/helper"
	"andikawhy/go-user-management/repository"

	"github.com/stretchr/testify/mock"
)

type UserUsecaseMock struct {
	mock.Mock
}

func (m *UserUsecaseMock) RemoveUser(deletedUserID uint64, currentUserId uint64) (*repository.UserResponse, *helper.StandardError) {
	args := m.Called()
	return args.Get(0).(*repository.UserResponse), args.Get(1).(*helper.StandardError)
}

func (m *UserUsecaseMock) ListUsers() (*[]repository.UserResponse, *helper.StandardError) {
	args := m.Called()
	return args.Get(0).(*[]repository.UserResponse), args.Get(1).(*helper.StandardError)
}
