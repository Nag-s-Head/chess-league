package auditlogs

import (
	"html/template"
	"net/http"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
	adminutils "github.com/Nag-s-Head/chess-league/handlers/admin/admin_utils"
)

func Render(db *db.Db) func(http.ResponseWriter, *http.Request, *model.AdminUser) (template.HTML, error) {
	return func(w http.ResponseWriter, r *http.Request, au *model.AdminUser) (template.HTML, error) {
		auditLogs, err := model.GetAuditLogsUiFriendly(db)
		if err != nil {
			return "", err
		}

		return adminutils.RenderAuditLogs(auditLogs)
	}
}
