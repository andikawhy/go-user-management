package main

import (
	"andikawhy/go-user-management/repository"
	"andikawhy/go-user-management/router"
	"andikawhy/go-user-management/usecase"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	loadEnvs()
	repository.ConnectDB()
}

func main() {
	ginRouter := gin.Default()

	ginRouter.GET("/", router.Root)
	ginRouter.POST("/register", router.Register)
	ginRouter.POST("/login", router.Login)
	ginRouter.GET("/users", usecase.ValidateToken, router.ListUsers)
	ginRouter.POST("/users", usecase.ValidateToken, router.CreateUser)
	ginRouter.DELETE("/users/:id", usecase.ValidateToken, router.RemoveUser)
	ginRouter.Run()
}

func loadEnvs() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
