package submitgame_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Nag-s-Head/chess-league/db/model"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	submitgame "github.com/Nag-s-Head/chess-league/handlers/submit_game"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestForm(t *testing.T) {
	tpl, err := submitgame.Render()
	require.NoError(t, err)
	require.NotNil(t, tpl)
}

func TestSubmit(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	name := uuid.New().String()
	err := model.InsertPlayer(db, model.NewPlayer(name))
	require.NoError(t, err)

	t.Run("No form data should error", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/mocked-url", strings.NewReader(""))
		w := httptest.NewRecorder()
		err = submitgame.DoSubmit(db, w, r)
		require.Error(t, err)
	})

	t.Run("Test render of player lookup", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/mocked-url?player-name=%s&player-as=white&other-player-name=not_found&winner=white", name), strings.NewReader(""))

		w := httptest.NewRecorder()
		err = submitgame.DoSubmit(db, w, r)
		require.NoError(t, err)
	})
}
