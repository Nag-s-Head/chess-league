package psqldb_test

import (
	"errors"
	"testing"

	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestConnect(t *testing.T) {
	t.Parallel()
	database := testutils.GetDb(t)
	defer database.Close()

	require.NotNil(t, database)
	require.NoError(t, database.GetSqlxDb().Ping())
}

func TestDoTx(t *testing.T) {
	t.Parallel()

	database := testutils.GetDb(t)
	defer database.Close()

	var done bool
	require.NoError(t, database.DoTx(func(tx *sqlx.Tx) error {
		done = true
		return nil
	}))

	require.True(t, done)
}

func TestDoTxError(t *testing.T) {
	t.Parallel()

	database := testutils.GetDb(t)
	defer database.Close()

	expectedErr := errors.New("Test error")
	err := database.DoTx(func(tx *sqlx.Tx) error {
		return expectedErr
	})
	require.Error(t, err)
	require.True(t, errors.Is(err, expectedErr))
}
