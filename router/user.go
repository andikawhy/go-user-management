package router

import (
	"andikawhy/go-user-management/repository"
	"andikawhy/go-user-management/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var createUserData repository.Register

	if err := c.ShouldBindJSON(&createUserData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := usecase.Register(createUserData)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func RemoveUser(c *gin.Context) {
	userId := c.Param("id")

	userIDInt, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert requested user ID"})
		return
	}

	currentUserId, exists := c.Get("currentUserId")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "current user not found"})
		return
	}

	currentUserIdInt, ok := currentUserId.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert current user ID"})
		return
	}

	user, err := usecase.RemoveUser(userIDInt, currentUserIdInt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func ListUsers(c *gin.Context) {
	users, err := usecase.ListUsers()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}
