package gamedetails_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
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
	admin := model.NewAdminUser(uuid.New().String(), uuid.New().String(), "0.0.0.0", "bob")
	require.NoError(t, model.InsertAdminUser(tx, *admin))
	require.NoError(t, tx.Commit())

	player := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayer(db, player))

	player2 := model.NewPlayer(uuid.New().String())
	require.NoError(t, model.InsertPlayer(db, player2))

	t.Run("Win", func(t *testing.T) {
		ikey, err := model.NextIKey(db)
		require.NoError(t, err)

		r := httptest.NewRequest(http.MethodGet, "/mocked-url", nil)
		tx, err := db.GetSqlxDb().BeginTxx(t.Context(), nil)
		require.NoError(t, err)
		game, _, _, err := model.CreateGame(tx, &player, &player2, true, ikey, model.Score_Win, r)
		require.NoError(t, err)
		require.NoError(t, tx.Commit())

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.SetPathValue("ikey", strconv.FormatInt(game.IKey, 10))
		rr := httptest.NewRecorder()

		tpl, err := game_details.Render(db)(rr, req, admin)
		require.NoError(t, err)
		require.NotNil(t, tpl)
		require.Contains(t, string(tpl), `value="swap-winner"`)
		require.Contains(t, string(tpl), `value="set-draw"`)
		require.Contains(t, string(tpl), `value="delete"`)
	})

	t.Run("Draw", func(t *testing.T) {
		ikey, err := model.NextIKey(db)
		require.NoError(t, err)

		r := httptest.NewRequest(http.MethodGet, "/mocked-url", nil)
		tx, err := db.GetSqlxDb().BeginTxx(t.Context(), nil)
		require.NoError(t, err)
		game, _, _, err := model.CreateGame(tx, &player, &player2, true, ikey, model.Score_Draw, r)
		require.NoError(t, err)
		require.NoError(t, tx.Commit())

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.SetPathValue("ikey", strconv.FormatInt(game.IKey, 10))
		rr := httptest.NewRecorder()

		tpl, err := game_details.Render(db)(rr, req, admin)
		require.NoError(t, err)
		require.NotNil(t, tpl)
		require.NotContains(t, string(tpl), `value="swap-winner"`)
		require.NotContains(t, string(tpl), `value="set-draw"`)
		require.Contains(t, string(tpl), `value="delete"`)
	})

	t.Run("Deleted", func(t *testing.T) {
		ikey, err := model.NextIKey(db)
		require.NoError(t, err)

		r := httptest.NewRequest(http.MethodGet, "/mocked-url", nil)
		tx, err := db.GetSqlxDb().BeginTxx(t.Context(), nil)
		require.NoError(t, err)
		game, _, _, err := model.CreateGame(tx, &player, &player2, true, ikey, model.Score_Win, r)
		require.NoError(t, err)
		require.NoError(t, tx.Commit())

		require.NoError(t, model.DeleteGame(db, admin.Id, game.IKey))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.SetPathValue("ikey", strconv.FormatInt(game.IKey, 10))
		rr := httptest.NewRecorder()

		tpl, err := game_details.Render(db)(rr, req, admin)
		require.NoError(t, err)
		require.NotNil(t, tpl)
		require.NotContains(t, string(tpl), `value="swap-winner"`)
		require.NotContains(t, string(tpl), `value="set-draw"`)
		require.NotContains(t, string(tpl), `value="delete"`)
		require.Contains(t, string(tpl), `value="undo-delete"`)
		require.Contains(t, string(tpl), "DELETED")
	})
}

func TestPostGameDetailsUndeleteRequiresConfirmation(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	tx, err := db.GetSqlxDb().BeginTxx(t.Context(), nil)
	require.NoError(t, err)
	admin := model.NewAdminUser(uuid.NewString(), uuid.NewString(), "127.0.0.1", "Go test")
	require.NoError(t, model.InsertAdminUser(tx, *admin))
	require.NoError(t, tx.Commit())

	white := model.NewPlayer("undelete-handler-white-" + uuid.NewString())
	require.NoError(t, model.InsertPlayer(db, white))
	black := model.NewPlayer("undelete-handler-black-" + uuid.NewString())
	require.NoError(t, model.InsertPlayer(db, black))

	ikey, err := model.NextIKey(db)
	require.NoError(t, err)
	createTx, err := db.GetSqlxDb().BeginTxx(t.Context(), nil)
	require.NoError(t, err)
	createRequest := httptest.NewRequest(http.MethodPost, "/", nil)
	game, _, _, err := model.CreateGame(createTx, &white, &black, true, ikey, model.Score_Win, createRequest)
	require.NoError(t, err)
	require.NoError(t, createTx.Commit())
	require.NoError(t, model.DeleteGame(db, admin.Id, game.IKey))

	post := game_details.PostGameDetails(db)(admin)
	confirmationForm := url.Values{"submit": {"undo-delete"}}
	confirmationRequest := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(confirmationForm.Encode()))
	confirmationRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	confirmationRequest.SetPathValue("ikey", strconv.FormatInt(game.IKey, 10))
	confirmationResponse := httptest.NewRecorder()
	post(confirmationResponse, confirmationRequest)

	require.Contains(t, confirmationResponse.Body.String(), "Tick to confirm")
	require.Contains(t, confirmationResponse.Body.String(), `value="undo-delete"`)
	stillDeleted, err := model.GetGameWithDetails(db, game.IKey)
	require.NoError(t, err)
	require.True(t, stillDeleted.Deleted)

	confirmedForm := url.Values{
		"submit":  {"undo-delete"},
		"confirm": {"confirmed"},
	}
	confirmedRequest := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(confirmedForm.Encode()))
	confirmedRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	confirmedRequest.SetPathValue("ikey", strconv.FormatInt(game.IKey, 10))
	confirmedResponse := httptest.NewRecorder()
	post(confirmedResponse, confirmedRequest)

	require.Equal(t, `<script>window.location.reload();</script>`, confirmedResponse.Body.String())
	undeleted, err := model.GetGameWithDetails(db, game.IKey)
	require.NoError(t, err)
	require.False(t, undeleted.Deleted)
}
