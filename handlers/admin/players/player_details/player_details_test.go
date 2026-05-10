package player_details_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Nag-s-Head/chess-league/db/model"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/Nag-s-Head/chess-league/handlers/admin/players/player_details"
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

	req := httptest.NewRequest(http.MethodGet, "/admin/players/f3529eed-e490-4bc8-af26-7ae84af6b371", nil)
	req.SetPathValue("id", player.Id.String())
	rr := httptest.NewRecorder()

	tpl, err := player_details.Render(db)(rr, req, admin)
	require.NoError(t, err)
	require.NotNil(t, tpl)
	require.Contains(t, string(tpl), "Merge Player Into")
}

func TestPostPlayerDetails_Merger(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	// Setup admin
	admin := model.NewAdminUser("admin", uuid.New().String(), "password", "salt")
	tx, err := db.GetSqlxDb().BeginTxx(t.Context(), nil)
	require.NoError(t, err)
	require.NoError(t, model.InsertAdminUser(tx, *admin))
	require.NoError(t, tx.Commit())

	// Setup players
	target := model.NewPlayer("Target" + uuid.New().String())
	dest := model.NewPlayer("Dest" + uuid.New().String())
	require.NoError(t, model.InsertPlayer(db, target))
	require.NoError(t, model.InsertPlayer(db, dest))

	t.Run("Merge Button Clicked", func(t *testing.T) {
		form := url.Values{}
		form.Set("submit", "merge")
		req := httptest.NewRequest(http.MethodPost, "/admin/players/"+target.Id.String(), strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("id", target.Id.String())
		rr := httptest.NewRecorder()

		player_details.PostPlayerDetails(db)(admin)(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		require.Contains(t, rr.Body.String(), "Merge Target")
		require.Contains(t, rr.Body.String(), "INTO...")
		require.Contains(t, rr.Body.String(), "value=\"merge-select\"")
	})

	t.Run("Destination Selected", func(t *testing.T) {
		form := url.Values{}
		form.Set("submit", "merge-select")
		form.Set("merge-player-dest", dest.Id.String())
		req := httptest.NewRequest(http.MethodPost, "/admin/players/"+target.Id.String(), strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("id", target.Id.String())
		rr := httptest.NewRecorder()

		player_details.PostPlayerDetails(db)(admin)(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		require.Contains(t, rr.Body.String(), "Tick to confirm that you want to merge Target")
		require.Contains(t, rr.Body.String(), "INTO Dest")
		require.Contains(t, rr.Body.String(), "value=\"merge-confirm\"")
		require.Contains(t, rr.Body.String(), dest.Id.String()) // Hidden input check
	})

	t.Run("Confirmed", func(t *testing.T) {
		form := url.Values{}
		form.Set("submit", "merge-confirm")
		form.Set("confirm", "confirmed")
		form.Set("merge-player-dest", dest.Id.String())
		req := httptest.NewRequest(http.MethodPost, "/admin/players/"+target.Id.String(), strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("id", target.Id.String())
		rr := httptest.NewRecorder()

		player_details.PostPlayerDetails(db)(admin)(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		require.Contains(t, rr.Body.String(), "window.location.href = '/admin/players/"+dest.Id.String()+"'")

		// Verify merge in DB
		p, err := model.GetPlayer(db, target.Id)
		require.NoError(t, err)
		require.True(t, p.Deleted)
	})
}
