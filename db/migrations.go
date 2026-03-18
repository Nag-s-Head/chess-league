package db

import (
	"errors"
	"strings"
	"unicode"

	"github.com/djpiper28/rpg-book/common/database/migrations"
	"github.com/google/uuid"
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

func InternalFixPlayerNameCapitals(name string) string {
	parts := strings.Split(name, " ")
	for i, str := range parts {
		strBytes := []byte(str)
		if len(str) > 0 {
			strBytes[0] = byte(unicode.ToUpper(rune(str[0])))
		}
		parts[i] = string(strBytes)
	}

	return strings.Join(parts, " ")
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
		{
			Sql: `
CREATE INDEX idx_players_name_norm ON players(name_normalised);
CREATE INDEX idx_players_elo ON players(elo);
			`,
		},
		{
			Sql: `
CREATE SEQUENCE game_ikey_sequence AS BIGINT;
			`,
		},
		{
			Sql: `
ALTER TABLE games ADD COLUMN ikey BIGINT NOT NULL UNIQUE;
			`,
		},
		{
			Sql: `
ALTER TABLE players ADD COLUMN deleted BOOL NOT NULL DEFAULT false;
			`,
		},
		{
			PreProcess: func(tx *sqlx.Tx) error {
				type Player struct {
					Id   uuid.UUID `db:"id"`
					Name string    `db:"name"`
				}

				rows, err := tx.Queryx("SELECT id, name FROM players;")
				if err != nil {
					return errors.Join(errors.New("Cannot get players"), err)
				}

				players := make([]Player, 0)
				for rows.Next() {
					var player Player
					err := rows.StructScan(&player)
					if err != nil {
						return errors.Join(errors.New("Cannot scan player"), err)
					}

					players = append(players, player)
				}

				for _, player := range players {
					player.Name = InternalFixPlayerNameCapitals(player.Name)

					_, err := tx.NamedExec("UPDATE players SET name=:name WHERE id=:id;", player)
					if err != nil {
						return errors.Join(errors.New("Cannot update player name"), err)
					}
				}

				return nil
			},
		},
		{
			Sql: `
CREATE INDEX idx_games_player_white ON games(player_white);
CREATE INDEX idx_games_player_black ON games(player_black);
CREATE INDEX idx_games_played ON games(played);
			`,
		},
		{
			Sql: `
CREATE TABLE admin_users (
  id UUID PRIMARY KEY,
	name TEXT NOT NULL,
	oauth_id TEXT NOT NULL UNIQUE,
	created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	session_key TEXT,
	last_login TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	last_ip TEXT,
	last_user_agent TEXT
);
			`,
		},
		{
			Sql: "CREATE INDEX idx_admin_users_session_key ON admin_users(session_key);",
		},
		{
			Sql: `
CREATE TABLE audit_logs (
	id UUID PRIMARY KEY,
	created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	done_by UUID NOT NULL REFERENCES admin_users(id),
	operation_name TEXT NOT NULL,
	operation_description TEXT NOT NULL
);

CREATE TABLE audit_log_player_affected (
	audit_log_id UUID NOT NULL REFERENCES audit_logs(id),
	player_id TEXT NOT NULL REFERENCES players(id),
	elo_change INT NOT NULL,
	UNIQUE(audit_log_id, player_id)
);

CREATE INDEX idx_audit_log_player_affected_audit_log_id ON audit_log_player_affected(audit_log_id);
CREATE INDEX idx_audit_log_player_affected_player_id ON audit_log_player_affected(player_id);


CREATE TABLE audit_log_game_affected (
	audit_log_id UUID NOT NULL REFERENCES audit_logs(id),
	game_ikey BIGINT NOT NULL REFERENCES games(ikey),
	UNIQUE(audit_log_id, game_ikey)
);

CREATE INDEX idx_audit_log_game_affected_audit_log_id ON audit_log_game_affected(audit_log_id);
CREATE INDEX idx_audit_log_game_affected_game_ikey ON audit_log_game_affected(game_ikey);
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
