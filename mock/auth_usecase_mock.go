package mocks

import (
	"andikawhy/go-user-management/helper"
	"andikawhy/go-user-management/repository"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type AuthUsecaseMock struct {
	mock.Mock
}

func (m *AuthUsecaseMock) Login(loginData repository.Login) (string, *helper.StandardError) {
	args := m.Called()
	return args.Get(0).(string), args.Get(1).(*helper.StandardError)
}

func (m *AuthUsecaseMock) Register(registerData repository.Register) (*repository.UserResponse, *helper.StandardError) {
	args := m.Called()
	return args.Get(0).(*repository.UserResponse), args.Get(1).(*helper.StandardError)
}

func (m *AuthUsecaseMock) ValidateToken(c *gin.Context) {

}
