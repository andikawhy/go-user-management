package main

import (
	"andikawhy/go-user-management/repository"
	"andikawhy/go-user-management/router"
	"andikawhy/go-user-management/usecase"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	loadEnvs()
	db := repository.ConnectDB()

	userRepository := repository.NewUserRepositoryImpl(db)

	userUsecase := usecase.NewUserUsecaseImpl(userRepository)
	authUsecase := usecase.NewAuthUsecaseImpl(userRepository)

	userRouter := router.NewUserRouterImpl(userUsecase, authUsecase)
	authRouter := router.NewAuthRouterImpl(userUsecase, authUsecase)

	ginRouter := router.SetupRouter(userRouter, authRouter, authUsecase)
	ginRouter.Run()
}

func loadEnvs() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
