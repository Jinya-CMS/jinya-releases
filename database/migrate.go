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

		// Introduced with JRWEB-4
		_, err = conn.Exec(`
create or replace function add_foreign_key_if_not_exists(from_table text, from_column text, to_table text, to_column text)
returns void language plpgsql as
$$
declare 
   fk_exists boolean;
begin
    fk_exists := case when exists (select true
	from information_schema.table_constraints tc
		inner join information_schema.constraint_column_usage ccu
			using (constraint_catalog, constraint_schema, constraint_name)
		inner join information_schema.key_column_usage kcu
			using (constraint_catalog, constraint_schema, constraint_name)
	where constraint_type = 'FOREIGN KEY'
	  and ccu.table_name = to_table
	  and ccu.column_name = to_column
	  and tc.table_name = from_table
	  and kcu.column_name = from_column) then true else false end;
	if not fk_exists then
		execute format('alter table %s add constraint %s_%s_fkey foreign key (%s) references %s(%s) on delete cascade', from_table, from_table, to_column, from_column, to_table, to_column);
	end if;
end
$$;
`)
		if err != nil {
			panic(err)
		}

		// Use function added in JRWEB-4 to create FK
		_, err = conn.Exec(`
select add_foreign_key_if_not_exists('track', 'application_id', 'application', 'id');
select add_foreign_key_if_not_exists('version', 'track_id', 'track', 'id');
select add_foreign_key_if_not_exists('version', 'application_id', 'application', 'id');
select add_foreign_key_if_not_exists('push_token', 'application_id', 'application', 'id');
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
