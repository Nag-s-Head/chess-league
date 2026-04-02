package adminuserdetails_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nag-s-Head/chess-league/db/model"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	adminuserdetails "github.com/Nag-s-Head/chess-league/handlers/admin/admin_users/admin_user_details"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	user := model.NewAdminUser(uuid.New().String(), uuid.New().String(), "0.0.0.0", "bob")
	tx, err := db.GetSqlxDb().BeginTxx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()

	require.NoError(t, model.InsertAdminUser(tx, *user))
	require.NoError(t, tx.Commit())

	req := httptest.NewRequest(http.MethodGet, "/admin/admins/042e73dc-8c82-4c9a-9aa4-0c0d593c1faa", nil)
	req.SetPathValue("id", user.Id.String())
	rr := httptest.NewRecorder()

	tpl, err := adminuserdetails.Render(db)(rr, req, user)
	require.NoError(t, err)
	require.NotNil(t, tpl)
}
