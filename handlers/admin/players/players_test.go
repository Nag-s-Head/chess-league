package players_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Nag-s-Head/chess-league/db/model"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/Nag-s-Head/chess-league/handlers/admin/players"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestEmptyQuery(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	name := uuid.New().String()
	player := model.NewPlayer(name)
	player.DEPRECATEDElo = 1234
	player.Liglicko2Rating = 1678.2
	require.NoError(t, model.InsertPlayer(db, player))

	r := httptest.NewRequest(http.MethodGet, "/mocked-url", strings.NewReader(""))
	w := httptest.NewRecorder()

	tpl, err := players.Render(db)(w, r, model.NewAdminUser("bob", "bob", "0.0.0.0", "bob"))
	require.NoError(t, err)
	require.NotNil(t, tpl)
	require.Contains(t, string(tpl), "1678 ELO")
	require.True(t, strings.Contains(string(tpl), name))
}

func TestInvalidQuery(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	name := uuid.New().String()
	player := model.NewPlayer(name)
	player.Liglicko2Rating = 1678.2
	require.NoError(t, model.InsertPlayer(db, player))

	r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/mocked-url?q=%s", url.QueryEscape("(((")), strings.NewReader(""))
	w := httptest.NewRecorder()

	_, err := players.Render(db)(w, r, model.NewAdminUser("bob", "bob", "0.0.0.0", "bob"))
	require.Error(t, err)
}

func TestValidQuery(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	name := uuid.New().String()
	player := model.NewPlayer(name)
	player.Liglicko2Rating = 1678.2
	require.NoError(t, model.InsertPlayer(db, player))

	name2 := uuid.New().String()
	player2 := model.NewPlayer(name2)
	require.NoError(t, model.InsertPlayer(db, player2))

	r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/mocked-url?q=%s", url.QueryEscape(name)), strings.NewReader(""))
	w := httptest.NewRecorder()

	tpl, err := players.Render(db)(w, r, model.NewAdminUser("bob", "bob", "0.0.0.0", "bob"))
	require.NoError(t, err)
	require.NotNil(t, tpl)
	require.Contains(t, string(tpl), "1678 ELO")
	require.NotContains(t, string(tpl), name2)
}
