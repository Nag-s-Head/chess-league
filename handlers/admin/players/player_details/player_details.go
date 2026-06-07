package player_details

import (
	"bytes"
	"embed"
	"errors"
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
var mergeSelectionTpl *template.Template = utils.GetTemplate(f, "merge_selection.html")

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

		buf := bytes.NewBuffer(nil)
		err = indexTpl.Execute(buf, Model{
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

			player, err := model.GetPlayer(db, id)
			if err != nil {
				adminutils.RenderError(w, errors.Join(errors.New("Cannot get player"), err))
				return
			}

			err = r.ParseForm()
			if err != nil {
				adminutils.RenderError(w, err)
				return
			}

			submitType := r.Form.Get("submit")
			switch submitType {
			case "merge":
				players, err := model.GetPlayers(db)
				if err != nil {
					adminutils.RenderError(w, err)
					return
				}
				target, err := model.GetPlayer(db, id)
				if err != nil {
					adminutils.RenderError(w, err)
					return
				}

				otherPlayers := make([]model.Player, 0)
				for _, p := range players {
					if p.Id != id && !p.Deleted {
						otherPlayers = append(otherPlayers, p)
					}
				}

				err = mergeSelectionTpl.Execute(w, struct {
					TargetName string
					Players    []model.Player
				}{
					TargetName: target.Name,
					Players:    otherPlayers,
				})
				if err != nil {
					adminutils.RenderError(w, err)
					return
				}
			case "merge-select":
				destIdStr := r.Form.Get("merge-player-dest")
				destId, err := uuid.Parse(destIdStr)
				if err != nil {
					adminutils.RenderError(w, err)
					return
				}

				target, err := model.GetPlayer(db, id)
				if err != nil {
					adminutils.RenderError(w, err)
					return
				}
				dest, err := model.GetPlayer(db, destId)
				if err != nil {
					adminutils.RenderError(w, err)
					return
				}

				action := fmt.Sprintf("merge %s INTO %s such that only %s is left", target.Name, dest.Name, dest.Name)
				err = utils.RenderConfirmationPage(w, action, "merge-confirm")
				if err != nil {
					adminutils.RenderError(w, err)
					return
				}
				w.Write([]byte(fmt.Sprintf("<input type='hidden' name='merge-player-dest' value='%s'>", destIdStr)))
			case "merge-confirm":
				if !utils.IsConfirmed(r) {
					adminutils.RenderError(w, errors.New("Not confirmed"))
					return
				}

				destIdStr := r.Form.Get("merge-player-dest")
				destId, err := uuid.Parse(destIdStr)
				if err != nil {
					adminutils.RenderError(w, err)
					return
				}

				err = model.MergePlayers(db, id, destId, au.Id)
				if err != nil {
					adminutils.RenderError(w, err)
					return
				}

				w.Write([]byte(fmt.Sprintf("<script>window.location.href = '/admin/players/%s'</script>", destId.String())))
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
				err = utils.RenderConfirmationPage(w, fmt.Sprintf("Delete player %s", player.Name), "delete-confirm")
				if err != nil {
					adminutils.RenderError(w, err)
					return
				}
			case "delete-confirm":
				err = model.DeletePlayer(db, id, au.Id)
				if err != nil {
					adminutils.RenderError(w, err)
					return
				}
				w.Write([]byte("<script>window.location.reload();</script>"))
			default:
				adminutils.RenderError(w, fmt.Errorf("%s is not a valid submit type", submitType))
				return
			}
		}
	}
}
