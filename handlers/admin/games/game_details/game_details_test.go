package gamedetails_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Nag-s-Head/chess-league/db/model"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	game_details "github.com/Nag-s-Head/chess-league/handlers/admin/games/game_details"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	tx, err := db.GetSqlxDb().BeginTxx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()

	admin := model.NewAdminUser(uuid.New().String(), uuid.New().String(), "0.0.0.0", "bob")
	require.NoError(t, model.InsertAdminUser(tx, *admin))

	player := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayerTx(tx, player))
	require.NoError(t, tx.Commit())

	player2 := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayer(db, player2))

	ikey, err := model.NextIKey(db)
	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodGet, "/mocked-url", nil)
	tx2, err := db.GetSqlxDb().BeginTxx(t.Context(), nil)
	require.NoError(t, err)
	game, _, _, err := model.CreateGame(tx2, &player, &player2, true, ikey, model.Score_Win, r)
	require.NoError(t, err)
	require.NoError(t, tx2.Commit())

	req := httptest.NewRequest(http.MethodGet, "/admin/games/details/1", nil)
	req.SetPathValue("ikey", strconv.FormatInt(game.IKey, 10))
	rr := httptest.NewRecorder()

	tpl, err := game_details.Render(db)(rr, req, admin)
	require.NoError(t, err)
	require.NotNil(t, tpl)
}
