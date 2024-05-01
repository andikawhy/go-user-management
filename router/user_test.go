package router_test

import (
	"andikawhy/go-user-management/helper"
	mocks "andikawhy/go-user-management/mock"
	"andikawhy/go-user-management/repository"
	"andikawhy/go-user-management/router"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

var mockUser = repository.UserResponse{
	ID:       100,
	Username: "username",
	Email:    "test@mail.com",
}

func TestRemoveUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockUserUsecase := new(mocks.UserUsecaseMock)
		userRouter := router.NewUserRouterImpl(mockUserUsecase, nil)

		mockError := &helper.StandardError{Error: nil, ErrorCode: http.StatusOK}

		mockUserUsecase.On("RemoveUser").Return(&mockUser, mockError)

		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set("currentUserId", uint64(2))
			c.Next()
		})
		router.DELETE("/users/:id", userRouter.RemoveUser)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/users/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.MatchRegex(t, w.Body.String(), "successfully remove user")
	})

	t.Run("Current user id parsing failure", func(t *testing.T) {
		mockUserUsecase := new(mocks.UserUsecaseMock)
		userRouter := router.NewUserRouterImpl(mockUserUsecase, nil)

		mockError := &helper.StandardError{Error: nil, ErrorCode: http.StatusOK}

		mockUserUsecase.On("RemoveUser").Return(&mockUser, mockError)

		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set("currentUserId", "abc")
			c.Next()
		})
		router.DELETE("/users/:id", userRouter.RemoveUser)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/users/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.MatchRegex(t, w.Body.String(), "Failed to convert current user ID")
	})

	t.Run("Requested User ID Conversion Fail", func(t *testing.T) {
		userRouter := router.NewUserRouterImpl(nil, nil)

		router := gin.Default()
		router.DELETE("/users/:id", userRouter.RemoveUser)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/users/abc", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.MatchRegex(t, w.Body.String(), "Failed to convert requested user ID")
	})

	t.Run("Current user id context not found", func(t *testing.T) {
		mockUserUsecase := new(mocks.UserUsecaseMock)
		userRouter := router.NewUserRouterImpl(mockUserUsecase, nil)

		mockError := &helper.StandardError{Error: nil, ErrorCode: http.StatusOK}

		mockUserUsecase.On("RemoveUser").Return(&mockUser, mockError)

		router := gin.Default()
		router.DELETE("/users/:id", userRouter.RemoveUser)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/users/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.MatchRegex(t, w.Body.String(), "current user not found")
	})

	t.Run("Error from use case", func(t *testing.T) {
		mockUserUsecase := new(mocks.UserUsecaseMock)
		userRouter := router.NewUserRouterImpl(mockUserUsecase, nil)

		mockError := &helper.StandardError{Error: errors.New("error message"), ErrorCode: http.StatusInternalServerError}

		mockUserUsecase.On("RemoveUser").Return(&mockUser, mockError)

		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set("currentUserId", uint64(2))
			c.Next()
		})
		router.DELETE("/users/:id", userRouter.RemoveUser)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/users/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.MatchRegex(t, w.Body.String(), "error message")
		mockUserUsecase.AssertExpectations(t)
	})
}

func TestListUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockUserUsecase := new(mocks.UserUsecaseMock)
		userRouter := router.NewUserRouterImpl(mockUserUsecase, nil)

		mockError := &helper.StandardError{Error: nil, ErrorCode: http.StatusOK}

		mockUserUsecase.On("ListUsers").Return(&[]repository.UserResponse{mockUser}, mockError)

		router := gin.Default()
		router.GET("/users", userRouter.ListUsers)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/users", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.MatchRegex(t, w.Body.String(), "successfully list users")
	})

	t.Run("Error", func(t *testing.T) {
		mockUserUsecase := new(mocks.UserUsecaseMock)
		userRouter := router.NewUserRouterImpl(mockUserUsecase, nil)

		mockError := &helper.StandardError{Error: errors.New("error message"), ErrorCode: http.StatusInternalServerError}

		mockUserUsecase.On("ListUsers").Return(&[]repository.UserResponse{mockUser}, mockError)

		router := gin.Default()
		router.GET("/users", userRouter.ListUsers)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/users", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.MatchRegex(t, w.Body.String(), "error message")
	})

}

func TestCreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockUserUsecase := new(mocks.UserUsecaseMock)
		mockAuthUsecase := new(mocks.AuthUsecaseMock)
		userRouter := router.NewUserRouterImpl(mockUserUsecase, mockAuthUsecase)

		mockError := &helper.StandardError{Error: nil, ErrorCode: http.StatusOK}

		mockAuthUsecase.On("Register").Return(&mockUser, mockError)

		router := gin.Default()
		router.POST("/users", userRouter.CreateUser)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"username": "username", "password": "password", "email": "test@mail.com"}`))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.MatchRegex(t, w.Body.String(), "successfully create user")
	})

	t.Run("Error", func(t *testing.T) {
		mockUserUsecase := new(mocks.UserUsecaseMock)
		mockAuthUsecase := new(mocks.AuthUsecaseMock)
		userRouter := router.NewUserRouterImpl(mockUserUsecase, mockAuthUsecase)

		mockError := &helper.StandardError{Error: errors.New("error message"), ErrorCode: http.StatusInternalServerError}

		mockAuthUsecase.On("Register").Return(&mockUser, mockError)

		router := gin.Default()
		router.POST("/users", userRouter.CreateUser)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"username": "username", "password": "password", "email": "test@mail.com"}`))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.MatchRegex(t, w.Body.String(), "error message")
	})

	t.Run("Bind JSON Error", func(t *testing.T) {
		userRouter := router.NewUserRouterImpl(nil, nil)

		router := gin.Default()
		router.POST("/users", userRouter.CreateUser)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader(``))

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

}
