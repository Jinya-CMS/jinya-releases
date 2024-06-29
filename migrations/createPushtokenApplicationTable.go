package migrations

import (
	"github.com/jmoiron/sqlx"
)

const createPushtokenApplicationTable = `
CREATE TABLE "pushtoken_application" (
	token_id uuid,
	application_id uuid,
	PRIMARY KEY (token_id, application_id)
)`

type CreatePushtokenApplicationTable struct{}

func (migration CreatePushtokenApplicationTable) Execute(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(createPushtokenApplicationTable)
	if err != nil {
		_ = tx.Rollback()

		return err
	}

	return tx.Commit()
}

func (migration CreatePushtokenApplicationTable) GetVersion() string {
	return "CreatePushtokenApplicationTable"
}
