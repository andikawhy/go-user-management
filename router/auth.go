package router

import (
	"andikawhy/go-user-management/repository"
	"andikawhy/go-user-management/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	userUsecase usecase.UserUsecase
	authUsecase usecase.AuthUsecase
}

func NewAuthController(userUsecase usecase.UserUsecase, authUsecase usecase.AuthUsecase) *AuthController {
	return &AuthController{
		userUsecase: userUsecase,
		authUsecase: authUsecase,
	}
}

func (t *AuthController) Register(c *gin.Context) {
	var registerData repository.Register

	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, registerError := t.authUsecase.Register(registerData)

	if registerError != nil {
		c.JSON(int(registerError.ErrorCode), gin.H{"error": registerError.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user, "message": "successfully register"})
}

func (t *AuthController) Login(c *gin.Context) {
	var loginData repository.Login

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, loginError := t.authUsecase.Login(loginData)

	if loginError != nil {
		c.JSON(int(loginError.ErrorCode), gin.H{"error": loginError.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "message": "successfully login"})
}
