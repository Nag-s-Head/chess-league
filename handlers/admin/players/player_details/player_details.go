package player_details

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
	adminutils "github.com/Nag-s-Head/chess-league/handlers/admin/admin_utils"
	"github.com/Nag-s-Head/chess-league/handlers/utils"
	"github.com/google/uuid"
)

//go:embed *.html
var f embed.FS
var indexTpl *template.Template = utils.GetTemplate(f, "player_details.html")
var renameFormTpl *template.Template = utils.GetTemplate(f, "rename_form.html")

type Model struct {
	Player    model.Player
	Details   model.GamesUiFriendly
	AuditLogs template.HTML
}

func Render(db *db.Db) func(http.ResponseWriter, *http.Request, *model.AdminUser) (template.HTML, error) {
	return func(w http.ResponseWriter, r *http.Request, au *model.AdminUser) (template.HTML, error) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			return "", err
		}

		player, err := model.GetPlayer(db, id)
		if err != nil {
			return "", err
		}

		games, err := model.GetGamesByPlayer(db, id)
		if err != nil {
			return "", err
		}

		details := model.MapGamesToUserFriendly(id, games)
		auditLogs, err := model.GetAuditLogsUiFriendlyForPlayer(db, id)
		if err != nil {
			return "", err
		}

		renderedAuditLogs, err := adminutils.RenderAuditLogs(auditLogs)
		if err != nil {
			return "", err
		}

		var buf bytes.Buffer
		err = indexTpl.Execute(&buf, Model{
			Player:    player,
			Details:   details,
			AuditLogs: renderedAuditLogs,
		})
		if err != nil {
			return "", err
		}
		return template.HTML(buf.String()), nil
	}
}

func PostPlayerDetails(db *db.Db) func(*model.AdminUser) func(http.ResponseWriter, *http.Request) {
	return func(au *model.AdminUser) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			idStr := r.PathValue("id")
			id, err := uuid.Parse(idStr)
			if err != nil {
				adminutils.RenderError(w, err)
				return
			}

			err = r.ParseForm()
			if err != nil {
				adminutils.RenderError(w, err)
				return
			}

			submitType := r.Form.Get("submit")
			switch submitType {
			case "rename":
				newName := r.Form.Get("player-name")
				if newName == "" {
					err := renameFormTpl.Execute(w, nil)
					if err != nil {
						adminutils.RenderError(w, err)
						return
					}
					return
				} else {
					err := model.RenamePlayer(db, id, newName, au.Id)
					if err != nil {
						adminutils.RenderError(w, err)
						return
					}

					w.Write([]byte(`
	<p class="text-green-500 font-bold">Success, reloading...</p>
	<script>window.location.reload()</script>
						`))
				}
			case "delete":
			default:
				adminutils.RenderError(w, fmt.Errorf("%s is not a valid submit type", submitType))
				return
			}
		}
	}
}
