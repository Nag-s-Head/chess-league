package model

import (
	"errors"
	"time"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/djpiper28/rpg-book/common/normalisation"
	"github.com/google/uuid"
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

func InsertPlayer(db *db.Db, player Player) error {
	_, err := db.GetSqlxDb().
		NamedExec(
			"INSERT INTO players (id, name, name_normalised, elo, join_time) VALUES (:id, :name, :name_normalised, :elo, :join_time);",
			player)

	if err != nil {
		return errors.Join(errors.New("Cannot insert player"), err)
	}
	return nil
}
