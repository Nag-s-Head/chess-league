package adminutils_test

import (
	"testing"

	"github.com/Nag-s-Head/chess-league/db/model"
	adminutils "github.com/Nag-s-Head/chess-league/handlers/admin/admin_utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRenderAuditLogsEmpty(t *testing.T) {
	tpl, err := adminutils.RenderAuditLogs([]model.AuditLogUiFriendly{})
	require.NoError(t, err)
	require.NotEmpty(t, tpl)
}

func TestRenderAuditLogsMany(t *testing.T) {
	tpl, err := adminutils.RenderAuditLogs([]model.AuditLogUiFriendly{
		{
			AuditLog:  *model.NewAuditLog(uuid.New(), "op_name", "op_desc"),
			AdminName: "admin_name",
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, tpl)
	require.Contains(t, tpl, "admin_name")
	require.Contains(t, tpl, "op_name")
	require.Contains(t, tpl, "op_desc")
}
