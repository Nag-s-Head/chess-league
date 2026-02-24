package db

import (
	"errors"

	"github.com/djpiper28/rpg-book/common/database/migrations"
	"github.com/jmoiron/sqlx"
)

type Db struct {
	db *sqlx.DB
}

func New() (*Db, error) {
	db, err := InternalConnect()
	if err != nil {
		return nil, err
	}

	return From(db) 
}

func From(in *sqlx.DB) (*Db, error) {
	db := &Db{db: in}
	migrator := migrations.New([]migrations.Migration{
		{
			Sql: `
-- user is a Postgres keyword this took too long to figure out
CREATE TABLE players (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	name_normalised TEXT NOT NULL UNIQUE,
	elo INTEGER DEFAULT 1000 CHECK(elo >= 0),
	join_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE games (
	player_white TEXT NOT NULL REFERENCES players(id),
	player_black TEXT NOT NULL REFERENCES players(id),
	score TEXT NOT NULL CHECK (score='1-0' OR score='0-1' OR score='1/2-1/2'),
	submitter TEXT NOT NULL REFERENCES players(id),
	played TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted BOOLEAN NOT NULL DEFAULT FALSE,
	elo_given INT NOT NULL,
	elo_taken INT NOT NULL,
	submit_ip TEXT,
	submit_user_agent TEXT
); 
			`,
		},
	})

	migrator.Rebinder = sqlx.DOLLAR
	err := migrator.Migrate(db)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot migrate database"), err)
	}

	return db, nil
}

func (d *Db) GetSqlxDb() *sqlx.DB {
	return d.db
}

func (d *Db) Close() {
	d.db.Close()
}
