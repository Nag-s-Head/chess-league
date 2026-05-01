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
	"github.com/stretchr/testify/require"
)

func TestCreateGameP1White(t *testing.T) {
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

	// Verify old states are set correctly
	require.Equal(t, model.StartingLiglicko2Rating, game.Liglicko2WhiteOldRating)
	require.Equal(t, model.StartingLiglicko2Volatility, game.Liglicko2WhiteOldVolatility)
	require.Equal(t, model.StartingLiglicko2Deviation, game.Liglicko2WhiteOldDeviation)
	require.NotZero(t, game.Liglicko2WhiteOldAt)

	require.Equal(t, model.StartingLiglicko2Rating, game.Liglicko2BlackOldRating)
	require.Equal(t, model.StartingLiglicko2Volatility, game.Liglicko2BlackOldVolatility)
	require.Equal(t, model.StartingLiglicko2Deviation, game.Liglicko2BlackOldDeviation)
	require.NotZero(t, game.Liglicko2BlackOldAt)

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

	// Verify old states are set correctly
	require.Equal(t, model.StartingLiglicko2Rating, game.Liglicko2WhiteOldRating)
	require.Equal(t, model.StartingLiglicko2Volatility, game.Liglicko2WhiteOldVolatility)
	require.Equal(t, model.StartingLiglicko2Deviation, game.Liglicko2WhiteOldDeviation)
	require.NotZero(t, game.Liglicko2WhiteOldAt)

	require.Equal(t, model.StartingLiglicko2Rating, game.Liglicko2BlackOldRating)
	require.Equal(t, model.StartingLiglicko2Volatility, game.Liglicko2BlackOldVolatility)
	require.Equal(t, model.StartingLiglicko2Deviation, game.Liglicko2BlackOldDeviation)
	require.NotZero(t, game.Liglicko2BlackOldAt)

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
	t.Parallel()

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

func TestCreateGameSetsOldStatesCorrectly(t *testing.T) {
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
	ikey1, _ := model.NextIKey(db)
	game1, _, _, err := model.CreateGame(tx, &p1, &p2, true, ikey1, model.Score_Win, r)
	require.NoError(t, err)

	// Capture P1 and P2 states after Game 1
	p1AfterG1 := p1
	p2AfterG1 := p2

	// Game 2: P1 vs P2, Draw
	ikey2, _ := model.NextIKey(db)
	game2, _, _, err := model.CreateGame(tx, &p1, &p2, true, ikey2, model.Score_Draw, r)
	require.NoError(t, err)

	// Verify Game 2's "old" states match states after Game 1
	require.InDelta(t, p1AfterG1.Liglicko2Rating, game2.Liglicko2WhiteOldRating, 1e-9)
	require.InDelta(t, p1AfterG1.Liglicko2Volatility, game2.Liglicko2WhiteOldVolatility, 1e-9)
	require.InDelta(t, p1AfterG1.Liglicko2Deviation, game2.Liglicko2WhiteOldDeviation, 1e-9)
	require.InDelta(t, p1AfterG1.Liglicko2At, game2.Liglicko2WhiteOldAt, 1e-9)

	require.InDelta(t, p2AfterG1.Liglicko2Rating, game2.Liglicko2BlackOldRating, 1e-9)
	require.InDelta(t, p2AfterG1.Liglicko2Volatility, game2.Liglicko2BlackOldVolatility, 1e-9)
	require.InDelta(t, p2AfterG1.Liglicko2Deviation, game2.Liglicko2BlackOldDeviation, 1e-9)
	require.InDelta(t, p2AfterG1.Liglicko2At, game2.Liglicko2BlackOldAt, 1e-9)

	// Also verify Game 1's "old" states were starting ratings
	require.Equal(t, model.StartingLiglicko2Rating, game1.Liglicko2WhiteOldRating)
	require.Equal(t, model.StartingLiglicko2Rating, game1.Liglicko2BlackOldRating)
}
