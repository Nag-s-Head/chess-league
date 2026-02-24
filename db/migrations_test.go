package db_test

import (
	"testing"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/stretchr/testify/require"
)

func TestMigrations(t *testing.T) {
	d, err := db.InternalConnect()
	require.NoError(t, err)
	migratedDb, err := db.From(d)
	require.NoError(t, err)
	defer migratedDb.Close()

	require.NotNil(t, migratedDb)
}

func getDb(t *testing.T) *db.Db {
	t.Helper()

	db, err := db.New()
	require.NoError(t, err)
	return db
}

func TestFrom(t *testing.T) {
	d := getDb(t)
	d.Close()
}
