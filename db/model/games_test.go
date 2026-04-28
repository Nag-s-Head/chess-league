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
	db := testutils.GetDb(t)
	defer db.Close()

	ikey1, err := model.NextIKey(db)
	require.NoError(t, err)

	ikey2, err := model.NextIKey(db)
	require.NoError(t, err)

	require.NotEqual(t, ikey1, ikey2)
}

func TestCreateGameP1White(t *testing.T) {
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
	game, eloWhite, eloBlack, err := model.CreateGame(tx, &p1, &p2, true, ikey, model.Score_Win, r)
	require.NoError(t, err)

	require.NotEqual(t, model.StartingElo, p1.DEPRECATEDElo)
	require.NotEqual(t, model.StartingElo, p2.DEPRECATEDElo)
	require.NotEqual(t, model.StartingLiglicko2Rating, p1.Liglicko2Rating)
	require.NotEqual(t, model.StartingLiglicko2Rating, p2.Liglicko2Rating)

	require.Equal(t, p1.Id, game.Submitter)
	require.Equal(t, p1.Id, game.PlayerWhite)
	require.Equal(t, p2.Id, game.PlayerBlack)

	require.Equal(t, int(p1.Liglicko2Rating-model.StartingLiglicko2Rating), eloWhite)
	require.Equal(t, int(p2.Liglicko2Rating-model.StartingLiglicko2Rating), eloBlack)
	require.InDelta(t, p1.Liglicko2Rating-model.StartingLiglicko2Rating, game.Liglicko2White, 1e-9)
	require.InDelta(t, p2.Liglicko2Rating-model.StartingLiglicko2Rating, game.Liglicko2Black, 1e-9)

	require.Equal(t, model.Score_Win, game.Score)
	require.Equal(t, false, game.Deleted)
	require.NotEqual(t, 0.0, game.Liglicko2White)
	require.NotEqual(t, 0.0, game.Liglicko2Black)

	require.Equal(t, ikey, game.IKey)
	require.NoError(t, tx.Commit())

	p1New, err := model.GetPlayer(db, p1.Id)
	require.NoError(t, err)
	require.Equal(t, p1.DEPRECATEDElo, p1New.DEPRECATEDElo)
	require.Equal(t, p1.Liglicko2Rating, p1New.Liglicko2Rating)

	p2New, err := model.GetPlayer(db, p2.Id)
	require.NoError(t, err)
	require.Equal(t, p2.DEPRECATEDElo, p2New.DEPRECATEDElo)
	require.Equal(t, p2.Liglicko2Rating, p2New.Liglicko2Rating)
}

func TestCreateGameP1Black(t *testing.T) {
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
	game, eloWhite, eloBlack, err := model.CreateGame(tx, &p1, &p2, false, ikey, model.Score_Win, r)
	require.NoError(t, err)

	require.NotEqual(t, model.StartingElo, p1.DEPRECATEDElo)
	require.NotEqual(t, model.StartingElo, p2.DEPRECATEDElo)
	require.NotEqual(t, model.StartingLiglicko2Rating, p1.Liglicko2Rating)
	require.NotEqual(t, model.StartingLiglicko2Rating, p2.Liglicko2Rating)

	require.Equal(t, p1.Id, game.Submitter)
	require.Equal(t, p2.Id, game.PlayerWhite)
	require.Equal(t, p1.Id, game.PlayerBlack)

	require.Equal(t, int(p2.Liglicko2Rating-model.StartingLiglicko2Rating), eloWhite)
	require.Equal(t, int(p1.Liglicko2Rating-model.StartingLiglicko2Rating), eloBlack)
	require.InDelta(t, p2.Liglicko2Rating-model.StartingLiglicko2Rating, game.Liglicko2White, 1e-9)
	require.InDelta(t, p1.Liglicko2Rating-model.StartingLiglicko2Rating, game.Liglicko2Black, 1e-9)

	require.Equal(t, model.Score_Win, game.Score)
	require.Equal(t, false, game.Deleted)
	require.NotEqual(t, 0.0, game.Liglicko2White)
	require.NotEqual(t, 0.0, game.Liglicko2Black)

	require.Equal(t, ikey, game.IKey)
	require.NoError(t, tx.Commit())

	p1New, err := model.GetPlayer(db, p1.Id)
	require.NoError(t, err)
	require.Equal(t, p1.DEPRECATEDElo, p1New.DEPRECATEDElo)
	require.Equal(t, p1.Liglicko2Rating, p1New.Liglicko2Rating)

	p2New, err := model.GetPlayer(db, p2.Id)
	require.NoError(t, err)
	require.Equal(t, p2.DEPRECATEDElo, p2New.DEPRECATEDElo)
	require.Equal(t, p2.Liglicko2Rating, p2New.Liglicko2Rating)
}

func TestCreateGameSameIkeyFails(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	tx, err := db.GetSqlxDb().BeginTxx(context.Background(), nil)
	require.NoError(t, err)
	defer tx.Commit()

	p1 := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayer(db, p1))

	p2 := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayer(db, p2))

	ikey, err := model.NextIKey(db)
	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodGet, "/mocked-url", strings.NewReader(""))
	_, _, _, err = model.CreateGame(tx, &p1, &p2, true, ikey, model.Score_Win, r)
	require.NoError(t, err)

	_, _, _, err = model.CreateGame(tx, &p1, &p2, true, ikey, model.Score_Win, r)
	require.Error(t, err)
}

func TestMapGamesToUiFriendly(t *testing.T) {
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
