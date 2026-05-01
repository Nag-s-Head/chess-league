package gamedetails

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
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

const (
	swapWinner = "swap-winner"
	setDraw    = "set-draw"
	deleted    = "delete"
)

func PostGameDetails(db *db.Db) func(*model.AdminUser) func(http.ResponseWriter, *http.Request) {
	return func(au *model.AdminUser) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			ikeyStr := r.PathValue("ikey")
			ikey, err := strconv.ParseInt(ikeyStr, 10, 64)
			if err != nil {
				adminutils.RenderError(w, errors.New("Cannot read ikey"))
			}

			err = r.ParseForm()
			if err != nil {
				adminutils.RenderError(w, err)
				return
			}

			submitType := r.Form.Get("submit")
			if !utils.IsConfirmed(r) {
				actionName := ""
				switch submitType {
				case swapWinner:
					actionName = "Swap The Winner"
				case setDraw:
					actionName = "Set The Game To A Draw"
				case deleted:
					actionName = "Delete The Game"
				default:
					adminutils.RenderError(w, fmt.Errorf("%s is not a valid submit type", submitType))
					return
				}

				err = utils.RenderConfirmationPage(w, actionName, submitType)
				if err != nil {
					adminutils.RenderError(w, fmt.Errorf("%s is not a valid submit type", submitType))
				}
				return
			}

			switch submitType {
			case swapWinner:
				err = model.SwapGameWinner(db, au.Id, ikey)
				if err != nil {
					adminutils.RenderError(w, errors.Join(errors.New("Cannot swap game winner"), err))
					return
				}
			case setDraw:
				err = model.SetGameToDraw(db, au.Id, ikey)
				if err != nil {
					adminutils.RenderError(w, errors.Join(errors.New("Cannot set game to a draw"), err))
					return
				}
			case deleted:
				err = model.DeleteGame(db, au.Id, ikey)
				if err != nil {
					adminutils.RenderError(w, errors.Join(errors.New("Cannot delete game"), err))
					return
				}
			default:
				adminutils.RenderError(w, fmt.Errorf("%s is not a valid submit type", submitType))
				return
			}

			w.Write([]byte("<script>window.location.reload();</script>"))
		}
	}
}
