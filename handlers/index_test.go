package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Nag-s-Head/chess-league/db/model"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/Nag-s-Head/chess-league/handlers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestIndex(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	r := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	w := httptest.NewRecorder()
	handlers.Index(db)(w, r)

	body, err := io.ReadAll(w.Result().Body)
	require.NoError(t, err)
	require.NotEmpty(t, body)
}

func TestPlayerDetails(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	playerName := "Test Player " + uuid.New().String()
	player := model.NewPlayer(playerName)
	err := model.InsertPlayer(db, player)
	require.NoError(t, err)

	t.Run("Valid player ID", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/player/"+player.Id.String(), strings.NewReader(""))
		r.SetPathValue("id", player.Id.String())
		w := httptest.NewRecorder()
		handlers.PlayerDetails(db)(w, r)

		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, w.Body.String(), "Test Player")
	})

	t.Run("Invalid player ID", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/player/not-a-uuid", strings.NewReader(""))
		r.SetPathValue("id", "not-a-uuid")
		w := httptest.NewRecorder()
		handlers.PlayerDetails(db)(w, r)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Non-existent player ID", func(t *testing.T) {
		id := uuid.New().String()
		r := httptest.NewRequest(http.MethodGet, "/player/"+id, strings.NewReader(""))
		r.SetPathValue("id", id)
		w := httptest.NewRecorder()
		handlers.PlayerDetails(db)(w, r)

		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
