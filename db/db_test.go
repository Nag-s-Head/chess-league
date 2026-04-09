package db_test

import (
	"testing"

	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/stretchr/testify/require"
)

func TestConnect(t *testing.T) {
	database := testutils.GetDb(t)
	defer database.Close()

	require.NotNil(t, database)
	require.NoError(t, database.GetSqlxDb().Ping())
}
