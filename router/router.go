package router

import (
	"andikawhy/go-user-management/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(userRouter UserRouter, authRouter AuthRouter, authUsecase usecase.AuthUsecase) *gin.Engine {
	ginRouter := gin.Default()

	ginRouter.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "OK")
	})
	ginRouter.POST("/api/v1/register", authRouter.Register)
	ginRouter.POST("/api/v1/login", authRouter.Login)
	ginRouter.GET("/api/v1/users", authUsecase.ValidateToken, userRouter.ListUsers)
	ginRouter.POST("/api/v1/users", authUsecase.ValidateToken, userRouter.CreateUser)
	ginRouter.DELETE("/api/v1/users/:id", authUsecase.ValidateToken, userRouter.RemoveUser)

	return ginRouter
}
