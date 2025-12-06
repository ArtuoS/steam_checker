package main

import (
	"steam_checker/internal/game"
	gameAdapters "steam_checker/internal/game/adapters"
	eventAdapters "steam_checker/internal/game/event/adapters"
	"steam_checker/internal/infra/db/postgres"
	"steam_checker/internal/infra/integration/steam"
	"steam_checker/internal/infra/middleware/auth"
	"steam_checker/internal/user"

	userAdapters "steam_checker/internal/user/adapters"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(gin.Recovery())

	db := postgres.InitDB()

	gameRepository := gameAdapters.NewPostgresRepository(db)
	eventRepository := eventAdapters.NewPostgresRepository(db)
	steamIntegration := steam.NewIntegration()

	gameService := game.NewService(gameRepository, eventRepository, steamIntegration)
	gameRouter := game.NewRouter(gameService)

	userRepository := userAdapters.NewPostgresRepository(db)

	userService := user.NewService(userRepository)
	userRouter := user.NewRouter(userService)

	authorized := router.Group("/")
	authorized.Use(auth.RequireAuth)

	{
		authorized.POST("/games/track", gameRouter.Track)
	}

	{
		router.POST("/users", userRouter.Create)
		router.POST("/auth", userRouter.Authenticate)
	}

	router.Run()
}
