package router

import (
	"andikawhy/go-user-management/repository"
	"andikawhy/go-user-management/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthRouter interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

type AuthRouterImpl struct {
	userUsecase usecase.UserUsecase
	authUsecase usecase.AuthUsecase
}

func NewAuthRouterImpl(userUsecase usecase.UserUsecase, authUsecase usecase.AuthUsecase) AuthRouter {
	return &AuthRouterImpl{
		userUsecase: userUsecase,
		authUsecase: authUsecase,
	}
}

func (t *AuthRouterImpl) Register(c *gin.Context) {
	var registerData repository.Register

	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, registerError := t.authUsecase.Register(registerData)

	if registerError.Error != nil {
		c.JSON(int(registerError.ErrorCode), gin.H{"error": registerError.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user, "message": "successfully register"})
}

func (t *AuthRouterImpl) Login(c *gin.Context) {
	var loginData repository.Login

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, loginError := t.authUsecase.Login(loginData)

	if loginError.Error != nil {
		c.JSON(int(loginError.ErrorCode), gin.H{"error": loginError.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "message": "successfully login"})
}
