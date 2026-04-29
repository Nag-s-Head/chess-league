package gamedetails

import (
	"bytes"
	"embed"
	"errors"
	"html/template"
	"net/http"
	"strconv"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
	adminutils "github.com/Nag-s-Head/chess-league/handlers/admin/admin_utils"
	"github.com/Nag-s-Head/chess-league/handlers/utils"
)

//go:embed *.html
var f embed.FS
var indexTpl *template.Template = utils.GetTemplate(f, "game_details.html")

type Model struct {
	AuditLogs template.HTML
	model.GameWithDetails
}

func Render(db *db.Db) func(http.ResponseWriter, *http.Request, *model.AdminUser) (template.HTML, error) {
	return func(w http.ResponseWriter, r *http.Request, au *model.AdminUser) (template.HTML, error) {
		ikeyStr := r.PathValue("ikey")
		ikey, err := strconv.ParseInt(ikeyStr, 10, 64)
		if err != nil {
			return "", errors.New("Cannot read ikey")
		}

		game, err := model.GetGameWithDetails(db, ikey)
		if err != nil {
			return "", err
		}

		auditLogs, err := model.GetAuditLogsUiFriendlyForGame(db, ikey)
		if err != nil {
			return "", err
		}

		renderedAuditLogs, err := adminutils.RenderAuditLogs(auditLogs)
		if err != nil {
			return "", err
		}

		buf := bytes.NewBuffer(nil)
		err = indexTpl.Execute(buf, Model{
			GameWithDetails: game,
			AuditLogs:       renderedAuditLogs,
		})
		if err != nil {
			return "", err
		}

		return template.HTML(buf.String()), nil
	}
}
