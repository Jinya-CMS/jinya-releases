package database

import (
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"jinya-releases/config"
)

func Connect() (*sqlx.DB, error) {
	connectionString := config.LoadedConfiguration.PostgresUrl
	db, err := sqlx.Connect("pgx", connectionString)

	if err != nil {
		return nil, err
	}

	return db, err
}
