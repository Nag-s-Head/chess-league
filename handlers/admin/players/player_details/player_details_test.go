package player_details_test

import (
	"net/http"
	"net/http/httptest"
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
}
