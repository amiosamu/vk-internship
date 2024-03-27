package postgres

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

type DataSource struct {
	DB *sqlx.DB
}

func InitDB() (*DataSource, error) {
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")

	db := os.Getenv("PG_DB")

	ssl := os.Getenv("PG_SSL")

	pgConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, db, ssl)

	open, err := sqlx.Open("postgres", pgConnString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database %w", err)
	}
	if err := open.Ping(); err != nil {
		return nil, fmt.Errorf("error ting to the database: %w", err)

	}
	return &DataSource{
		DB: open,
	}, nil
}

func (d *DataSource) Close() error {
	if err := d.DB.Close(); err != nil {
		return fmt.Errorf("error closing PostgreSQL: %w", err)
	}
	return nil
}
