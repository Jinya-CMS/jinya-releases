package migrations

import (
	"github.com/jmoiron/sqlx"
)

const createPushTokenTable = `
CREATE TABLE "push_token" (
	id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	token text NOT NULL,
	application_id uuid NOT NULL,
    FOREIGN KEY (application_id) REFERENCES application(id)
)`

type CreatePushTokenTable struct{}

func (migration CreatePushTokenTable) Execute(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(createPushTokenTable)
	if err != nil {
		_ = tx.Rollback()

		return err
	}

	return tx.Commit()
}

func (migration CreatePushTokenTable) GetVersion() string {
	return "CreatePushTokenTable"
}
