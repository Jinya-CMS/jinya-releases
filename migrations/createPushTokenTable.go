package migrations

import (
	"github.com/jmoiron/sqlx"
)

const createPushTokenTable = `
create table "push_token" (
	id uuid primary key default uuid_generate_v4(),
	token text not null,
	application_id uuid not null,
    foreign key (application_id) references application(id)
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
