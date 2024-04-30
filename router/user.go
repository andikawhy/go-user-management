package router

import (
	"andikawhy/go-user-management/repository"
	"andikawhy/go-user-management/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	userUsecase usecase.UserUsecase
	authUsecase usecase.AuthUsecase
}

func NewUserRouter(userUsecase usecase.UserUsecase, authUsecase usecase.AuthUsecase) UserRouter {
	return UserRouter{
		userUsecase: userUsecase,
		authUsecase: authUsecase,
	}
}

func (t *UserRouter) CreateUser(c *gin.Context) {
	var createUserData repository.Register

	if err := c.ShouldBindJSON(&createUserData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, registerError := t.authUsecase.Register(createUserData)

	if registerError.Error != nil {
		c.JSON(int(registerError.ErrorCode), gin.H{"error": registerError.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user, "message": "successfully create user"})
}

func (t *UserRouter) RemoveUser(c *gin.Context) {
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

	user, removeUserError := t.userUsecase.RemoveUser(userIDInt, currentUserIdInt)

	if removeUserError.Error != nil {
		c.JSON(int(removeUserError.ErrorCode), gin.H{"error": removeUserError.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user, "message": "successfully remove user"})
}

func (t *UserRouter) ListUsers(c *gin.Context) {
	users, err := t.userUsecase.ListUsers()

	if err.Error != nil {
		c.JSON(int(err.ErrorCode), gin.H{"error": err.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users, "message": "successfully list users"})
}
