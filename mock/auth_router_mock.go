package mocks

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type AuthRouterMock struct {
	mock.Mock
}

func (m *AuthRouterMock) Register(c *gin.Context) {
	m.Called(c)
	c.JSON(http.StatusOK, gin.H{"status": "registered"})
}

func (m *AuthRouterMock) Login(c *gin.Context) {
	m.Called(c)
	c.JSON(http.StatusOK, gin.H{"status": "logged in"})
}
