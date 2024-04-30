package router

import (
	"andikawhy/go-user-management/repository"
	"andikawhy/go-user-management/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var registerData repository.Register

	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := usecase.Register(registerData)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func Login(c *gin.Context) {
	var loginData repository.Login

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := usecase.Login(loginData)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
