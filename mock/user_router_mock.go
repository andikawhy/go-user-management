package mocks

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type UserRouterMock struct {
	mock.Mock
}

func (m *UserRouterMock) ListUsers(c *gin.Context) {
	m.Called(c)
	c.JSON(http.StatusOK, gin.H{"status": "users listed"})
}

func (m *UserRouterMock) CreateUser(c *gin.Context) {
	m.Called(c)
	c.JSON(http.StatusOK, gin.H{"status": "user created"})
}

func (m *UserRouterMock) RemoveUser(c *gin.Context) {
	m.Called(c)
	c.JSON(http.StatusOK, gin.H{"status": "user removed"})
}
