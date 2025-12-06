package postgres

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DB struct {
	Connection *sqlx.DB
}

func InitDB() *DB {
	if err := godotenv.Load("../../.env"); err != nil {
		fmt.Println("Error is occurred on .env file, please check")
		panic(err)
	}

	host := os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	pass := os.Getenv("DB_PASSWORD")

	psqlSetup := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, pass,
	)

	db, err := sqlx.Open("postgres", psqlSetup)
	if err != nil {
		panic(err)
	}

	return &DB{
		Connection: db,
	}
}
