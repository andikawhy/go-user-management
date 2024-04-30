package usecase

import (
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

func Register(registerData repository.Register) (*repository.User, error) {
	var userFound repository.User
	repository.DB.Where("username=?", registerData.Username).Find(&userFound)

	if userFound.ID != 0 {
		return nil, errors.New("user already exist")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(registerData.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &repository.User{
		Username: registerData.Username,
		Email:    registerData.Email,
		Password: string(passwordHash),
	}

	if err := repository.DB.Create(&user).Error; err != nil {
		return nil, errors.New("create user db failure")
	}

	return user, nil
}

func Login(loginData repository.Login) (string, error) {
	var userFound repository.User
	repository.DB.Where("username=?", loginData.Username).Find(&userFound)

	if userFound.ID == 0 {
		return "", errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(loginData.Password)); err != nil {
		return "", errors.New("invalid password")
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

func ValidateToken(c *gin.Context) {
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

	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var user repository.User
	repository.DB.Where("ID=?", claims["id"]).Find(&user)

	if user.ID == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("currentUserId", user.ID)

	c.Next()
}
