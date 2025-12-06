package main

import (
	"steam_checker/internal/game"
	gameAdapters "steam_checker/internal/game/adapters"
	eventAdapters "steam_checker/internal/game/event/adapters"
	"steam_checker/internal/infra/db/postgres"
	"steam_checker/internal/infra/integration/steam"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	db := postgres.InitDB()

	gameRepository := gameAdapters.NewPostgresRepository(db)
	eventRepository := eventAdapters.NewPostgresRepository(db)
	steamIntegration := steam.NewIntegration()

	gameService := game.NewService(gameRepository, eventRepository, steamIntegration)
	gameRouter := game.NewRouter(gameService)

	{
		router.POST("/games/track", gameRouter.Track)
	}

	router.Run()
}
