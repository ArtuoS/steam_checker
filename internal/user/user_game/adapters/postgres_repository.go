package adapters

import (
	"context"
	"steam_checker/internal/infra/db/postgres"
	"steam_checker/internal/user/user_game"
)

type PostgresRepository struct {
	DB *postgres.DB
}

func NewPostgresRepository(db *postgres.DB) *PostgresRepository {
	return &PostgresRepository{
		DB: db,
	}
}

func (p *PostgresRepository) Create(ctx context.Context, userGame *user_game.UserGame) error {
	_, err := p.DB.Connection.Exec(ctx,
		`INSERT INTO user_games (id, user_id, game_id) VALUES ($1, $2, $3) 
			ON CONFLICT DO NOTHING`,
		userGame.ID, userGame.UserID, userGame.GameID)
	if err != nil {
		return err
	}

	return nil
}
