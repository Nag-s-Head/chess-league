package testutils 

import (
	"testing"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/stretchr/testify/require"
)

func GetDb(t *testing.T) *db.Db {
	t.Helper()

	db, err := db.New()
	require.NoError(t, err)
	return db
}
