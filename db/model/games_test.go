package model_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

	p1 := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayer(db, p1))

	p2 := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayer(db, p2))

	ikey, err := model.NextIKey(db)
	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodGet, "/mocked-url", strings.NewReader(""))
	game, err := model.CreateGame(db, &p1, &p2, true, ikey, model.Score_Win, r)
	require.NoError(t, err)

	require.NotEqual(t, model.StartingElo, p1.Elo)
	require.NotEqual(t, model.StartingElo, p2.Elo)

	require.Equal(t, p1.Id, game.Submitter)
	require.Equal(t, p1.Id, game.PlayerWhite)
	require.Equal(t, p2.Id, game.PlayerBlack)

	require.Equal(t, p1.Elo-model.StartingElo, game.EloGiven)
	require.Equal(t, p2.Elo-model.StartingElo, game.EloTaken)

	require.Equal(t, model.Score_Win, game.Score)
	require.Equal(t, false, game.Deleted)

	require.Equal(t, ikey, game.IKey)

	p1New, err := model.GetPlayer(db, p1.Id)
	require.NoError(t, err)
	require.Equal(t, p1.Elo, p1New.Elo)

	p2New, err := model.GetPlayer(db, p2.Id)
	require.NoError(t, err)
	require.Equal(t, p2.Elo, p2New.Elo)
}

func TestCreateGameP1Black(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	p1 := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayer(db, p1))

	p2 := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayer(db, p2))

	ikey, err := model.NextIKey(db)
	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodGet, "/mocked-url", strings.NewReader(""))
	game, err := model.CreateGame(db, &p1, &p2, false, ikey, model.Score_Win, r)
	require.NoError(t, err)

	require.NotEqual(t, model.StartingElo, p1.Elo)
	require.NotEqual(t, model.StartingElo, p2.Elo)

	require.Equal(t, p1.Id, game.Submitter)
	require.Equal(t, p2.Id, game.PlayerWhite)
	require.Equal(t, p1.Id, game.PlayerBlack)

	require.Equal(t, p2.Elo-model.StartingElo, game.EloGiven)
	require.Equal(t, p1.Elo-model.StartingElo, game.EloTaken)

	require.Equal(t, model.Score_Win, game.Score)
	require.Equal(t, false, game.Deleted)

	require.Equal(t, ikey, game.IKey)

	p1New, err := model.GetPlayer(db, p1.Id)
	require.NoError(t, err)
	require.Equal(t, p1.Elo, p1New.Elo)

	p2New, err := model.GetPlayer(db, p2.Id)
	require.NoError(t, err)
	require.Equal(t, p2.Elo, p2New.Elo)
}

func TestCreateGameSameIkeyFails(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	p1 := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayer(db, p1))

	p2 := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayer(db, p2))

	ikey, err := model.NextIKey(db)
	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodGet, "/mocked-url", strings.NewReader(""))
	_, err = model.CreateGame(db, &p1, &p2, true, ikey, model.Score_Win, r)
	require.NoError(t, err)

	_, err = model.CreateGame(db, &p1, &p2, true, ikey, model.Score_Win, r)
	require.Error(t, err)
}
