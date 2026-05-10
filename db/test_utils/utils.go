package testutils

import (
	"errors"
	"sync"
	"testing"

	"github.com/Nag-s-Head/chess-league/db"
)

var migrationLock sync.Mutex

func getDb(t *testing.T, tries int) *db.Db {
	t.Helper()

	migrationLock.Lock()
	defer migrationLock.Unlock()
	db, err := db.New()

	if errors.Is(err, errors.New("Cannot create migrations table")) {
		t.Log("Migration table createion error - likely due to a race condition. Retrying...")
		return getDb(t, tries+1)
	}
	return db
}

func GetDb(t *testing.T) *db.Db {
	t.Helper()

	db := getDb(t, 0)
	return db
}
