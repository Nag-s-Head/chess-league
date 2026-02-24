package db_test

import (
	"testing"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/stretchr/testify/require"
)

func TestConnect(t *testing.T) {
	database, err := db.InternalConnect()
	require.NoError(t, err)
	defer database.Close()

	require.NotNil(t, database)
	require.NoError(t, database.Ping())
}
