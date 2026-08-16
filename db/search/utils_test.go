package search_test

import (
	"database/sql"
	"errors"
	"strings"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func assertQueryResultAfterFuzz(t *testing.T, err error) {
	t.Helper()

	pqErr, _ := errors.AsType[*pq.Error](err)
	if pqErr != nil ||
		errors.Is(err, sql.ErrNoRows) ||
		errors.Is(err, sql.ErrTxDone) ||
		strings.Contains(err.Error(), "sql:") {

		// Cases where the SQL error is fine
		if strings.Contains(err.Error(), "invalid input syntax for type") {
			return
		}

		require.FailNow(t, "Unexpected pg error", err)
	}
}
