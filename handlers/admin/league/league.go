package league

import (
	"bytes"
	"embed"
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
	adminutils "github.com/Nag-s-Head/chess-league/handlers/admin/admin_utils"
	"github.com/Nag-s-Head/chess-league/handlers/utils"
	"github.com/google/uuid"
)

//go:embed league.html
var f embed.FS
var indexTpl *template.Template = utils.GetTemplate(f, "league.html")

func Render(db *db.Db) func(http.ResponseWriter, *http.Request, *model.AdminUser) (template.HTML, error) {
	return func(w http.ResponseWriter, r *http.Request, au *model.AdminUser) (template.HTML, error) {
		players, err := model.GetUiFriendlyLeaguePlayers(db)
		if err != nil {
			return "", err
		}

		var buf bytes.Buffer
		err = indexTpl.Execute(&buf, players)
		if err != nil {
			return "", err
		}

		return template.HTML(buf.String()), nil
	}
}

func PostLeaguePlayers(db *db.Db) func(*model.AdminUser) func(http.ResponseWriter, *http.Request) {
	return func(au *model.AdminUser) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				adminutils.RenderError(w, errors.Join(errors.New("Cannot parse form"), err))
				return
			}

			const playerChoicePrefix = "player-"
			playerIds := make([]uuid.UUID, 0)

			for key := range r.Form {
				if strings.HasPrefix(key, playerChoicePrefix) {
					idRaw := strings.TrimPrefix(key, playerChoicePrefix)
					id, err := uuid.Parse(idRaw)
					if err != nil {
						adminutils.RenderError(w, errors.Join(errors.New("Cannot parse ID of player"), err))
						return
					}

					playerIds = append(playerIds, id)
				}
			}

			err = model.SetLeaguePlayers(db, au.Id, playerIds)
			if err != nil {
				adminutils.RenderError(w, errors.Join(errors.New("Cannot set league players"), err))
				return
			}

			slog.Info("Updated the league players", "admin", au)
			w.Write([]byte("<script>window.location.reload();</script>"))
		}
	}
}
