package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/djpiper28/rpg-book/common/normalisation"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const StartingElo = 1000
const (
	StartingLiglicko2Rating     = 1500.0
	StartingLiglicko2Deviation  = 500.0
	StartingLiglicko2Volatility = 0.09
)

type Player struct {
	Id             uuid.UUID `db:"id"`
	Name           string    `db:"name"`
	NameNormalised string    `db:"name_normalised"`
	DEPRECATEDElo  int       `db:"elo"` // Deprecated: for use with old elo system
	// Liglicko2Rating is the player's current liglicko2 rating scalar.
	Liglicko2Rating float64 `db:"liglicko2_rating"`
	// Liglicko2Deviation is the player's current liglicko2 rating deviation (RD).
	Liglicko2Deviation float64 `db:"liglicko2_deviation"`
	// Liglicko2Volatility is the player's current liglicko2 volatility (sigma).
	Liglicko2Volatility float64 `db:"liglicko2_volatility"`
	// Liglicko2At is the liglicko2 timestamp used by the algorithm.
	// It is stored as "rating periods since Unix epoch", where this app currently
	// defines 1 rating period as 1 day.
	Liglicko2At float64   `db:"liglicko2_at"`
	JoinTime    time.Time `db:"join_time"`
	Deleted     bool      `db:"deleted"`
}

type PlayerWithGameCount struct {
	Player
	GameCount int `db:"game_count"`
}

func NewPlayer(name string) Player {
	return Player{
		Id:                  uuid.New(),
		Name:                name,
		NameNormalised:      normalisation.Normalise(name),
		DEPRECATEDElo:       StartingElo,
		Liglicko2Rating:     StartingLiglicko2Rating,
		Liglicko2Deviation:  StartingLiglicko2Deviation,
		Liglicko2Volatility: StartingLiglicko2Volatility,
		Liglicko2At:         liglicko2InstantFromTime(time.Now()),
		JoinTime:            time.Now(),
		Deleted:             false,
	}
}

func InsertPlayerTx(tx *sqlx.Tx, player Player) error {
	_, err := tx.
		NamedExec(
			`INSERT INTO players (id, name, name_normalised, elo, liglicko2_rating, liglicko2_deviation, liglicko2_volatility, liglicko2_at, join_time)
VALUES (:id, :name, :name_normalised, :elo, :liglicko2_rating, :liglicko2_deviation, :liglicko2_volatility, :liglicko2_at, :join_time);`,
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

func GetPlayersByEloWithGameCount(db *db.Db) ([]PlayerWithGameCount, error) {
	rows, err := db.GetSqlxDb().Queryx(`
		SELECT 
		  players.*, COUNT(games.ikey) AS game_count FROM players 
		LEFT JOIN 
		    games 
		  ON 
		    games.player_white=players.id OR games.player_black=players.id
		GROUP by players.id
		ORDER BY liglicko2_rating DESC, name;`)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot get players"), err)
	}

	players := make([]PlayerWithGameCount, 0)
	for rows.Next() {
		var player PlayerWithGameCount
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

func GetTotalPlayerCount(db *db.Db) (int, error) {
	var count int
	err := db.GetSqlxDb().Get(&count, "SELECT count(*) FROM players WHERE deleted=false")
	if err != nil {
		return 0, errors.Join(errors.New("Cannot get total player count"), err)
	}

	return count, nil
}

func RenamePlayer(db *db.Db, id uuid.UUID, newName string, adminId uuid.UUID) error {
	tx, err := db.GetSqlxDb().BeginTxx(context.Background(), nil)
	if err != nil {
		return errors.Join(errors.New("Cannot start transaction"), err)
	}
	defer tx.Rollback()

	var oldName string
	err = tx.Get(&oldName, "SELECT name FROM players WHERE id=$1;", id)
	if err != nil {
		return errors.Join(errors.New("Cannot get old player name"), err)
	}

	_, err = tx.
		Exec("UPDATE players SET name=$1, name_normalised=$2 WHERE id=$3;",
			newName,
			normalisation.Normalise(newName),
			id)
	if err != nil {
		return errors.Join(errors.New("Cannot update player"), err)
	}

	auditLog := NewAuditLog(adminId, "Player rename", fmt.Sprintf("Renamed from '%s' to '%s'", oldName, newName))
	err = InsertAuditLog(tx, auditLog)
	if err != nil {
		return errors.Join(errors.New("Cannot insert audit log"), err)
	}

	err = InsertAuditLogPlayerAffected(tx, NewAuditLogPlayerAffected(auditLog.Id, id, 0))
	if err != nil {
		return errors.Join(errors.New("Cannot insert audit log player affected"), err)
	}

	err = tx.Commit()
	if err != nil {
		return errors.Join(errors.New("Cannot commit transaction"), err)
	}

	slog.Info("Player renamed", "oldName", oldName, "newName", newName, "by", adminId)

	return nil
}
