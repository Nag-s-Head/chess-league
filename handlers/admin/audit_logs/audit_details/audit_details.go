package auditdetails

import (
	"bytes"
	"embed"
	"errors"
	"html/template"
	"net/http"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
	"github.com/Nag-s-Head/chess-league/handlers/utils"
	"github.com/google/uuid"
)

//go:embed *.html
var f embed.FS
var indexTpl *template.Template = utils.GetTemplate(f, "audit_details.html")

func Render(db db.Db) func(http.ResponseWriter, *http.Request, *model.AdminUser) (template.HTML, error) {
	return func(w http.ResponseWriter, r *http.Request, au *model.AdminUser) (template.HTML, error) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			return "", errors.New("Cannot read audit log id")
		}

		tx, err := db.GetSqlxDb().BeginTxx(r.Context(), nil)
		if err != nil {
			return "", err
		}
		defer tx.Rollback()

		detailedAuditLog, err := model.GetAuditLog(tx, id)
		if err != nil {
			return "", err
		}

		buf := bytes.NewBuffer(nil)
		err = indexTpl.Execute(buf, detailedAuditLog)
		if err != nil {
			return "", err
		}

		return template.HTML(buf.String()), nil
	}
}
