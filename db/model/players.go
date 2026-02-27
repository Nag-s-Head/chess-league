package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/djpiper28/rpg-book/common/normalisation"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const StartingElo = 1000

type Player struct {
	Id             uuid.UUID `db:"id"`
	Name           string    `db:"name"`
	NameNormalised string    `db:"name_normalised"`
	Elo            int       `db:"elo"`
	JoinTime       time.Time `db:"join_time"`
}

func NewPlayer(name string) Player {
	return Player{
		Id:             uuid.New(),
		Name:           name,
		NameNormalised: normalisation.Normalise(name),
		Elo:            StartingElo,
		JoinTime:       time.Now(),
	}
}

func InsertPlayerTx(tx *sqlx.Tx, player Player) error {
	_, err := tx.
		NamedExec(
			"INSERT INTO players (id, name, name_normalised, elo, join_time) VALUES (:id, :name, :name_normalised, :elo, :join_time);",
			player)

	if err != nil {
		return errors.Join(fmt.Errorf("Cannot insert player %s", player.Name), err)
	}
	return nil
}

func InsertPlayer(db *db.Db, player Player) error {
	tx, err := db.GetSqlxDb().BeginTxx(context.Background(), nil)
	if err != nil {
		return errors.Join(errors.New("Could not start transaction"), err)
	}
	defer tx.Rollback()

	err = InsertPlayerTx(tx, player)
	if err != nil {
		return errors.Join(fmt.Errorf("Cannot insert player %s", player.Name), err)
	}

	err = tx.Commit()
	if err != nil {
		return errors.Join(errors.New("Could not commit transaction"), err)
	}
	return nil
}

func GetPlayer(db *db.Db, id uuid.UUID) (Player, error) {
	row := db.GetSqlxDb().QueryRowx(
		"SELECT * FROM players WHERE id=$1;",
		id)

	var player Player
	err := row.StructScan(&player)
	if err != nil {
		return Player{}, errors.Join(errors.New("Cannot get player"), err)
	}

	return player, nil
}

func SearchPlayerByName(db *db.Db, name string) ([]Player, error) {
	rows, err := db.GetSqlxDb().Queryx(`SELECT * FROM players WHERE name_normalised LIKE $1 ORDER BY name_normalised ASC;`, "%"+normalisation.Normalise(name)+"%")
	if err != nil {
		return nil, errors.Join(errors.New("Cannot search players by rough name"), err)
	}

	players := make([]Player, 0)
	for rows.Next() {
		var player Player
		err = rows.StructScan(&player)
		if err != nil {
			return nil, errors.Join(errors.New("Cannot struct scan player"), err)
		}

		players = append(players, player)
	}

	return players, nil
}

func GetPlayers(db *db.Db) ([]Player, error) {
	rows, err := db.GetSqlxDb().Queryx("SELECT * FROM players ORDER BY name_normalised ASC;")
	if err != nil {
		return nil, errors.Join(errors.New("Cannot get players"), err)
	}

	players := make([]Player, 0)
	for rows.Next() {
		var player Player
		err = rows.StructScan(&player)
		if err != nil {
			return nil, errors.Join(errors.New("Cannot struct scan player"), err)
		}

		players = append(players, player)
	}

	return players, nil
}

func GetPlayersByElo(db *db.Db) ([]Player, error) {
	rows, err := db.GetSqlxDb().Queryx("SELECT * FROM players ORDER BY elo DESC;")
	if err != nil {
		return nil, errors.Join(errors.New("Cannot get players"), err)
	}

	players := make([]Player, 0)
	for rows.Next() {
		var player Player
		err = rows.StructScan(&player)
		if err != nil {
			return nil, errors.Join(errors.New("Cannot struct scan player"), err)
		}

		players = append(players, player)
	}

	return players, nil
}

func getOrCreatePlayer(tx *sqlx.Tx, name string) (Player, error) {
	var player Player
	row := tx.QueryRowx("SELECT * FROM players WHERE name_normalised=$1;", normalisation.Normalise(name))

	err := row.StructScan(&player)
	if errors.Is(sql.ErrNoRows, err) {
		player = NewPlayer(name)
		err := InsertPlayerTx(tx, player)
		if err != nil {
			return Player{}, errors.Join(errors.New("Could not create player 1"), err)
		}
	} else if err != nil {
		return Player{}, errors.Join(errors.New("Could not scan player"), err)
	}

	return player, nil
}
