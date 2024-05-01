package router_test

import (
	mocks "andikawhy/go-user-management/mock"
	"andikawhy/go-user-management/router"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/mock"
)

func TestSetupRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authRouterMock := new(mocks.AuthRouterMock)
	userRouterMock := new(mocks.UserRouterMock)
	authUsecaseMock := new(mocks.AuthUsecaseMock)

	authRouterMock.On("Register", mock.Anything)
	authRouterMock.On("Login", mock.Anything)
	userRouterMock.On("ListUsers", mock.Anything)
	userRouterMock.On("CreateUser", mock.Anything)
	userRouterMock.On("RemoveUser", mock.Anything)
	authUsecaseMock.On("ValidateToken", mock.Anything)

	router := router.SetupRouter(userRouterMock, authRouterMock, authUsecaseMock)

	t.Run("GET /", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "\"OK\"", w.Body.String())
	})

	t.Run("POST /api/v1/register", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := bytes.NewBufferString(`{"username":"testuser","password":"testpass"}`)
		req, _ := http.NewRequest("POST", "/api/v1/register", body)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("POST /api/v1/login", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := bytes.NewBufferString(`{"username":"testuser","password":"testpass"}`)
		req, _ := http.NewRequest("POST", "/api/v1/login", body)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GET /api/v1/users", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/users", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("POST /api/v1/users", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := bytes.NewBufferString(`{"name":"newuser","email":"newuser@test.com"}`)
		req, _ := http.NewRequest("POST", "/api/v1/users", body)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("DELETE /api/v1/users/:id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/v1/users/123", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
