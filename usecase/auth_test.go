package usecase_test

import (
	"andikawhy/go-user-management/helper"
	mocks "andikawhy/go-user-management/mock"
	"andikawhy/go-user-management/repository"
	"andikawhy/go-user-management/usecase"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang-jwt/jwt/v4"
)

func TestLogin(t *testing.T) {
	t.Run("test normal login", func(t *testing.T) {
		findByUsernameResponse := mockUser

		userRepositoryMock := new(mocks.UserRepositoryMock)

		userRepositoryMock.On("FindByUsername").Return(findByUsernameResponse)

		authUsecase := usecase.NewAuthUsecaseImpl(userRepositoryMock)
		loginResult, err := authUsecase.Login(repository.Login{Username: "username", Password: "password"})

		assert.Equal(t, len(loginResult) > 0, true)
		assert.Equal(t, err, nil)
	})

	t.Run("test user not found login", func(t *testing.T) {
		findByUsernameResponse := repository.User{}

		userRepositoryMock := new(mocks.UserRepositoryMock)

		userRepositoryMock.On("FindByUsername").Return(findByUsernameResponse)

		authUsecase := usecase.NewAuthUsecaseImpl(userRepositoryMock)
		loginResult, err := authUsecase.Login(repository.Login{Username: "username", Password: "password"})

		assert.Equal(t, len(loginResult) > 0, false)
		assert.Equal(t, err, helper.StandardError{Error: errors.New("user not found"), ErrorCode: 400})
	})

	t.Run("test wrong password login", func(t *testing.T) {
		findByUsernameResponse := mockUser

		userRepositoryMock := new(mocks.UserRepositoryMock)

		userRepositoryMock.On("FindByUsername").Return(findByUsernameResponse)

		authUsecase := usecase.NewAuthUsecaseImpl(userRepositoryMock)
		loginResult, err := authUsecase.Login(repository.Login{Username: "username", Password: "wrong password"})

		assert.Equal(t, len(loginResult) > 0, false)
		assert.Equal(t, err, helper.StandardError{Error: errors.New("wrong password"), ErrorCode: 401})
	})
}

func TestRegister(t *testing.T) {
	t.Run("test normal register", func(t *testing.T) {
		expectedResponse := mockUserResponse

		userRepositoryMock := new(mocks.UserRepositoryMock)

		userRepositoryMock.On("FindByUsername").Return(repository.User{})
		userRepositoryMock.On("Save").Return(mockUser)

		authUsecase := usecase.NewAuthUsecaseImpl(userRepositoryMock)
		registerResult, err := authUsecase.Register(repository.Register{Username: "username", Password: "password", Email: "test@mail.com"})

		assert.Equal(t, err, nil)
		assert.Equal(t, registerResult, expectedResponse)
	})

	t.Run("user already exist", func(t *testing.T) {
		userRepositoryMock := new(mocks.UserRepositoryMock)

		userRepositoryMock.On("FindByUsername").Return(mockUser)
		userRepositoryMock.On("Save").Return(mockUser)

		authUsecase := usecase.NewAuthUsecaseImpl(userRepositoryMock)
		registerResult, err := authUsecase.Register(repository.Register{Username: "username", Password: "password", Email: "test@mail.com"})

		assert.Equal(t, err, helper.StandardError{Error: errors.New("user already exist"), ErrorCode: http.StatusBadRequest})
		assert.Equal(t, registerResult, nil)
	})
}
func TestValidateToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	userRepositoryMock := new(mocks.UserRepositoryMock)
	authUsecase := usecase.NewAuthUsecaseImpl(userRepositoryMock)
	router.Use(authUsecase.ValidateToken)

	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	os.Setenv("SECRET", "testkey")

	t.Run("Authorization header missing", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.MatchRegex(t, w.Body.String(), "authorization header is missing")
	})

	t.Run("Invalid token format", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Add("Authorization", "invalid token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.MatchRegex(t, w.Body.String(), "invalid token format")
	})

	t.Run("Invalid token", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodHS256)
		tokenString, _ := token.SignedString([]byte("wrongkey"))

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Add("Authorization", "Bearer "+tokenString)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.MatchRegex(t, w.Body.String(), "invalid or expired token")
	})

	t.Run("Valid token and user exists", func(t *testing.T) {
		userRepositoryMock.On("FindByUsername").Return(mockUser)
		claims := jwt.MapClaims{
			"username": "validUser",
			"exp":      float64(time.Now().Add(time.Hour).Unix()),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET")))

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Add("Authorization", "Bearer "+tokenString)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
