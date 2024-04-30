package router_test

import (
	"andikawhy/go-user-management/helper"
	mocks "andikawhy/go-user-management/mock"
	"andikawhy/go-user-management/router"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockUserUsecase := new(mocks.UserUsecaseMock)
		mockAuthUsecase := new(mocks.AuthUsecaseMock)
		authRouter := router.NewAuthRouter(mockUserUsecase, mockAuthUsecase)

		mockError := &helper.StandardError{Error: nil, ErrorCode: http.StatusOK}

		mockAuthUsecase.On("Register").Return(&mockUser, mockError)

		router := gin.Default()
		router.POST("/register", authRouter.Register)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"username": "username", "password": "password", "email": "test@mail.com"}`))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.MatchRegex(t, w.Body.String(), "successfully register")
	})

	t.Run("Error", func(t *testing.T) {
		mockUserUsecase := new(mocks.UserUsecaseMock)
		mockAuthUsecase := new(mocks.AuthUsecaseMock)
		authRouter := router.NewAuthRouter(mockUserUsecase, mockAuthUsecase)

		mockError := &helper.StandardError{Error: errors.New("error message"), ErrorCode: http.StatusInternalServerError}

		mockAuthUsecase.On("Register").Return(&mockUser, mockError)

		router := gin.Default()
		router.POST("/register", authRouter.Register)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"username": "username", "password": "password", "email": "test@mail.com"}`))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.MatchRegex(t, w.Body.String(), "error message")
	})

	t.Run("Bind JSON Error", func(t *testing.T) {
		authRouter := router.NewAuthRouter(nil, nil)

		router := gin.Default()
		router.POST("/register", authRouter.Register)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/register", strings.NewReader(``))

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockUserUsecase := new(mocks.UserUsecaseMock)
		mockAuthUsecase := new(mocks.AuthUsecaseMock)
		authRouter := router.NewAuthRouter(mockUserUsecase, mockAuthUsecase)

		mockError := &helper.StandardError{Error: nil, ErrorCode: http.StatusOK}

		mockAuthUsecase.On("Login").Return(string("token"), mockError)

		router := gin.Default()
		router.POST("/login", authRouter.Login)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username": "username", "password": "password"}`))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.MatchRegex(t, w.Body.String(), "successfully login")
	})

	t.Run("Error", func(t *testing.T) {
		mockUserUsecase := new(mocks.UserUsecaseMock)
		mockAuthUsecase := new(mocks.AuthUsecaseMock)
		authRouter := router.NewAuthRouter(mockUserUsecase, mockAuthUsecase)

		mockError := &helper.StandardError{Error: errors.New("error message"), ErrorCode: http.StatusInternalServerError}

		mockAuthUsecase.On("Login").Return(string("token"), mockError)

		router := gin.Default()
		router.POST("/login", authRouter.Login)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username": "username", "password": "password"}`))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.MatchRegex(t, w.Body.String(), "error message")
	})

	t.Run("Bind JSON Error", func(t *testing.T) {
		authRouter := router.NewAuthRouter(nil, nil)

		router := gin.Default()
		router.POST("/login", authRouter.Login)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(``))

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

}
