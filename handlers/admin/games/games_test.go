package games_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
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

	_, _, _, err = model.CreateGame(tx, &player1, &player2, true, ikey2, model.Score_Draw, &http.Request{
		RemoteAddr: "0.0.0.0",
	})
	require.NoError(t, err)
	require.NoError(t, tx.Commit())

	t.Run("No Search", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/mocked-url", strings.NewReader(""))
		w := httptest.NewRecorder()

		tpl, err := games.Render(db)(w, r, model.NewAdminUser("bob", "bob", "0.0.0.0", "bob"))
		require.NoError(t, err)
		require.Contains(t, tpl, player1.Name)
		require.Contains(t, tpl, player2.Name)
	})

	t.Run("Invalid Search", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/mocked-url?q=%s", url.QueryEscape("(((")), strings.NewReader(""))
		w := httptest.NewRecorder()

		_, err := games.Render(db)(w, r, model.NewAdminUser("bob", "bob", "0.0.0.0", "bob"))
		require.Error(t, err)

		t.Run("Search", func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/mocked-url?q=%s", url.QueryEscape("score=1-0")), strings.NewReader(""))
			w := httptest.NewRecorder()

			tpl, err := games.Render(db)(w, r, model.NewAdminUser("bob", "bob", "0.0.0.0", "bob"))
			require.NoError(t, err)
			require.NotContains(t, string(tpl), "Draw")
		})
	})
}
