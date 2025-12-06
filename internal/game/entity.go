package game

import (
	"steam_checker/internal/game/event"
	"time"

	"github.com/google/uuid"
)

type Game struct {
	ID        uuid.UUID     `json:"id" db:"id"`
	AppID     int           `json:"app_id" db:"app_id"`
	Name      string        `json:"name" db:"name"`
	Events    []event.Event `json:"events" db:"-"`
	Packages  []int         `json:"packages" db:"packages"`
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time    `json:"deleted_at" db:"deleted_at"`
}
