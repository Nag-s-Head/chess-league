package adminutils

import (
	"bytes"
	"errors"
	"html/template"

	"github.com/Nag-s-Head/chess-league/db/model"
	"github.com/Nag-s-Head/chess-league/handlers/utils"
)

var auditTpl *template.Template = utils.GetTemplate(f, "audit_logs.html")

func RenderAuditLogs(auditLogs []model.AuditLogUiFriendly) (template.HTML, error) {
	w := bytes.NewBuffer(make([]byte, 0))
	err := auditTpl.Execute(w, auditLogs)
	if err != nil {
		return "", errors.Join(errors.New("Could not render audit logs"), err)
	}

	return template.HTML(w.String()), nil
}
