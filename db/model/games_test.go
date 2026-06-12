package model_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Nag-s-Head/chess-league/db/model"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNextIkey(t *testing.T) {
	t.Parallel()

	db := testutils.GetDb(t)
	defer db.Close()

	ikey1, err := model.NextIKey(db)
	require.NoError(t, err)

	ikey2, err := model.NextIKey(db)
	require.NoError(t, err)

	require.NotEqual(t, ikey1, ikey2)
}

func TestMapGamesToUiFriendly(t *testing.T) {
	t.Parallel()

	db := testutils.GetDb(t)
	defer db.Close()

	tx, err := db.GetSqlxDb().BeginTxx(context.Background(), nil)
	require.NoError(t, err)
	defer tx.Commit()

	p1 := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayer(db, p1))

	games, err := model.GetGamesByPlayer(db, p1.Id)
	require.NoError(t, err)

	details := model.MapGamesToUserFriendly(p1.Id, games)
	require.NotEmpty(t, details)
}

func TestMapGamesToUiFriendlyDrawUsesLiglicko2PerColor(t *testing.T) {
	t.Parallel()

	player := model.NewPlayer(uuid.New().String())
	opponent := model.NewPlayer(uuid.New().String())

	games := []model.GameWithPlayerNames{
		{
			Game: model.Game{
				PlayerWhite:    player.Id,
				PlayerBlack:    opponent.Id,
				Score:          model.Score_Draw,
				Played:         time.Now(),
				Liglicko2White: 4.2,
				Liglicko2Black: -4.2,
			},
			WhiteName: player.Name,
			BlackName: opponent.Name,
		},
	}

	details := model.MapGamesToUserFriendly(player.Id, games)
	require.Len(t, details.Games, 1)
	require.Equal(t, "Draw", details.Games[0].Outcome)
	require.InDelta(t, 4.2, details.Games[0].Liglicko2Change, 1e-9)
}

func TestGetGame(t *testing.T) {
	t.Parallel()

	db := testutils.GetDb(t)
	defer db.Close()

	tx, err := db.GetSqlxDb().BeginTxx(context.Background(), nil)
	require.NoError(t, err)
	defer tx.Rollback()

	p1 := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayer(db, p1))

	p2 := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayer(db, p2))

	ikey, err := model.NextIKey(db)
	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodGet, "/mocked-url", strings.NewReader(""))
	game, _, _, err := model.CreateGame(tx, &p1, &p2, false, ikey, model.Score_Win, r)
	require.NoError(t, err)

	require.Equal(t, ikey, game.IKey)
	require.NoError(t, tx.Commit())

	game2, err := model.GetGameWithDetails(db, game.IKey)
	require.NoError(t, err)
	require.Equal(t, game.IKey, game2.IKey)
}

func TestGetGamesByPlayerPairCombsNoGamesPlayed(t *testing.T) {
	t.Parallel()

	db := testutils.GetDb(t)
	defer db.Close()

	games, err := model.GetGamesByPlayerPairCombs(db, []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}, []uuid.UUID{uuid.New(), uuid.New(), uuid.New()})
	require.NoError(t, err)
	require.Len(t, games, 0)
}

func TestGetGamesByPlayerPairCombsSingleGamePlayed(t *testing.T) {
	t.Parallel()

	db := testutils.GetDb(t)
	defer db.Close()

	tx, err := db.GetSqlxDb().BeginTxx(context.Background(), nil)
	require.NoError(t, err)
	defer tx.Rollback()

	p1 := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayer(db, p1))

	p2 := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayer(db, p2))

	r := httptest.NewRequest(http.MethodGet, "/mocked-url", strings.NewReader(""))

	// Game 1: P1 vs P2, P1 wins
	ikey, _ := model.NextIKey(db)
	game, _, _, err := model.CreateGame(tx, &p1, &p2, true, ikey, model.Score_Win, r)
	require.NoError(t, err)
	require.NoError(t, tx.Commit())

	games, err := model.GetGamesByPlayerPairCombs(db, []uuid.UUID{p1.Id, uuid.New(), uuid.New()}, []uuid.UUID{p2.Id, uuid.New(), uuid.New()})
	require.NoError(t, err)
	require.Len(t, games, 1)
	require.Equal(t, games[0].IKey, game.IKey)
}

func TestGetGamesByPlayerPairCombsManyGames(t *testing.T) {
	t.Parallel()

	db := testutils.GetDb(t)
	defer db.Close()

	tx, err := db.GetSqlxDb().BeginTxx(context.Background(), nil)
	require.NoError(t, err)
	defer tx.Rollback()

	pAIds := make([]uuid.UUID, 0)
	pBIds := make([]uuid.UUID, 0)

	for range 5 {
		p1 := model.NewPlayer(uuid.New().String())
		require.NoError(t, model.InsertPlayer(db, p1))

		pAIds = append(pAIds, p1.Id)

		p2 := model.NewPlayer(uuid.New().String())
		require.NoError(t, model.InsertPlayer(db, p2))

		pBIds = append(pBIds, p2.Id)

		r := httptest.NewRequest(http.MethodGet, "/mocked-url", strings.NewReader(""))

		for range 5 {
			ikey, _ := model.NextIKey(db)
			_, _, _, err := model.CreateGame(tx, &p1, &p2, true, ikey, model.Score_Win, r)
			require.NoError(t, err)
		}
	}

	require.NoError(t, tx.Commit())

	games, err := model.GetGamesByPlayerPairCombs(db, pAIds, pBIds)
	require.NoError(t, err)
	require.True(t, len(games) > 5)
}
