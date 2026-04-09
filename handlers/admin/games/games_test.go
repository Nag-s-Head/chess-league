package games_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/Nag-s-Head/chess-league/db/model"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/Nag-s-Head/chess-league/handlers/admin/games"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	player1 := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayer(db, player1))

	player2 := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayer(db, player2))

	ikey1, err := model.NextIKey(db)
	require.NoError(t, err)

	ikey2, err := model.NextIKey(db)
	require.NoError(t, err)

	tx, err := db.GetSqlxDb().BeginTxx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()

	_, _, _, err = model.CreateGame(tx, &player1, &player2, true, ikey1, model.Score_Win, &http.Request{
		RemoteAddr: "0.0.0.0",
	})
	require.NoError(t, err)

	_, _, _, err = model.CreateGame(tx, &player1, &player2, true, ikey2, model.Score_Win, &http.Request{
		RemoteAddr: "0.0.0.0",
	})
	require.NoError(t, err)
	require.NoError(t, tx.Commit())

	tpl, err := games.Render(db)(nil, nil, model.NewAdminUser("bob", "bob", "0.0.0.0", "bob"))
	require.NoError(t, err)
	require.NotNil(t, tpl)
	require.True(t, strings.Contains(string(tpl), "ELO"))
}
