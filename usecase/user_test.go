package usecase_test

import (
	"andikawhy/go-user-management/helper"
	mocks "andikawhy/go-user-management/mock"
	"andikawhy/go-user-management/repository"
	"andikawhy/go-user-management/usecase"
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
)

var mockUser = repository.User{
	ID:       100,
	Username: "username",
	Password: "$2a$10$i2iCoVi58g3cc144o3pb/.VQDcQajDWkmixVl7B.aNHJcoAK.MyXa",
	Email:    "test@mail.com",
}

var mockUserResponse = repository.UserResponse{
	ID:       100,
	Username: "username",
	Email:    "test@mail.com",
}

func TestListUsers(t *testing.T) {
	t.Run("test normal case list users", func(t *testing.T) {
		findAllMockResponse := []repository.User{mockUser}
		expectedResponse := []repository.UserResponse{mockUserResponse}

		userRepositoryMock := new(mocks.UserRepositoryMock)

		userRepositoryMock.On("FindAll").Return(findAllMockResponse)

		userUsecase := usecase.NewUserUsecaseImpl(userRepositoryMock)
		users, err := userUsecase.ListUsers()

		assert.Equal(t, expectedResponse, users)
		assert.Equal(t, err, nil)
	})

	t.Run("test normal case list empty users", func(t *testing.T) {
		findAllMockResponse := []repository.User{}
		expectedResponse := &[]repository.UserResponse{}

		userRepositoryMock := new(mocks.UserRepositoryMock)

		userRepositoryMock.On("FindAll").Return(findAllMockResponse)

		userUsecase := usecase.NewUserUsecaseImpl(userRepositoryMock)
		users, err := userUsecase.ListUsers()

		assert.Equal(t, expectedResponse, users)
		assert.Equal(t, err, nil)
	})
}

func TestRemoveUser(t *testing.T) {
	t.Run("test normal case remove user", func(t *testing.T) {
		deleteMockResponse := mockUser
		expectedResponse := mockUserResponse

		userRepositoryMock := new(mocks.UserRepositoryMock)

		userRepositoryMock.On("FindById").Return(deleteMockResponse)
		userRepositoryMock.On("Delete").Return(deleteMockResponse)

		userUsecase := usecase.NewUserUsecaseImpl(userRepositoryMock)
		users, err := userUsecase.RemoveUser(100, 101)

		assert.Equal(t, expectedResponse, users)
		assert.Equal(t, err, nil)
	})

	t.Run("negative: current user == deleted user", func(t *testing.T) {
		deleteMockResponse := mockUser

		userRepositoryMock := new(mocks.UserRepositoryMock)

		userRepositoryMock.On("FindById").Return(deleteMockResponse)

		userUsecase := usecase.NewUserUsecaseImpl(userRepositoryMock)
		users, err := userUsecase.RemoveUser(100, 100)

		assert.Equal(t, users, nil)
		assert.Equal(t, err, helper.StandardError{Error: errors.New("cannot delete current user"), ErrorCode: 400})
	})

	t.Run("negative: user not found", func(t *testing.T) {
		userRepositoryMock := new(mocks.UserRepositoryMock)

		userRepositoryMock.On("FindById").Return(repository.User{})

		userUsecase := usecase.NewUserUsecaseImpl(userRepositoryMock)
		users, err := userUsecase.RemoveUser(100, 101)

		assert.Equal(t, users, nil)
		assert.Equal(t, err, helper.StandardError{Error: errors.New("user not found"), ErrorCode: 400})
	})
}
