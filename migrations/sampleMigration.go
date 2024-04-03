package migrations

import "github.com/jmoiron/sqlx"

type SampleMigration struct {
}

func (s SampleMigration) Execute(db *sqlx.DB) error {
	panic("implement me")
}

func (s SampleMigration) GetVersion() string {
	panic("implement me")
}
