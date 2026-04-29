package db

import (
	"errors"
	"fmt"
	"time"

	liglicko2core "github.com/Nag-s-Head/chess-league/db/liglicko2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type legacyLiglicko2Player struct {
	Id       uuid.UUID `db:"id"`
	JoinTime time.Time `db:"join_time"`
}

type legacyLiglicko2Game struct {
	IKey        int64     `db:"ikey"`
	PlayerWhite uuid.UUID `db:"player_white"`
	PlayerBlack uuid.UUID `db:"player_black"`
	Score       string    `db:"score"`
	Played      time.Time `db:"played"`
}

// This is the initial migration from old ELO to liglicko2
func InternalMigrateLegacyLiglicko2(tx *sqlx.Tx) error {
	statements := []string{
		`ALTER TABLE players ADD COLUMN IF NOT EXISTS liglicko2_rating DOUBLE PRECISION NOT NULL DEFAULT 1500;`,
		`ALTER TABLE players ADD COLUMN IF NOT EXISTS liglicko2_deviation DOUBLE PRECISION NOT NULL DEFAULT 500 CHECK (liglicko2_deviation >= 0);`,
		`ALTER TABLE players ADD COLUMN IF NOT EXISTS liglicko2_volatility DOUBLE PRECISION NOT NULL DEFAULT 0.09 CHECK (liglicko2_volatility >= 0);`,
		`ALTER TABLE players ADD COLUMN IF NOT EXISTS liglicko2_at DOUBLE PRECISION NOT NULL DEFAULT 0;`,
		`ALTER TABLE games ADD COLUMN IF NOT EXISTS liglicko2_white DOUBLE PRECISION NOT NULL DEFAULT 0;`,
		`ALTER TABLE games ADD COLUMN IF NOT EXISTS liglicko2_black DOUBLE PRECISION NOT NULL DEFAULT 0;`,
		`CREATE INDEX IF NOT EXISTS idx_players_liglicko2_rating ON players(liglicko2_rating);`,
	}

	for _, sql := range statements {
		if _, err := tx.Exec(sql); err != nil {
			return errors.Join(errors.New("cannot apply liglicko2 schema update"), err)
		}
	}

	var players []legacyLiglicko2Player
	if err := tx.Select(&players, "SELECT id, join_time FROM players;"); err != nil {
		return errors.Join(errors.New("cannot load players for liglicko2 backfill"), err)
	}

	ratings := make(map[uuid.UUID]liglicko2core.Rating, len(players))
	for _, p := range players {
		ratings[p.Id] = liglicko2core.Rating{
			Rating:     liglicko2core.DefaultRating,
			Deviation:  liglicko2core.DefaultDeviation,
			Volatility: liglicko2core.DefaultVolatility,
			At:         liglicko2core.InstantFromTime(p.JoinTime),
		}
	}

	var games []legacyLiglicko2Game
	if err := tx.Select(&games, `
SELECT ikey, player_white, player_black, score, played
FROM games
ORDER BY played ASC, ikey ASC;`); err != nil {
		return errors.Join(errors.New("cannot load games for liglicko2 backfill"), err)
	}

	for _, g := range games {
		white, okW := ratings[g.PlayerWhite]
		black, okB := ratings[g.PlayerBlack]
		if !okW || !okB {
			return fmt.Errorf("cannot backfill game %d due to missing player rows", g.IKey)
		}

		score, err := liglicko2ScoreFromResult(g.Score)
		if err != nil {
			return errors.Join(fmt.Errorf("cannot parse score for game %d", g.IKey), err)
		}

		now := liglicko2core.InstantFromTime(g.Played)
		nextWhite, err := liglicko2core.UpdateSingle(white, black, score, now, liglicko2core.FirstAdvantage)
		if err != nil {
			return errors.Join(fmt.Errorf("cannot backfill white liglicko2 for game %d", g.IKey), err)
		}
		nextBlack, err := liglicko2core.UpdateSingle(black, white, 1.0-score, now, -liglicko2core.FirstAdvantage)
		if err != nil {
			return errors.Join(fmt.Errorf("cannot backfill black liglicko2 for game %d", g.IKey), err)
		}

		deltaWhite := nextWhite.Rating - white.Rating
		deltaBlack := nextBlack.Rating - black.Rating

		if _, err := tx.Exec(
			"UPDATE games SET liglicko2_white=$1, liglicko2_black=$2 WHERE ikey=$3;",
			deltaWhite,
			deltaBlack,
			g.IKey,
		); err != nil {
			return errors.Join(fmt.Errorf("cannot update liglicko2 deltas for game %d", g.IKey), err)
		}

		ratings[g.PlayerWhite] = nextWhite
		ratings[g.PlayerBlack] = nextBlack
	}

	for playerID, r := range ratings {
		if _, err := tx.Exec(`
UPDATE players
SET liglicko2_rating=$1,
	liglicko2_deviation=$2,
	liglicko2_volatility=$3,
	liglicko2_at=$4
WHERE id=$5;`,
			r.Rating,
			r.Deviation,
			r.Volatility,
			r.At,
			playerID,
		); err != nil {
			return errors.Join(fmt.Errorf("cannot update player liglicko2 state for %s", playerID), err)
		}
	}

	return nil
}

