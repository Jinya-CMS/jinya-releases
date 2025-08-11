package database

import (
	"jinya-releases/config"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

func Connect() (*sqlx.DB, error) {
	connectionString := config.LoadedConfiguration.PostgresUrl
	db, err := sqlx.Connect("pgx", connectionString)

	if err != nil {
		return nil, err
	}

	return db, err
}
