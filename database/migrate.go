package database

import (
	"context"
	"jinya-releases/config"

	"github.com/DerKnerd/gorp"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

var dbMap *gorp.DbMap

func GetDbMap() *gorp.DbMap {
	return dbMap
}

func SetupDatabase() {
	if dbMap == nil {
		pool, err := pgxpool.New(context.Background(), config.LoadedConfiguration.PostgresUrl)
		if err != nil {
			panic(err)
		}

		conn := stdlib.OpenDBFromPool(pool)

		dialect := gorp.PostgresDialect{}

		dbMap = &gorp.DbMap{Db: conn, Dialect: dialect}

		AddTableWithName[Application]("application")
		AddTableWithName[PushToken]("push_token")
		track := AddTableWithName[Track]("track")
		track.SetUniqueTogether("slug", "application_id")
		track.SetUniqueTogether("name", "application_id")
		version := AddTableWithName[Version]("version")
		version.SetUniqueTogether("version", "application_id", "track_id")

		err = dbMap.CreateTablesIfNotExists()
		if err != nil {
			panic(err)
		}

		_, err = conn.Exec(`
alter table track
	drop constraint if exists track_application_id_fkey;
alter table version
	drop constraint if exists version_track_id_fkey;
alter table version
	drop constraint if exists version_application_id_fkey;
alter table push_token
	drop constraint if exists push_token_application_id_fkey;
`)
		if err != nil {
			panic(err)
		}

		_, err = conn.Exec(`
alter table track
	add constraint track_application_id_fkey foreign key (application_id) references application(id); 
alter table version
	add constraint version_track_id_fkey foreign key (track_id) references track(id);
alter table version
	add constraint version_application_id_fkey foreign key (application_id) references application(id);
alter table push_token
	add constraint push_token_application_id_fkey foreign key (application_id) references application(id);
`)
		if err != nil {
			panic(err)
		}

		_, err = conn.Exec(`
create extension if not exists "uuid-ossp";

alter table application 
    alter column id set default uuid_generate_v4();
alter table track 
    alter column id set default uuid_generate_v4();
alter table version 
    alter column id set default uuid_generate_v4();
alter table push_token 
    alter column id set default uuid_generate_v4();
`)
		if err != nil {
			panic(err)
		}

		_, err = conn.Exec(`
drop table if exists migrations;
alter table version 
    drop constraint if exists version_version_key;
alter table version 
	drop column if exists url;
`)
		if err != nil {
			panic(err)
		}
	}
}
