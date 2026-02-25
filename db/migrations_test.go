package db_test

import (
	"testing"

	"github.com/Nag-s-Head/chess-league/db"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
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

func TestFrom(t *testing.T) {
	d := testutils.GetDb(t)
	d.Close()
}
