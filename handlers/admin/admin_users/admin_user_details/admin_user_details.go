package adminuserdetails

import (
	"bytes"
	"embed"
	"html/template"
	"net/http"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
	adminutils "github.com/Nag-s-Head/chess-league/handlers/admin/admin_utils"
	"github.com/Nag-s-Head/chess-league/handlers/utils"
	"github.com/google/uuid"
)

//go:embed admin_user_details.html
var f embed.FS
var indexTpl *template.Template = utils.GetTemplate(f, "admin_user_details.html")

type Model struct {
	model.AdminUser
	AuditLogs template.HTML
}

func Render(db *db.Db) func(http.ResponseWriter, *http.Request, *model.AdminUser) (template.HTML, error) {
	return func(w http.ResponseWriter, r *http.Request, _ *model.AdminUser) (template.HTML, error) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			return "", err
		}

		user, err := model.GetAdminUser(db, id)
		if err != nil {
			return "", err
		}

		auditlogs, err := model.GetAuditLogsUiFriendlyForAdmin(db, id)
		if err != nil {
			return "", err
		}

		renderedAuditLogs, err := adminutils.RenderAuditLogs(auditlogs)
		if err != nil {
			return "", err
		}

		var buf bytes.Buffer
		err = indexTpl.Execute(&buf, Model{
			AdminUser: *user,
			AuditLogs: renderedAuditLogs,
		})
		if err != nil {
			return "", err
		}
		return template.HTML(buf.String()), nil
	}
}
