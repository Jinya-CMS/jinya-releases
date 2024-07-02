package migrations

import (
	"github.com/jmoiron/sqlx"
)

const createApplicationTable = `
create table "application" (
	id uuid primary key default uuid_generate_v4(),
    name text not null unique,
    slug text not null unique,
    logo text null
)`

type CreateApplicationTable struct{}

func (migration CreateApplicationTable) Execute(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(createApplicationTable)
	if err != nil {
		_ = tx.Rollback()

		return err
	}

	return tx.Commit()
}

func (migration CreateApplicationTable) GetVersion() string {
	return "CreateApplicationTable"
}
