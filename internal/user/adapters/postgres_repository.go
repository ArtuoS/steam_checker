package adapters

import (
	"context"
	"steam_checker/internal/infra/db/postgres"
	"steam_checker/internal/user"
)

type PostgresRepository struct {
	DB *postgres.DB
}

func NewPostgresRepository(db *postgres.DB) *PostgresRepository {
	return &PostgresRepository{
		DB: db,
	}
}

func (p *PostgresRepository) GetAll(ctx context.Context) ([]user.User, error) {
	var models []user.User
	rows, err := p.DB.Connection.Query(ctx, `SELECT * FROM users`)
	if err != nil {
		return models, err
	}

	if err = rows.Scan(&models); err != nil {
		return models, err
	}

	return models, nil
}

func (p *PostgresRepository) Create(ctx context.Context, input *user.User) error {
	if _, err := p.DB.Connection.Exec(ctx,
		`INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING`,
		input.ID, input.Name, input.Email, input.Password); err != nil {
		return err
	}

	return nil
}

func (p *PostgresRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	rows, err := p.DB.Connection.Query(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, email)
	if err != nil {
		return false, err
	}

	return rows.Next(), nil
}

func (p *PostgresRepository) GetByEmail(ctx context.Context, email string) (user.User, error) {
	var model user.User
	row := p.DB.Connection.QueryRow(ctx,
		`SELECT id, name, email, password FROM users WHERE email = $1`, email) // Especifique as colunas
	if err := row.Scan(&model.ID, &model.Name, &model.Email, &model.Password); err != nil {
		return model, err
	}

	return model, nil
}
