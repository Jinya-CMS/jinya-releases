package migrations

import (
	"github.com/jmoiron/sqlx"
)

const createApplicationTable = `
CREATE TABLE "application" (
	id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name text NOT NULL,
    logo text NULL,
    slug text UNIQUE NOT NULL,
    homepage_template text NULL,
    trackpage_template text NULL,
    additional_css text NULL,
    additional_javascript text NULL
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
