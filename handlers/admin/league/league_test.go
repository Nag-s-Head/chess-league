package league_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Nag-s-Head/chess-league/db/model"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/Nag-s-Head/chess-league/handlers/admin/league"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	player1 := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayer(db, player1))

	player2 := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayer(db, player2))

	tpl, err := league.Render(db)(nil, nil, model.NewAdminUser("bob", "bob", "0.0.0.0", "bob"))
	require.NoError(t, err)
	require.NotNil(t, tpl)
	require.True(t, strings.Contains(string(tpl), "Manage League"))
}

func TestPostLeaguePlayers(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	// Create players
	player1 := model.NewPlayer("Alice-" + uuid.New().String())
	require.NoError(t, model.InsertPlayer(db, player1))
	player2 := model.NewPlayer("Bob-" + uuid.New().String())
	require.NoError(t, model.InsertPlayer(db, player2))

	admin := model.NewAdminUser(uuid.New().String(), uuid.New().String(), "test", "test")
	require.NoError(t, db.DoTx(func(tx *sqlx.Tx) error {
		require.NoError(t, model.InsertAdminUser(tx, *admin))
		return nil
	}))

	t.Run("successful update", func(t *testing.T) {
		form := url.Values{}
		form.Add(fmt.Sprintf("player-%s", player1.Id), "on")

		req := httptest.NewRequest(http.MethodPost, "/admin/league", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		handler := league.PostLeaguePlayers(db)(admin)
		handler(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "<script>window.location.reload();</script>")

		// Verify DB update
		require.NoError(t, db.DoTx(func(tx *sqlx.Tx) error {
			leaguePlayers, err := model.GetLeaguePlayers(tx)
			require.NoError(t, err)
			require.Len(t, leaguePlayers, 1)
			assert.Equal(t, player1.Id, leaguePlayers[0].Id)
			return nil
		}))
	})

	t.Run("invalid player id", func(t *testing.T) {
		form := url.Values{}
		form.Add("player-not-a-uuid", "on")

		req := httptest.NewRequest(http.MethodPost, "/admin/league", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		handler := league.PostLeaguePlayers(db)(admin)
		handler(w, req)

		// RenderError should have been called, but what's the status code?
		// RenderError in adminutils doesn't seem to set an error status code explicitly in the provided snippet,
		// but it writes the error template.
		assert.Contains(t, w.Body.String(), "Cannot parse ID of player")
	})

	t.Run("non-existent player id", func(t *testing.T) {
		nonExistentId := uuid.New()
		form := url.Values{}
		form.Add(fmt.Sprintf("player-%s", nonExistentId), "on")

		req := httptest.NewRequest(http.MethodPost, "/admin/league", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		handler := league.PostLeaguePlayers(db)(admin)
		handler(w, req)

		assert.Contains(t, w.Body.String(), "Cannot set league players")
	})
}
