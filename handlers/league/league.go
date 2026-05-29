package league

import (
	"bytes"
	"embed"
	"errors"
	"html/template"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
	"github.com/Nag-s-Head/chess-league/handlers/utils"
	"github.com/jmoiron/sqlx"
)

//go:embed league.html
var f embed.FS
var indexTpl *template.Template = utils.GetTemplate(f, "league.html")

type Group struct {
	GroupName string
	Players   []model.Player
}

type Model struct {
	Groups []Group
}

func Render(db *db.Db) (template.HTML, error) {
	var tpl template.HTML
	err := db.DoTx(func(tx *sqlx.Tx) error {
		players, err := model.GetLeaguePlayers(tx)
		if err != nil {
			return err
		}

		model := Model{Groups: []Group{
			{
				GroupName: "A",
				Players:   make([]model.Player, 0),
			},
			{
				GroupName: "B",
				Players:   make([]model.Player, 0),
			},
			{
				GroupName: "C",
				Players:   make([]model.Player, 0),
			},
			{
				GroupName: "D",
				Players:   make([]model.Player, 0),
			},
			{
				GroupName: "E",
				Players:   make([]model.Player, 0),
			},
		}}

		for i, player := range players {
			index := i / 3
			if index >= len(model.Groups) {
				break
			}

			model.Groups[index].Players = append(model.Groups[index].Players, player)
		}

		var buf bytes.Buffer
		err = indexTpl.Execute(&buf, model)
		if err != nil {
			return err
		}

		template, err := template.HTML(buf.String()), nil
		if err != nil {
			return err
		}

		tpl = template
		return nil
	})

	if err != nil {
		return "", errors.Join(errors.New("Cannot render league page"), err)
	}

	return tpl, nil
}
