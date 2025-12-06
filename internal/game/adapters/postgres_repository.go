package adapters

import (
	"context"
	"fmt"
	"steam_checker/internal/game"
	"steam_checker/internal/infra/db/postgres"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
)

type PostgresRepository struct {
	DB *postgres.DB
}

func NewPostgresRepository(db *postgres.DB) *PostgresRepository {
	return &PostgresRepository{
		DB: db,
	}
}

func (p *PostgresRepository) GetAll(ctx context.Context) ([]game.Game, error) {
	var models []game.Game
	rows, err := p.DB.Connection.Query(ctx, `SELECT * FROM games`)
	if err != nil {
		return models, err
	}

	if err = rows.Scan(&models); err != nil {
		return models, err
	}

	return models, nil
}

func (p *PostgresRepository) Create(ctx context.Context, input *game.Game) error {
	tx, err := p.DB.Connection.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	if _, err := tx.Exec(ctx,
		`INSERT INTO games (id, app_id, name, packages) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING`,
		input.ID, input.AppID, input.Name, pq.Array(input.Packages)); err != nil {
		return err
	}

	if len(input.Events) > 0 {
		valueArgs := make([]interface{}, 0, len(input.Events)*4)
		valuePlaceholders := make([]string, 0, len(input.Events))

		for i, evt := range input.Events {
			base := i * 4
			valuePlaceholders = append(valuePlaceholders,
				fmt.Sprintf("($%d, $%d, $%d, $%d)", base+1, base+2, base+3, base+4),
			)
			valueArgs = append(valueArgs, evt.ID, input.ID, evt.Type, evt.Data)
		}

		query := fmt.Sprintf(`
            INSERT INTO events (id, game_id, type, data) VALUES %s ON CONFLICT (id) DO NOTHING`,
			strings.Join(valuePlaceholders, ","))

		if _, err := tx.Exec(ctx, query, valueArgs...); err != nil {
			return err
		}
	}

	return nil
}

func (p *PostgresRepository) Exists(ctx context.Context, appID int) (bool, error) {
	rows, err := p.DB.Connection.Query(ctx,
		`SELECT EXISTS(SELECT 1 FROM games WHERE app_id = $1)`, appID)
	if err != nil {
		return false, err
	}

	return rows.Next(), nil
}
