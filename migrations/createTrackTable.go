package migrations

import "github.com/jmoiron/sqlx"

const createTrackTable = `
create table "track" (
	id uuid primary key default uuid_generate_v4(),
	application_id uuid not null,
    name text not null unique,
    slug text not null unique,
    is_default bool not null,
    foreign key (application_id) references application(id)
)
`

type CreateTrackTable struct{}

func (migration CreateTrackTable) Execute(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(createTrackTable)
	if err != nil {
		_ = tx.Rollback()

		return err
	}

	return tx.Commit()
}

func (migration CreateTrackTable) GetVersion() string {
	return "CreateTrackTable"
}
