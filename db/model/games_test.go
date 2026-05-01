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
