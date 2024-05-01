package usecase

import (
	"andikawhy/go-user-management/helper"
	"andikawhy/go-user-management/repository"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Login(loginData repository.Login) (string, *helper.StandardError)
	Register(registerData repository.Register) (*repository.UserResponse, *helper.StandardError)
	ValidateToken(c *gin.Context)
}

type AuthUsecaseImpl struct {
	UserRepository repository.UserRepository
}

func (t *AuthUsecaseImpl) Register(registerData repository.Register) (*repository.UserResponse, *helper.StandardError) {
	userFound := t.UserRepository.FindByUsername(registerData.Username)

	if userFound.ID != 0 {
		return nil, &helper.StandardError{Error: errors.New("user already exist"), ErrorCode: http.StatusBadRequest}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(registerData.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, &helper.StandardError{Error: err, ErrorCode: http.StatusInternalServerError}
	}

	user := repository.User{
		Username: registerData.Username,
		Email:    registerData.Email,
		Password: string(passwordHash),
	}

	createdUser := t.UserRepository.Save(user)

	userResponse := repository.UserResponse{
		ID:        createdUser.ID,
		Email:     createdUser.Email,
		Username:  createdUser.Username,
		CreatedAt: createdUser.CreatedAt,
	}

	return &userResponse, nil
}

func (t *AuthUsecaseImpl) Login(loginData repository.Login) (string, *helper.StandardError) {
	userFound := t.UserRepository.FindByUsername(loginData.Username)

	if userFound.ID == 0 {
		return "", &helper.StandardError{Error: errors.New("user not found"), ErrorCode: http.StatusBadRequest}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(loginData.Password)); err != nil {
		return "", &helper.StandardError{Error: errors.New("wrong password"), ErrorCode: http.StatusUnauthorized}
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       userFound.ID,
		"username": userFound.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", &helper.StandardError{Error: errors.New("failed to generate token"), ErrorCode: http.StatusInternalServerError}
	}

	return token, nil
}

func (t *AuthUsecaseImpl) ValidateToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is missing"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	authToken := strings.Split(authHeader, " ")
	if len(authToken) != 2 || authToken[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tokenString := authToken[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		c.Abort()
		return
	}

	user := t.UserRepository.FindByUsername(claims["username"].(string))

	if user == (repository.User{}) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("currentUserId", user.ID)

	c.Next()
}

func NewAuthUsecaseImpl(userRepository repository.UserRepository) AuthUsecase {
	return &AuthUsecaseImpl{
		UserRepository: userRepository,
	}
}
