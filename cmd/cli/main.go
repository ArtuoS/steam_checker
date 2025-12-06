package main

import (
	"context"
	"fmt"
	"steam_checker/internal/game"
	gameAdapters "steam_checker/internal/game/adapters"
	eventAdapters "steam_checker/internal/game/event/adapters"
	"steam_checker/internal/infra/db/postgres"
	"steam_checker/internal/infra/integration/steam"
)

func main() {
	db := postgres.InitDB()

	gameRepository := gameAdapters.NewPostgresRepository(db)
	eventRepository := eventAdapters.NewPostgresRepository(db)
	steamIntegration := steam.NewIntegration()

	gameService := game.NewService(gameRepository, eventRepository, steamIntegration)
	err := gameService.Create(context.Background(), &game.CreateInput{
		AppID: 1904480,
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
