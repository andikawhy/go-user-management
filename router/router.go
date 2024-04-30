package router

import (
	"andikawhy/go-user-management/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userController *UserController, authController *AuthController, authUsecase usecase.AuthUsecase) *gin.Engine {
	ginRouter := gin.Default()

	ginRouter.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "{\"status\":\"OK\"}")
	})
	ginRouter.POST("/api/v1/register", authController.Register)
	ginRouter.POST("/api/v1/login", authController.Login)
	ginRouter.GET("/api/v1/users", authUsecase.ValidateToken, userController.ListUsers)
	ginRouter.POST("/api/v1/users", authUsecase.ValidateToken, userController.CreateUser)
	ginRouter.DELETE("/api/v1/users/:id", authUsecase.ValidateToken, userController.RemoveUser)

	return ginRouter
}
