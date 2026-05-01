package model_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Nag-s-Head/chess-league/db/model"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestReplayFromConsistency(t *testing.T) {
	t.Parallel()

	db := testutils.GetDb(t)
	defer db.Close()

	// 1. Setup Players
	p1 := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayer(db, p1))

	p2 := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayer(db, p2))

	r := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	r.RemoteAddr = "127.0.0.1"

	// 2. Create a sequence of games
	tx, err := db.GetSqlxDb().BeginTxx(context.Background(), nil)
	require.NoError(t, err)

	createGame := func(tx *sqlx.Tx, white, black *model.Player, score model.Score) model.Game {
		ikey, _ := model.NextIKey(db)
		game, _, _, err := model.CreateGame(tx, white, black, true, ikey, score, r)
		require.NoError(t, err)
		return game
	}

	g1 := createGame(tx, &p1, &p2, model.Score_Win)
	createGame(tx, &p2, &p1, model.Score_Win) // P2 wins
	createGame(tx, &p1, &p2, model.Score_Draw)

	require.NoError(t, tx.Commit())

	// 3. Capture final state
	p1Orig, _ := model.GetPlayer(db, p1.Id)
	p2Orig, _ := model.GetPlayer(db, p2.Id)

	var gamesOrig []model.Game
	err = db.GetSqlxDb().Select(&gamesOrig, "SELECT * FROM games WHERE player_white IN ($1, $2) ORDER BY played ASC", p1.Id, p2.Id)
	require.NoError(t, err)
	require.Len(t, gamesOrig, 3)

	// 4. Replay and verify consistency
	tx, err = db.GetSqlxDb().BeginTxx(context.Background(), nil)
	require.NoError(t, err)
	defer tx.Rollback()

	affectedGames, affectedPlayers, err := model.ReplayFrom(tx, g1.Played)
	require.NoError(t, err)

	// Find our players in the affected list
	var p1Replayed, p2Replayed *model.Player
	for _, p := range affectedPlayers {
		if p.Id == p1.Id {
			p1Replayed = p
		} else if p.Id == p2.Id {
			p2Replayed = p
		}
	}

	require.NotNil(t, p1Replayed)
	require.NotNil(t, p2Replayed)

	// Verify final ratings are identical to original ones
	require.InDelta(t, p1Orig.Liglicko2Rating, p1Replayed.Liglicko2Rating, 1e-9)
	require.InDelta(t, p1Orig.Liglicko2Deviation, p1Replayed.Liglicko2Deviation, 1e-9)
	require.InDelta(t, p1Orig.Liglicko2Volatility, p1Replayed.Liglicko2Volatility, 1e-9)

	require.InDelta(t, p2Orig.Liglicko2Rating, p2Replayed.Liglicko2Rating, 1e-9)
	require.InDelta(t, p2Orig.Liglicko2Deviation, p2Replayed.Liglicko2Deviation, 1e-9)
	require.InDelta(t, p2Orig.Liglicko2Volatility, p2Replayed.Liglicko2Volatility, 1e-9)

	// Verify all game deltas and old states are also identical
	for _, gReplayed := range affectedGames {
		for _, gOrig := range gamesOrig {
			if gReplayed.IKey == gOrig.IKey {
				require.InDelta(t, gOrig.Liglicko2White, gReplayed.Liglicko2White, 1e-9)
				require.InDelta(t, gOrig.Liglicko2Black, gReplayed.Liglicko2Black, 1e-9)
				require.InDelta(t, gOrig.Liglicko2WhiteOldRating, gReplayed.Liglicko2WhiteOldRating, 1e-9)
				require.InDelta(t, gOrig.Liglicko2BlackOldRating, gReplayed.Liglicko2BlackOldRating, 1e-9)
			}
		}
	}
}

func TestReplayAfterEdit(t *testing.T) {
	t.Parallel()

	db := testutils.GetDb(t)
	defer db.Close()

	// 1. Setup Players
	p1 := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayer(db, p1))

	p2 := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayer(db, p2))

	r := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	r.RemoteAddr = "127.0.0.1"

	// 2. Create two games
	tx, err := db.GetSqlxDb().BeginTxx(context.Background(), nil)
	require.NoError(t, err)

	ikey1, _ := model.NextIKey(db)
	g1, _, _, err := model.CreateGame(tx, &p1, &p2, true, ikey1, model.Score_Win, r)
	require.NoError(t, err)

	ikey2, _ := model.NextIKey(db)
	_, _, _, err = model.CreateGame(tx, &p1, &p2, true, ikey2, model.Score_Loss, r)
	require.NoError(t, err)

	require.NoError(t, tx.Commit())

	// Capture state before edit
	p1AfterWin, _ := model.GetPlayer(db, p1.Id)

	// 3. Change Game 1 to Draw and replay
	tx, err = db.GetSqlxDb().BeginTxx(context.Background(), nil)
	require.NoError(t, err)
	defer tx.Rollback()

	_, err = tx.Exec("UPDATE games SET score = $1 WHERE ikey = $2", model.Score_Draw, ikey1)
	require.NoError(t, err)

	affectedGames, affectedPlayers, err := model.ReplayFrom(tx, g1.Played)
	require.NoError(t, err)

	var g1Final, g2Final *model.Game
	for i, g := range affectedGames {
		if g.IKey == ikey1 {
			g1Final = &affectedGames[i]
		} else if g.IKey == ikey2 {
			g2Final = &affectedGames[i]
		}
	}

	require.NotNil(t, g1Final)
	require.NotNil(t, g2Final)

	// Game 1 was Draw, so delta should be 0 (both were 1500)
	require.InDelta(t, 0.0, g1Final.Liglicko2White, 1e-9)

	// Game 2 OldRating should now be 1500 (since Game 1 was Draw)
	// instead of the > 1500 rating from the original Win
	require.InDelta(t, 1500.0, g2Final.Liglicko2WhiteOldRating, 1e-9)

	// Final P1 rating should be different from before
	var p1Final *model.Player
	for _, p := range affectedPlayers {
		if p.Id == p1.Id {
			p1Final = p
		}
	}
	require.NotEqual(t, p1AfterWin.Liglicko2Rating, p1Final.Liglicko2Rating)
}
