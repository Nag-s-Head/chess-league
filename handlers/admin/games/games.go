package games

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

type Model struct {
	Query string
	Games []model.GameWithOutcome
}

//go:embed games.html
var f embed.FS
var indexTpl *template.Template = utils.GetTemplate(f, "games.html")

func Render(db db.Db) func(http.ResponseWriter, *http.Request, *model.AdminUser) (template.HTML, error) {
	return func(w http.ResponseWriter, r *http.Request, au *model.AdminUser) (template.HTML, error) {
		var m Model

		games, err := model.GetGamesWithOutcomes(db)
		if err != nil {
			return "", err
		}

		m.Games = games

		queries, ok := r.URL.Query()["q"]
		if ok {
			m.Query = queries[0]
		}

		if strings.Trim(m.Query, " ") != "" {
			searchResults, err := search.SearchGames(db, m.Query)
			if err != nil {
				return "", err
			}

			m.Games = make([]model.GameWithOutcome, 0)
			for _, game := range games {
				for _, searchRes := range searchResults {
					if game.Ikey == searchRes.IKey {
						m.Games = append(m.Games, game)
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
