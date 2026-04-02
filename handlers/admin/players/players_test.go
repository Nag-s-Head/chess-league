package players_test

import (
	"testing"

	"github.com/Nag-s-Head/chess-league/db/model"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/Nag-s-Head/chess-league/handlers/admin/players"
	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	tpl, err := players.Render(db)(nil, nil, model.NewAdminUser("bob", "bob", "0.0.0.0", "bob"))
	require.NoError(t, err)
	require.NotNil(t, tpl)
}
