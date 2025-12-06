package adapters

import (
	"context"
	"steam_checker/internal/game/event"
	"steam_checker/internal/infra/db/postgres"

	"github.com/google/uuid"
)

type PostgresRepository struct {
	DB *postgres.DB
}

func NewPostgresRepository(db *postgres.DB) *PostgresRepository {
	return &PostgresRepository{
		DB: db,
	}
}

func (p *PostgresRepository) GetByGameID(ctx context.Context, gameID uuid.UUID) ([]event.Event, error) {
	var models []event.Event
	rows, err := p.DB.Connection.QueryContext(ctx, "SELECT * FROM events WHERE game_id = ?", gameID)
	if err != nil {
		return models, err
	}

	if err = rows.Scan(&models); err != nil {
		return models, err
	}

	return models, nil
}
