package players_test

import (
	"strings"
	"testing"

	"github.com/Nag-s-Head/chess-league/db/model"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/Nag-s-Head/chess-league/handlers/admin/players"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	name := uuid.New().String()
	player := model.NewPlayer(name)
	player.DEPRECATEDElo = 1234
	player.Liglicko2Rating = 1678.2
	require.NoError(t, model.InsertPlayer(db, player))

	tpl, err := players.Render(db)(nil, nil, model.NewAdminUser("bob", "bob", "0.0.0.0", "bob"))
	require.NoError(t, err)
	require.NotNil(t, tpl)
	require.Contains(t, string(tpl), "1678 ELO")
	require.True(t, strings.Contains(string(tpl), name))
}
