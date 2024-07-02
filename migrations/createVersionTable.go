package migrations

import "github.com/jmoiron/sqlx"

const createVersionTable = `
create collation en_natural (
	locale = 'en-us-u-kn-true',
	provider = 'icu'
);

create table "version" (
	id uuid primary key default uuid_generate_v4(),
	application_id uuid not null,
	track_id uuid not null, 
    version text collate en_natural not null unique,
    upload_date date not null,
    foreign key (application_id) references application(id),
    foreign key (track_id) references track(id),
	unique (track_id, version)
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
