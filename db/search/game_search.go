package search

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
	"github.com/djpiper28/rpg-book/common/normalisation"
	"github.com/djpiper28/rpg-book/common/search/parser"
	sqlsearch "github.com/djpiper28/rpg-book/common/search/sql_search"
)

type gamePlayerNameMapper struct {
	ColumnName string
}

func (g *gamePlayerNameMapper) Map(operator parser.GeneratorOperator, value string) (query string, args []any, err error) {
	return fmt.Sprintf("%s.name_normalised LIKE ?", g.ColumnName), []any{"%" + normalisation.Normalise(value) + "%"}, nil
}

type gamesAnyPlayerNameMapper struct{}

func (g *gamesAnyPlayerNameMapper) Map(operator parser.GeneratorOperator, value string) (query string, args []any, err error) {
	mappedValue := "%" + normalisation.Normalise(value) + "%"
	return "players_white.name_normalised LIKE ? OR players_black.name_normalised LIKE ?", []any{mappedValue, mappedValue}, nil
}

func SearchGames(db db.Db, query string) ([]model.Game, error) {
	tableData := sqlsearch.SqlTableData{
		FieldsToScan: []string{
			"games.ikey",
			"games.player_white",
			"games.player_black",
			"games.played",
			"games.deleted",
			"games.score",
			"games.submitter",
			"games.submit_ip",
			"games.submit_user_agent",
			"games.liglicko2_white",
			"games.liglicko2_black",
		},
		TableName: "games",
		JoinClauses: `
		INNER JOIN players AS players_white ON players_white.id = games.player_white
		INNER JOIN players AS players_black ON players_black.id = games.player_black
		`,
	}

	columnMap := sqlsearch.SqlColmnMap{
		TextColumns: map[string]string{
			"score": "games.score",
		},
		BooleanColumns: map[string]string{
			"deleted": "games.deleted",
		},
		NumberColumns: map[string]string{
			"liglicko2_white": "games.liglicko2_white",
			"liglicko2_black": "games.liglicko2_black",
			"ikey":            "games.ikey",
		},
		CustomColumns: map[string]sqlsearch.CustomColumn{
			"white_player": &gamePlayerNameMapper{ColumnName: "players_white"},
			"black_player": &gamePlayerNameMapper{ColumnName: "players_black"},
			"any_player":   &gamesAnyPlayerNameMapper{},
		},
		BasicQueryColumn:    "any_player",
		BasicQueryOperation: parser.GeneratorOperator_Includes,
	}

	nodes, err := parser.Parse(query)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot parse query"), err)
	}

	query, args, err := sqlsearch.AsSql(nodes, tableData, columnMap)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot generate SQL for query"), err)
	}

	query = db.GetSqlxDb().Rebind(query)

	slog.Debug("Executing query", "query", query, "args", args)
	tx, err := db.GetSqlxDb().BeginTxx(context.Background(), &sql.TxOptions{ReadOnly: true})
	if err != nil {
		slog.Error("Cannot execute query", "query", query, "args", args, "err", err)
		return nil, errors.Join(errors.New("Cannot start read-only transaction"), err)
	}

	defer tx.Rollback()

	games := make([]model.Game, 0)
	err = tx.Select(&games, query, args...)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot execute query SQL"), err)
	}

	return games, nil
}
