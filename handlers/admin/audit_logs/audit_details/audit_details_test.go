package auditdetails_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nag-s-Head/chess-league/db/model"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	auditdetails "github.com/Nag-s-Head/chess-league/handlers/admin/audit_logs/audit_details"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	tx, err := db.GetSqlxDb().BeginTxx(t.Context(), nil)
	require.NoError(t, err)
	admin := model.NewAdminUser("bob", uuid.New().String(), "0.0.0.0", "bob")
	require.NoError(t, model.InsertAdminUser(tx, *admin))

	auditLog := model.NewAuditLog(admin.Id, "Test Op", "Test Desc")
	require.NoError(t, model.InsertAuditLog(tx, auditLog))
	require.NoError(t, tx.Commit())

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.SetPathValue("id", auditLog.Id.String())
	rr := httptest.NewRecorder()

	tpl, err := auditdetails.Render(db)(rr, req, admin)
	require.NoError(t, err)
	require.NotNil(t, tpl)
	require.Contains(t, string(tpl), "Test Op")
	require.Contains(t, string(tpl), "Test Desc")
	require.Contains(t, string(tpl), "bob")
}
