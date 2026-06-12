package model

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func GetLeaguePlayers(tx *sqlx.Tx) ([]Player, error) {
	var players []Player
	err := tx.Select(&players,
		`SELECT players.* 
		FROM players 
		INNER JOIN league_players ON 
		  league_players.player_id = players.id 
		ORDER BY liglicko2_rating DESC;`)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot get league players"), err)
	}
	return players, nil
}

func buildPlayerListString(players []*Player) string {
	var builder strings.Builder
	for i, player := range players {
		builder.WriteString(player.Name)
		if i == len(players)-1 {
			builder.WriteRune('.')
		} else if i == len(players)-2 {
			builder.WriteString(", and ")
		} else {
			builder.WriteString(", ")
		}
	}

	return builder.String()
}

func SetLeaguePlayers(db db.Db, adminId uuid.UUID, players []uuid.UUID) error {
	err := db.DoTx(func(tx *sqlx.Tx) error {
		existingPlayers, err := GetLeaguePlayers(tx)
		if err != nil {
			return errors.Join(errors.New("Cannot get existing league players"), err)
		}

		playersAdded := make([]*Player, 0)
		playersRemoved := make([]*Player, 0)

		for _, oldPlayer := range existingPlayers {
			found := false
			for _, newPlayer := range players {
				if newPlayer == oldPlayer.Id {
					found = true
					break
				}
			}

			if !found {
				playersRemoved = append(playersRemoved, &oldPlayer)
			}
		}

		for _, newPlayer := range players {
			found := false
			for _, oldPlayer := range existingPlayers {
				if newPlayer == oldPlayer.Id {
					found = true
					break
				}
			}

			if !found {
				player, err := GetPlayerTx(tx, newPlayer)
				if err != nil {
					return errors.Join(errors.New("Cannot get the data for the added player"), err)
				}

				playersAdded = append(playersAdded, &player)
			}
		}

		auditLog := NewAuditLog(adminId, "League players list updated",
			fmt.Sprintf("The list of players in the league has been updated.\nRemoved players: %s;\nAdded players: %s",
				buildPlayerListString(playersRemoved),
				buildPlayerListString(playersAdded)))
		err = InsertAuditLog(tx, auditLog)
		if err != nil {
			return err
		}

		for _, player := range playersAdded {
			err = InsertAuditLogPlayerAffected(tx, NewAuditLogPlayerAffected(auditLog.Id, player.Id))
			if err != nil {
				return errors.Join(errors.New("Cannot create audit logs for added players"), err)
			}

			_, err = tx.Exec("INSERT INTO league_players (player_id) VALUES ($1);", player.Id)
			if err != nil {
				return errors.Join(errors.New("Cannot insert player into league"), err)
			}
		}

		for _, player := range playersRemoved {
			err = InsertAuditLogPlayerAffected(tx, NewAuditLogPlayerAffected(auditLog.Id, player.Id))
			if err != nil {
				return errors.Join(errors.New("Cannot create audit logs for removed players"), err)
			}

			_, err = tx.Exec("DELETE FROM league_players WHERE player_id=$1;", player.Id)
			if err != nil {
				return errors.Join(errors.New("Cannot delete player from leageue"), err)
			}
		}

		return nil
	})

	if err != nil {
		return errors.Join(errors.New("Cannot set league players"), err)
	}

	return nil
}

type LeaguePlayerUiFriendly struct {
	InLeague bool `db:"in_league"`
	Player
}

func GetUiFriendlyLeaguePlayers(db db.Db) ([]LeaguePlayerUiFriendly, error) {
	var players []LeaguePlayerUiFriendly
	err := db.GetSqlxDb().Select(&players, `
		SELECT 
		  players.*, 
		  EXISTS(
		    SELECT 1 FROM league_players WHERE league_players.player_id = players.id
  		) AS in_league 
		FROM players
		ORDER BY players.name ASC;`)

	if err != nil {
		return nil, errors.Join(errors.New("Cannot get players"), err)
	}

	return players, nil
}
