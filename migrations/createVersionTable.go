package migrations

import "github.com/jmoiron/sqlx"

const createVersionTable = `
	CREATE TABLE "version" (
	id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	application_id uuid NOT NULL,
	track_id uuid NOT NULL, 
    version text NOT NULL UNIQUE,
    url text NOT NULL,
    upload_date date NOT NULL,
    FOREIGN KEY (application_id) REFERENCES application(id),
    FOREIGN KEY (track_id) references track(id)
)
`

type CreateVersionTable struct{}

func (migration CreateVersionTable) Execute(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(createVersionTable)
	if err != nil {
		_ = tx.Rollback()

		return err
	}

	return tx.Commit()
}

func (migration CreateVersionTable) GetVersion() string {
	return "CreateVersionTable"
}
