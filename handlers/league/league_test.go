package league_test

import (
	"testing"

	"github.com/Nag-s-Head/chess-league/db/model"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/Nag-s-Head/chess-league/handlers/league"
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

	tpl, err := league.Render(db)
	require.NoError(t, err)
	require.NotNil(t, tpl)
	require.Contains(t, tpl, "Team Groupings")
}
