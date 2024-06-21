package migrations

import (
	"github.com/jmoiron/sqlx"
)

const createPushtokenTable = `
CREATE TABLE "pushtoken" (
	id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	token uuid DEFAULT uuid_generate_v4()
	
)`

type CreatePushtokenTable struct{}

func (migration CreatePushtokenTable) Execute(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(createPushtokenTable)
	if err != nil {
		_ = tx.Rollback()

		return err
	}

	return tx.Commit()
}

func (migration CreatePushtokenTable) GetVersion() string {
	return "CreatePushtokenTable"
}
