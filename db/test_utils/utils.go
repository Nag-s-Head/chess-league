package testutils

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/stretchr/testify/require"
)

var migrationLock sync.Mutex

func getDb(t *testing.T, tries int) *db.Db {
	t.Helper()

	if tries > 20 {
		require.FailNow(t, fmt.Sprintf("Cannot retry anymore due to maximum retries %d", tries))
	}

	migrationLock.Lock()
	defer migrationLock.Unlock()
	db, err := db.New()

	if errors.Is(err, errors.New("Cannot create migrations table")) {
		t.Log("Migration table createion error - likely due to a race condition. Retrying...")
		time.Sleep(time.Second / 10)
		return getDb(t, tries+1)
	}
	return db
}

func GetDb(t *testing.T) *db.Db {
	t.Helper()

	db := getDb(t, 0)
	require.NotNil(t, db)
	return db
}
