package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/liglicko2"
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

func (p *Player) ApplyRating(rating liglicko2.Rating) {
	p.Liglicko2Rating = rating.Rating
	p.Liglicko2Deviation = rating.Deviation
	p.Liglicko2Volatility = rating.Volatility
	p.Liglicko2At = rating.At
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
		Liglicko2At:         Liglicko2InstantFromTime(time.Now()),
		JoinTime:            time.Now(),
		Deleted:             false,
	}
}

func InsertPlayerTx(tx *sqlx.Tx, player Player) error {
	_, err := tx.
		NamedExec(
			`INSERT INTO players (id, name, name_normalised, elo, liglicko2_rating, liglicko2_deviation, liglicko2_volatility, liglicko2_at, join_time, deleted)
VALUES (:id, :name, :name_normalised, :elo, :liglicko2_rating, :liglicko2_deviation, :liglicko2_volatility, :liglicko2_at, :join_time, :deleted);`,
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

func GetPlayerTx(tx *sqlx.Tx, id uuid.UUID) (Player, error) {
	row := tx.QueryRowx(
		"SELECT * FROM players WHERE id=$1;",
		id)

	var player Player
	err := row.StructScan(&player)
	if err != nil {
		return Player{}, errors.Join(errors.New("Cannot get player"), err)
	}

	return player, nil
}

func GetPlayer(db *db.Db, id uuid.UUID) (Player, error) {
	var returnPlayer Player
	err := db.DoTx(func(tx *sqlx.Tx) error {
		player, err := GetPlayerTx(tx, id)
		if err != nil {
			return err
		}

		returnPlayer = player
		return nil
	})

	if err != nil {
		return Player{}, errors.Join(errors.New("Cannot get player"), err)
	}

	return returnPlayer, nil
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

func getPlayerById(txx *sqlx.Tx, id uuid.UUID) (Player, error) {
	var player Player
	err := txx.Get(&player, "SELECT * FROM players WHERE id = $1;", id)
	if err != nil {
		return Player{}, errors.Join(errors.New("Cannot get player"), err)
	}

	return player, nil
}

func GetPlayersByElo(db *db.Db, showDeleted bool) ([]Player, error) {
	rows, err := db.GetSqlxDb().Queryx("SELECT * FROM players WHERE deleted=FALSE OR deleted=$1 ORDER BY liglicko2_rating DESC, name;", showDeleted)
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
		    (
		      games.player_white=players.id OR games.player_black=players.id
		    )
		      AND 
		    games.deleted = false
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

	var oldPlayer Player
	err = tx.Get(&oldPlayer, "SELECT * FROM players WHERE id=$1;", id)
	if err != nil {
		return errors.Join(errors.New("Cannot get old player name"), err)
	}

	nameNormalised := normalisation.Normalise(newName)

	if oldPlayer.Deleted {
		nameNormalised = oldPlayer.Id.String()
	}

	_, err = tx.
		Exec("UPDATE players SET name=$1, name_normalised=$2 WHERE id=$3;",
			newName,
			nameNormalised,
			id)
	if err != nil {
		return errors.Join(errors.New("Cannot update player"), err)
	}

	auditLog := NewAuditLog(adminId, "Player rename", fmt.Sprintf("Renamed from '%s' to '%s'", oldPlayer.Name, newName))
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

	slog.Info("Player renamed", "oldPlayer", oldPlayer, "newName", newName, "by", adminId)

	return nil
}

func DeletePlayer(db *db.Db, playerId, adminId uuid.UUID) error {
	tx, err := db.GetSqlxDb().BeginTxx(context.Background(), nil)
	if err != nil {
		return errors.Join(errors.New("Cannot start transaction"), err)
	}
	defer tx.Rollback()

	var player Player
	err = tx.Get(&player, "SELECT * FROM players WHERE id=$1", playerId)
	if err != nil {
		return errors.Join(errors.New("Cannot get player to run validation against"), err)
	}

	if player.Deleted {
		return errors.New("Cannot delete a player who has already been deleted")
	}

	name := player.Id.String()
	_, err = tx.Exec("UPDATE players SET deleted=TRUE, name=$1, name_normalised=$2 WHERE id=$3", name, normalisation.Normalise(name), playerId)
	if err != nil {
		return errors.Join(errors.New("Cannot deleted player"), err)
	}

	auditLog := NewAuditLog(adminId, "Player deletion", fmt.Sprintf("Deleted player %s", player.Name))
	err = InsertAuditLog(tx, auditLog)
	if err != nil {
		return errors.Join(errors.New("Cannot insret audit log"), err)
	}

	err = InsertAuditLogPlayerAffected(tx, NewAuditLogPlayerAffected(auditLog.Id, playerId, 0))
	if err != nil {
		return errors.Join(errors.New("Cannot insret audit log player affected"), err)
	}

	err = tx.Commit()
	if err != nil {
		return errors.Join(errors.New("Cannot commit transaction"), err)
	}
	return nil
}

func MergePlayers(db *db.Db, target, dest, adminId uuid.UUID) error {
	tx, err := db.GetSqlxDb().BeginTxx(context.Background(), nil)
	if err != nil {
		return errors.Join(errors.New("Cannot start transaction"), err)
	}
	defer tx.Rollback()

	var targetPlayer Player
	err = tx.Get(&targetPlayer, "SELECT * FROM players WHERE id=$1;", target)
	if err != nil {
		return errors.Join(errors.New("Cannot get player target name for audit logs"), err)
	}

	if targetPlayer.Deleted {
		return errors.New("Cannot merge as target player is deleted")
	}

	var destPlayer Player
	err = tx.Get(&destPlayer, "SELECT * FROM players WHERE id=$1", dest)
	if err != nil {
		return errors.Join(errors.New("Cannot get player dest name for audit logs"), err)
	}

	if destPlayer.Deleted {
		return errors.New("Cannot merge as destination player is deleted")
	}

	auditLog := NewAuditLog(adminId, "Player merger", fmt.Sprintf("Merging player %s (%s) into %s (%s).", target, targetPlayer.Name, dest, destPlayer.Name))
	err = InsertAuditLog(tx, auditLog)
	if err != nil {
		return errors.Join(errors.New("Cannot insert player merger audit log"), err)
	}

	err = InsertAuditLogPlayerAffected(tx, NewAuditLogPlayerAffected(auditLog.Id, target, 0))
	if err != nil {
		return errors.Join(errors.New("Cannot insert player affected audit log (target)"), err)
	}

	err = InsertAuditLogPlayerAffected(tx, NewAuditLogPlayerAffected(auditLog.Id, dest, 0))
	if err != nil {
		return errors.Join(errors.New("Cannot insert player affected audit log (dest)"), err)
	}

	var ikey, firstTargetGameIkey int64
	err = tx.Get(&firstTargetGameIkey, "SELECT ikey FROM games WHERE (player_white=$1 OR player_black=$1) ORDER BY played ASC LIMIT 1;", target)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Join(errors.New("Cannot get first game from of the target player"), err)
	}
	noGames := errors.Is(err, sql.ErrNoRows)

	// The games should be replayed from the first dest game before target joined such that liglicko2 parameters are correct
	err = tx.Get(&ikey, "SELECT ikey FROM games WHERE (player_white=$1 OR player_black=$1) AND ikey < $2 ORDER BY played DESC LIMIT 1;", dest, firstTargetGameIkey)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return errors.Join(errors.New("Cannot get first game from of the target player"), err)
	}

	if errors.Is(err, sql.ErrNoRows) {
		ikey = firstTargetGameIkey
	}

	// Update the player to be deleted and tagged as merged
	_, err = tx.Exec("UPDATE players SET deleted=TRUE, name=name || '-MERGED', name_normalised=id WHERE id=$1", target)
	if err != nil {
		return errors.Join(errors.New("Cannot set player to deleted with merged status"), err)
	}

	if noGames {
		return tx.Commit()
	}

	// Set target player to dest in all games
	var affectedGames []int64
	err = tx.Select(&affectedGames, "SELECT ikey FROM games WHERE  (player_white=$1 OR player_black=$1)", target)
	if err != nil {
		return errors.Join(errors.New("Cannot get a list of affected games for player merger"), err)
	}

	for _, gameIkey := range affectedGames {
		err = InsertAuditLogGameAffected(tx, &AuditLogGameAffected{
			AuditLogId: auditLog.Id,
			GameIkey:   gameIkey,
		})

		if err != nil {
			return errors.Join(errors.New("Cannot insert game affected audit log"), err)
		}
	}

	_, err = tx.Exec("UPDATE games SET player_white=$1 WHERE player_white=$2", dest, target)
	if err != nil {
		return errors.Join(errors.New("Cannot update white players"), err)
	}

	_, err = tx.Exec("UPDATE games SET player_black=$1 WHERE player_black=$2", dest, target)
	if err != nil {
		return errors.Join(errors.New("Cannot update black players"), err)
	}

	// DELETE target vs dest games
	_, err = tx.Exec("UPDATE games SET deleted=TRUE WHERE (player_white=$1 AND player_black=$1)", dest)
	if err != nil {
		return errors.Join(errors.New("Cannot delete games where the player played against themselves"), err)
	}

	games, players, err := ReplayFrom(tx, ikey)
	if err != nil {
		return errors.Join(errors.New("Cannot replay games to calculate new ELOs"), err)
	}

	for _, game := range games {
		if slices.Contains(affectedGames, game.IKey) {
			continue
		}

		err = InsertAuditLogGameAffected(tx, &AuditLogGameAffected{
			AuditLogId: auditLog.Id,
			GameIkey:   game.IKey,
		})
		if err != nil {
			return errors.Join(errors.New("Cannot insert game affected audit log"), err)
		}
	}

	for _, player := range players {
		if player.Id == target || player.Id == dest {
			continue
		}

		err = InsertAuditLogPlayerAffected(tx, NewAuditLogPlayerAffected(auditLog.Id, player.Id, 0))
		if err != nil {
			return errors.Join(errors.New("Cannot insert player affected audit log"), err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.Join(errors.New("Cannot commit transaction"), err)
	}

	return nil
}
