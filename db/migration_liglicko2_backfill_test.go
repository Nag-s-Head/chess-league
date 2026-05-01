package db_test

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestMigrateFromPreLiglicko2SchemaBackfillsLiglicko2(t *testing.T) {
	testDb := testutils.GetDb(t)
	defer testDb.Close()

	baseConn := testDb.GetSqlxDb()

	schemaName := "migration_liglicko2_" + strings.ReplaceAll(uuid.New().String(), "-", "")
	_, err := baseConn.Exec(`CREATE SCHEMA ` + schemaName)
	require.NoError(t, err)
	defer func() {
		_, _ = baseConn.Exec(`DROP SCHEMA IF EXISTS ` + schemaName + ` CASCADE`)
	}()

	scopedDsn, err := databaseURLWithSearchPath(schemaName)
	require.NoError(t, err)
	scopedConn, err := sqlx.Connect("postgres", scopedDsn)
	require.NoError(t, err)
	defer scopedConn.Close()

	_, err = scopedConn.Exec(`
CREATE TABLE players (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	name_normalised TEXT NOT NULL UNIQUE,
	elo INTEGER DEFAULT 1000 CHECK(elo >= 0),
	join_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted BOOL NOT NULL DEFAULT false
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
	submit_user_agent TEXT,
	ikey BIGINT NOT NULL UNIQUE
);

CREATE TABLE migrations (
  version INTEGER PRIMARY KEY,
  date TIMESTAMPTZ NOT NULL
);
`)
	require.NoError(t, err)

	playerOneID := uuid.New()
	playerTwoID := uuid.New()
	playerOneJoin := time.Date(2025, 1, 5, 11, 0, 0, 0, time.UTC)
	playerTwoJoin := time.Date(2025, 1, 10, 9, 30, 0, 0, time.UTC)
	gameOnePlayed := time.Date(2025, 2, 1, 10, 0, 0, 0, time.UTC)
	gameTwoPlayed := time.Date(2025, 2, 8, 20, 0, 0, 0, time.UTC)

	_, err = scopedConn.Exec(`
INSERT INTO players (id, name, name_normalised, elo, join_time, deleted) VALUES
($1, 'Alice', 'alice', 1042, $2, false),
($3, 'Bob', 'bob', 958, $4, false);
`,
		playerOneID.String(), playerOneJoin,
		playerTwoID.String(), playerTwoJoin,
	)
	require.NoError(t, err)

	_, err = scopedConn.Exec(`
INSERT INTO games (player_white, player_black, score, submitter, played, deleted, elo_given, elo_taken, submit_ip, submit_user_agent, ikey) VALUES
($1, $2, '1-0', $1, $3, false, 16, -16, '', '', 1),
($2, $1, '1/2-1/2', $2, $4, false, 0, 0, '', '', 2);
`,
		playerOneID.String(),
		playerTwoID.String(),
		gameOnePlayed,
		gameTwoPlayed,
	)
	require.NoError(t, err)

	_, err = scopedConn.Exec(`INSERT INTO migrations(version, date) VALUES (9, NOW());`)
	require.NoError(t, err)

	migratedDb, err := db.From(scopedConn)
	require.NoError(t, err)
	defer migratedDb.Close()

	require.True(t, columnExists(t, scopedConn, schemaName, "players", "liglicko2_rating"))
	require.True(t, columnExists(t, scopedConn, schemaName, "players", "liglicko2_deviation"))
	require.True(t, columnExists(t, scopedConn, schemaName, "players", "liglicko2_volatility"))
	require.True(t, columnExists(t, scopedConn, schemaName, "players", "liglicko2_at"))
	require.True(t, columnExists(t, scopedConn, schemaName, "games", "liglicko2_white"))
	require.True(t, columnExists(t, scopedConn, schemaName, "games", "liglicko2_black"))

	require.True(t, columnExists(t, scopedConn, schemaName, "games", "liglicko2_white_old_rating"))
	require.True(t, columnExists(t, scopedConn, schemaName, "games", "liglicko2_white_old_deviation"))
	require.True(t, columnExists(t, scopedConn, schemaName, "games", "liglicko2_white_old_volatility"))
	require.True(t, columnExists(t, scopedConn, schemaName, "games", "liglicko2_white_old_at"))

	require.True(t, columnExists(t, scopedConn, schemaName, "games", "liglicko2_black_old_rating"))
	require.True(t, columnExists(t, scopedConn, schemaName, "games", "liglicko2_black_old_deviation"))
	require.True(t, columnExists(t, scopedConn, schemaName, "games", "liglicko2_black_old_volatility"))
	require.True(t, columnExists(t, scopedConn, schemaName, "games", "liglicko2_black_old_at"))

	var migrationVersion int
	err = scopedConn.Get(&migrationVersion, `SELECT version FROM migrations ORDER BY version DESC LIMIT 1;`)
	require.NoError(t, err)
	require.GreaterOrEqual(t, migrationVersion, 11)

	playerOneExpected := model.Player{
		Liglicko2Rating:     model.StartingLiglicko2Rating,
		Liglicko2Deviation:  model.StartingLiglicko2Deviation,
		Liglicko2Volatility: model.StartingLiglicko2Volatility,
		Liglicko2At:         liglicko2InstantFromTimeForTest(playerOneJoin),
	}
	playerTwoExpected := model.Player{
		Liglicko2Rating:     model.StartingLiglicko2Rating,
		Liglicko2Deviation:  model.StartingLiglicko2Deviation,
		Liglicko2Volatility: model.StartingLiglicko2Volatility,
		Liglicko2At:         liglicko2InstantFromTimeForTest(playerTwoJoin),
	}

	// Capture states before games
	p1OldGame1 := playerOneExpected
	p2OldGame1 := playerTwoExpected

	gameOneWhiteDeltaExpected, gameOneBlackDeltaExpected, err := model.CalculateLiglicko2(&playerOneExpected, &playerTwoExpected, model.Outcome_Win, gameOnePlayed)
	require.NoError(t, err)

	p1OldGame2 := playerOneExpected
	p2OldGame2 := playerTwoExpected
	gameTwoWhiteDeltaExpected, gameTwoBlackDeltaExpected, err := model.CalculateLiglicko2(&playerTwoExpected, &playerOneExpected, model.Outcome_Draw, gameTwoPlayed)
	require.NoError(t, err)

	type gameLiglicko2 struct {
		IKey           int64   `db:"ikey"`
		Liglicko2White float64 `db:"liglicko2_white"`
		Liglicko2Black float64 `db:"liglicko2_black"`

		Liglicko2WhiteOldRating     float64 `db:"liglicko2_white_old_rating"`
		Liglicko2WhiteOldDeviation  float64 `db:"liglicko2_white_old_deviation"`
		Liglicko2WhiteOldVolatility float64 `db:"liglicko2_white_old_volatility"`
		Liglicko2WhiteOldAt         float64 `db:"liglicko2_white_old_at"`

		Liglicko2BlackOldRating     float64 `db:"liglicko2_black_old_rating"`
		Liglicko2BlackOldDeviation  float64 `db:"liglicko2_black_old_deviation"`
		Liglicko2BlackOldVolatility float64 `db:"liglicko2_black_old_volatility"`
		Liglicko2BlackOldAt         float64 `db:"liglicko2_black_old_at"`
	}
	var migratedGames []gameLiglicko2
	err = scopedConn.Select(&migratedGames, `
SELECT 
    ikey, 
    liglicko2_white, 
    liglicko2_black, 
    liglicko2_white_old_rating, 
    liglicko2_white_old_deviation, 
    liglicko2_white_old_volatility, 
    liglicko2_white_old_at, 
    liglicko2_black_old_rating, 
    liglicko2_black_old_deviation, 
    liglicko2_black_old_volatility, 
    liglicko2_black_old_at 
FROM games ORDER BY ikey ASC;`)
	require.NoError(t, err)
	require.Len(t, migratedGames, 2)

	// Game 1
	require.InDelta(t, gameOneWhiteDeltaExpected, migratedGames[0].Liglicko2White, 1e-9)
	require.InDelta(t, gameOneBlackDeltaExpected, migratedGames[0].Liglicko2Black, 1e-9)
	require.InDelta(t, p1OldGame1.Liglicko2Rating, migratedGames[0].Liglicko2WhiteOldRating, 1e-9)
	require.InDelta(t, p1OldGame1.Liglicko2Deviation, migratedGames[0].Liglicko2WhiteOldDeviation, 1e-9)
	require.InDelta(t, p1OldGame1.Liglicko2Volatility, migratedGames[0].Liglicko2WhiteOldVolatility, 1e-9)
	require.InDelta(t, p1OldGame1.Liglicko2At, migratedGames[0].Liglicko2WhiteOldAt, 1e-9)
	require.InDelta(t, p2OldGame1.Liglicko2Rating, migratedGames[0].Liglicko2BlackOldRating, 1e-9)
	require.InDelta(t, p2OldGame1.Liglicko2Deviation, migratedGames[0].Liglicko2BlackOldDeviation, 1e-9)
	require.InDelta(t, p2OldGame1.Liglicko2Volatility, migratedGames[0].Liglicko2BlackOldVolatility, 1e-9)
	require.InDelta(t, p2OldGame1.Liglicko2At, migratedGames[0].Liglicko2BlackOldAt, 1e-9)

	// Game 2
	require.InDelta(t, gameTwoWhiteDeltaExpected, migratedGames[1].Liglicko2White, 1e-9)
	require.InDelta(t, gameTwoBlackDeltaExpected, migratedGames[1].Liglicko2Black, 1e-9)
	require.InDelta(t, p2OldGame2.Liglicko2Rating, migratedGames[1].Liglicko2WhiteOldRating, 1e-9)
	require.InDelta(t, p2OldGame2.Liglicko2Deviation, migratedGames[1].Liglicko2WhiteOldDeviation, 1e-9)
	require.InDelta(t, p2OldGame2.Liglicko2Volatility, migratedGames[1].Liglicko2WhiteOldVolatility, 1e-9)
	require.InDelta(t, p2OldGame2.Liglicko2At, migratedGames[1].Liglicko2WhiteOldAt, 1e-9)
	require.InDelta(t, p1OldGame2.Liglicko2Rating, migratedGames[1].Liglicko2BlackOldRating, 1e-9)
	require.InDelta(t, p1OldGame2.Liglicko2Deviation, migratedGames[1].Liglicko2BlackOldDeviation, 1e-9)
	require.InDelta(t, p1OldGame2.Liglicko2Volatility, migratedGames[1].Liglicko2BlackOldVolatility, 1e-9)
	require.InDelta(t, p1OldGame2.Liglicko2At, migratedGames[1].Liglicko2BlackOldAt, 1e-9)

	type playerLiglicko2 struct {
		ID                 string  `db:"id"`
		Elo                int     `db:"elo"`
		Liglicko2Rating    float64 `db:"liglicko2_rating"`
		Liglicko2Deviation float64 `db:"liglicko2_deviation"`
		Liglicko2Volatile  float64 `db:"liglicko2_volatility"`
		Liglicko2At        float64 `db:"liglicko2_at"`
	}
	var migratedPlayers []playerLiglicko2
	err = scopedConn.Select(&migratedPlayers, `
SELECT id, elo, liglicko2_rating, liglicko2_deviation, liglicko2_volatility, liglicko2_at
FROM players
ORDER BY name_normalised ASC;`)
	require.NoError(t, err)
	require.Len(t, migratedPlayers, 2)

	// Elo should remain unchanged by liglicko2 migration.
	require.Equal(t, 1042, migratedPlayers[0].Elo)
	require.Equal(t, 958, migratedPlayers[1].Elo)

	require.InDelta(t, playerOneExpected.Liglicko2Rating, migratedPlayers[0].Liglicko2Rating, 1e-9)
	require.InDelta(t, playerOneExpected.Liglicko2Deviation, migratedPlayers[0].Liglicko2Deviation, 1e-9)
	require.InDelta(t, playerOneExpected.Liglicko2Volatility, migratedPlayers[0].Liglicko2Volatile, 1e-9)
	require.InDelta(t, playerOneExpected.Liglicko2At, migratedPlayers[0].Liglicko2At, 1e-9)

	require.InDelta(t, playerTwoExpected.Liglicko2Rating, migratedPlayers[1].Liglicko2Rating, 1e-9)
	require.InDelta(t, playerTwoExpected.Liglicko2Deviation, migratedPlayers[1].Liglicko2Deviation, 1e-9)
	require.InDelta(t, playerTwoExpected.Liglicko2Volatility, migratedPlayers[1].Liglicko2Volatile, 1e-9)
	require.InDelta(t, playerTwoExpected.Liglicko2At, migratedPlayers[1].Liglicko2At, 1e-9)
}

func columnExists(t *testing.T, sqlDb *sqlx.DB, schemaName, tableName, columnName string) bool {
	t.Helper()

	var count int
	err := sqlDb.Get(&count, `
SELECT COUNT(*)
FROM information_schema.columns
WHERE table_schema = $1 AND table_name = $2 AND column_name = $3;`,
		schemaName,
		tableName,
		columnName,
	)
	require.NoError(t, err)
	return count == 1
}

func databaseURLWithSearchPath(searchPath string) (string, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return "", errors.New("DATABASE_URL is not set")
	}

	if strings.HasPrefix(databaseURL, "postgres://") || strings.HasPrefix(databaseURL, "postgresql://") {
		separator := "?"
		if strings.Contains(databaseURL, "?") {
			separator = "&"
		}
		return databaseURL + separator + "search_path=" + searchPath, nil
	}

	return fmt.Sprintf("%s search_path=%s", databaseURL, searchPath), nil
}

func liglicko2InstantFromTimeForTest(timestamp time.Time) float64 {
	return float64(timestamp.UnixNano()) / float64(24*time.Hour)
}
