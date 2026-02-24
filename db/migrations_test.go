package db_test

import (
	"testing"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func getTestDb(t *testing.T) *sqlx.DB {
	t.Helper()

	ret, err := db.InternalConnect()
	require.NoError(t, err)
	return ret
}

func TestMigrations(t *testing.T) {
	d := getTestDb(t)
	migratedDb, err := db.From(d)
	require.NoError(t, err)
	defer migratedDb.Close()

	require.NotNil(t, migratedDb)
}