func liglicko2ScoreFromResult(score string) (float64, error) {
	switch score {
	case "1-0":
		return 1.0, nil
	case "0-1":
		return 0.0, nil
	case "1/2-1/2":
		return 0.5, nil
	default:
		return 0.0, fmt.Errorf("unknown score: %s", score)
	}
}

// This migration stores the deviation, volatility, and last rating calculation time ("at") for each of the old games so that replay after game edit works
func InternalMigrateLiglicko2ToAddOldStatesToGames(tx *sqlx.Tx) error {
	// Game type is copied from model when the migration was written to keep it working forver
	type Game struct {
		PlayerWhite uuid.UUID `db:"player_white"`
		PlayerBlack uuid.UUID `db:"player_black"`
		Score       string    `db:"score"`
		Submitter   uuid.UUID `db:"submitter"`
		Played      time.Time `db:"played"`
		IKey        int64     `db:"ikey"`

		// Liglicko2White and Liglicko2Black are per-game liglicko2 deltas for each side.
		// They preserve sign, so draws between uneven players can still show non-zero
		// changes.
		Liglicko2White              float64 `db:"liglicko2_white"`
		Liglicko2WhiteOldRating     float64 `db:"liglicko2_white_old_rating"`
		Liglicko2WhiteOldVolatility float64 `db:"liglicko2_white_old_volatility"`
		Liglicko2WhiteOldDeviation  float64 `db:"liglicko2_white_old_deviation"`
		Liglicko2WhiteOldAt         float64 `db:"liglicko2_white_old_at"`

		Liglicko2Black              float64 `db:"liglicko2_black"`
		Liglicko2BlackOldRating     float64 `db:"liglicko2_black_old_rating"`
		Liglicko2BlackOldVolatility float64 `db:"liglicko2_black_old_volatility"`
		Liglicko2BlackOldDeviation  float64 `db:"liglicko2_black_old_deviation"`
		Liglicko2BlackOldAt         float64 `db:"liglicko2_black_old_at"`
	}

	// Player type is copied from model when the migration was written to keep it working forever
	type Player struct {
		Id uuid.UUID `db:"id"`
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
	}

	_, err := tx.Exec(`
-- The liglicko2 state of each player at the time before the game is saved so that replays (due to deletion or result swap) do not require replaying the full database
-- It is more cumbersome, but is more correct
ALTER TABLE games ADD COLUMN IF NOT EXISTS liglicko2_white_old_rating DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE games ADD COLUMN IF NOT EXISTS liglicko2_white_old_volatility DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE games ADD COLUMN IF NOT EXISTS liglicko2_white_old_deviation DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE games ADD COLUMN IF NOT EXISTS liglicko2_white_old_at DOUBLE PRECISION NOT NULL DEFAULT 0;

ALTER TABLE games ADD COLUMN IF NOT EXISTS liglicko2_black_old_rating DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE games ADD COLUMN IF NOT EXISTS liglicko2_black_old_volatility DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE games ADD COLUMN IF NOT EXISTS liglicko2_black_old_deviation DOUBLE PRECISION NOT NULL DEFAULT 0;
ALTER TABLE games ADD COLUMN IF NOT EXISTS liglicko2_black_old_at DOUBLE PRECISION NOT NULL DEFAULT 0;
		`)
	if err != nil {
		return errors.Join(errors.New("Cannot add new columns to games table"), err)
	}

	var playersList []Player
	err = tx.Select(&playersList, "SELECT id, join_time FROM players;")
	if err != nil {
		return errors.Join(errors.New("Cannot get get players"), err)
	}

	players := make(map[uuid.UUID]Player)
	for _, player := range playersList {
		players[player.Id] = player
	}

	ratings := make(map[uuid.UUID]liglicko2core.Rating, len(players))
	for _, p := range players {
		ratings[p.Id] = liglicko2core.Rating{
			Rating:     liglicko2core.DefaultRating,
			Deviation:  liglicko2core.DefaultDeviation,
			Volatility: liglicko2core.DefaultVolatility,
			At:         liglicko2core.InstantFromTime(p.JoinTime),
		}
	}

	var games []Game
	err = tx.Select(&games, "SELECT ikey, player_white, player_black, score, played FROM games ORDER BY played ASC;")
	if err != nil {
		return errors.Join(errors.New("Cannot get get games"), err)
	}

	for _, game := range games {
		white, found := players[game.PlayerWhite]
		if !found {
			return fmt.Errorf("Cannot find player %s", game.PlayerWhite)
		}

		whiteRating, found := ratings[game.PlayerWhite]
		if !found {
			return fmt.Errorf("Cannot find player %s", game.PlayerWhite)
		}

		black, found := players[game.PlayerBlack]
		if !found {
			return fmt.Errorf("Cannot find player %s", game.PlayerWhite)
		}

		blackRating, found := ratings[game.PlayerBlack]
		if !found {
			return fmt.Errorf("Cannot find player %s", game.PlayerWhite)
		}

		now := liglicko2core.InstantFromTime(game.Played)

		// Save the current player state in the database
		_, err := tx.Exec(`
UPDATE games 
SET
	liglicko2_white_old_rating=$1,
	liglicko2_white_old_volatility=$2,
	liglicko2_white_old_deviation=$3,

	liglicko2_black_old_rating=$4,
	liglicko2_black_old_volatility=$5,
	liglicko2_black_old_deviation=$6,

	liglicko2_white_old_at=$7,
	liglicko2_black_old_at=$7
WHERE
  ikey = $8;`,
			white.Liglicko2Rating,
			white.Liglicko2Volatility,
			white.Liglicko2Deviation,
			black.Liglicko2Rating,
			black.Liglicko2Volatility,
			black.Liglicko2Deviation,
			now,
			game.IKey,
		)

		score, err := liglicko2ScoreFromResult(game.Score)
		if err != nil {
			return errors.Join(fmt.Errorf("cannot parse score for game %d", game.IKey), err)
		}

		_, err = liglicko2core.UpdateSingle(whiteRating, blackRating, score, now, liglicko2core.FirstAdvantage)
		if err != nil {
			return errors.Join(fmt.Errorf("cannot backfill white liglicko2 for game %d", game.IKey), err)
		}

		_, err = liglicko2core.UpdateSingle(blackRating, whiteRating, 1.0-score, now, -liglicko2core.FirstAdvantage)
		if err != nil {
			return errors.Join(fmt.Errorf("cannot backfill black liglicko2 for game %d", game.IKey), err)
		}
	}

	return nil
}
