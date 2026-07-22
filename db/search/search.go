package search

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
	"github.com/djpiper28/rpg-book/common/normalisation"
	"github.com/djpiper28/rpg-book/common/search/parser"
	sqlsearch "github.com/djpiper28/rpg-book/common/search/sql_search"
)

type playerNameNormMapper struct{}

func (p playerNameNormMapper) Map(operator parser.GeneratorOperator, value string) (query string, args []any, err error) {
	return "players.name_normalised LIKE ?", []any{"%" + normalisation.Normalise(value) + "%"}, nil
}

func SearchPlayers(db db.Db, query string) ([]model.Player, error) {
	tableData := sqlsearch.SqlTableData{
		FieldsToScan: []string{
			"id",
			"name",
			"liglicko2_rating",
			"liglicko2_deviation",
			"liglicko2_volatility",
			"liglicko2_at",
			"join_time",
			"deleted",
		},
		TableName: "players",
	}

	columnMap := sqlsearch.SqlColmnMap{
		TextColumns: map[string]string{
			"name": "players.name",
		},
		BooleanColumns: map[string]string{
			"deleted": "deleted",
		},
		NumberColumns: map[string]string{
			"rating":     "liglicko2_rating",
			"deviation":  "liglicko2_deviation",
			"volatility": "liglicko2_volatility",
		},
		CustomColumns: map[string]sqlsearch.CustomColumn{
			"name_norm": playerNameNormMapper{},
		},
		BasicQueryColumn:    "name_norm",
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

	var players []model.Player
	err = tx.Select(&players, query, args...)
	if err != nil {
		return nil, errors.Join(errors.New("Cannot execute query SQL"), err)
	}

	return players, nil
}
