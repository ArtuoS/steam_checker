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
	userGameAdapters "steam_checker/internal/user/user_game/adapters"

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
	userGameRepository := userGameAdapters.NewPostgresRepository(db)

	userService := user.NewService(userRepository, gameService, userGameRepository)
	userRouter := user.NewRouter(userService)

	authorized := router.Group("/")
	authorized.Use(auth.RequireAuth)

	{
		authorized.POST("/games/track", gameRouter.Track)
	}

	{
		router.POST("/users", userRouter.Create)
		authorized.POST("/users/games/:app_id/track", userRouter.Track)

		router.POST("/auth", userRouter.Authenticate)
	}

	router.Run()
}
