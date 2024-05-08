package migrations

import "github.com/jmoiron/sqlx"

const createTrackTable = `
	CREATE TABLE "track" (
	id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	application_id uuid NOT NULL,
    name text NOT NULL UNIQUE,
    slug text NOT NULL UNIQUE,
    is_default bool NOT NULL,
    FOREIGN KEY (application_id) REFERENCES application(id)
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
