package postgres

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DB struct {
	Connection *pgxpool.Pool
}

func InitDB() *DB {
	if err := godotenv.Load("../../.env"); err != nil {
		fmt.Println("Error is occurred on .env file, please check")
		panic(err)
	}

	config, err := pgxpool.ParseConfig("")
	if err != nil {
		panic(fmt.Errorf("failed to parse pgx config: %w", err))
	}

	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	config.ConnConfig.Host = os.Getenv("DB_HOST")
	config.ConnConfig.Port = uint16(port)
	config.ConnConfig.User = os.Getenv("DB_USER")
	config.ConnConfig.Password = os.Getenv("DB_PASSWORD")
	config.ConnConfig.Database = os.Getenv("DB_NAME")

	// Recommended pool settings
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		panic(fmt.Errorf("unable to connect to database: %w", err))
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		panic(fmt.Errorf("database ping failed: %w", err))
	}

	fmt.Println("Connected to PostgreSQL successfully with pgx!")
	return &DB{Connection: pool}
}
