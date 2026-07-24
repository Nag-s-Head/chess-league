package players

import (
	"bytes"
	"embed"
	"html/template"
	"net/http"
	"strings"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
	"github.com/Nag-s-Head/chess-league/db/search"
	"github.com/Nag-s-Head/chess-league/handlers/utils"
)

//go:embed players.html
var f embed.FS
var indexTpl *template.Template = utils.GetTemplate(f, "players.html")

type Model struct {
	Query   string
	Players []model.PlayerWithGameCount
}

func Render(db db.Db) func(http.ResponseWriter, *http.Request, *model.AdminUser) (template.HTML, error) {
	return func(w http.ResponseWriter, r *http.Request, au *model.AdminUser) (template.HTML, error) {
		var m Model

		users, err := model.GetPlayersByEloWithGameCount(db)
		if err != nil {
			return "", err
		}

		m.Players = users
		queries, ok := r.URL.Query()["q"]
		if ok {
			m.Query = queries[0]
		}

		if strings.Trim(m.Query, " ") != "" {
			searchResults, err := search.SearchPlayers(db, m.Query)
			if err != nil {
				return "", err
			}

			m.Players = make([]model.PlayerWithGameCount, 0)
			for _, user := range users {
				for _, searchRes := range searchResults {
					if user.Id == searchRes.Id {
						m.Players = append(m.Players, user)
						break
					}
				}
			}
		}

		var buf bytes.Buffer
		err = indexTpl.Execute(&buf, m)
		if err != nil {
			return "", err
		}
		return template.HTML(buf.String()), nil
	}
}
