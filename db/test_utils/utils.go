package testutils

import (
	"sync"
	"testing"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/stretchr/testify/require"
)

var migrationLock sync.Mutex

func GetDb(t *testing.T) *db.Db {
	t.Helper()

	migrationLock.Lock()
	defer migrationLock.Unlock()
	db, err := db.New()
	require.NoError(t, err)
	return db
}
