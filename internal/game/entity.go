package game

import (
	"steam_checker/internal/game/event"

	"github.com/google/uuid"
)

type Game struct {
	ID       uuid.UUID     `json:"id" db:"id"`
	AppID    int           `json:"app_id" db:"app_id"`
	Name     string        `json:"name" db:"name"`
	Events   []event.Event `json:"events" db:"-"`
	Packages []int         `json:"packages" db:"packages"`
}
