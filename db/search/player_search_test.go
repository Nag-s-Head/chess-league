package search_test

import (
	"log/slog"
	"strings"
	"sync"
	"testing"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/liglicko2"
	"github.com/Nag-s-Head/chess-league/db/model"
	"github.com/Nag-s-Head/chess-league/db/search"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSearchPlayers(t *testing.T) {
	t.Parallel()

	db := testutils.GetDb(t)
	defer db.Close()

	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Dariuz"+uuid.NewString())))
	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Chas"+uuid.NewString())))
	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Danny"+uuid.NewString())))
	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Greg"+uuid.NewString())))
	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Sophie"+uuid.NewString())))
	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Anne"+uuid.NewString())))
	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Charlotte"+uuid.NewString())))
	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Andy Burnham"+uuid.NewString())))

	beryl := model.NewPlayer("Beryl" + uuid.NewString())
	beryl.Deleted = true
	beryl.Liglicko2Rating = 3900
	require.NoError(t, model.InsertPlayer(db, beryl))

	t.Run("Default Search", func(t *testing.T) {
		players, err := search.SearchPlayers(db, "dariuz")
		require.NoError(t, err)
		require.NotEmpty(t, players)
		require.Contains(t, players[0].Name, "Dariuz")
	})

	t.Run("Default Search Fuzzy", func(t *testing.T) {
		players, err := search.SearchPlayers(db, "dar")
		require.NoError(t, err)
		require.NotEmpty(t, players)
		require.Contains(t, players[0].Name, "Dariuz")
	})

	t.Run("Test Rating Alias", func(t *testing.T) {
		players, err := search.SearchPlayers(db, "rating>=3900")
		require.NoError(t, err)
		require.NotEmpty(t, players)
		require.Contains(t, players[0].Name, "Beryl")
	})

	t.Run("Test Deleted Column", func(t *testing.T) {
		players, err := search.SearchPlayers(db, "deleted=true")
		require.NoError(t, err)
		require.NotEmpty(t, players)
		require.True(t, players[0].Deleted)
	})

	t.Run("No Results", func(t *testing.T) {
		players, err := search.SearchPlayers(db, "deleted=true and deleted=false")
		require.NoError(t, err)
		require.NotNil(t, players)
		require.Len(t, players, 0)
	})

	t.Run("Test Advanced Query", func(t *testing.T) {
		// Sanity checks
		require.Less(t, 600.0, liglicko2.DefaultRating, "Test sanity check")
		require.Greater(t, 1600.0, liglicko2.DefaultRating, "Test sanity check")

		players, err := search.SearchPlayers(db, `deleted=false and (name_norm=greg OR name_norm="chas") and (rating>600 and rating<1600)`)
		require.NoError(t, err)
		require.NotEmpty(t, players)
		require.True(t, strings.Contains(players[0].Name, "Greg") || strings.Contains(players[0].Name, "Chas"))
		require.False(t, players[0].Deleted)
	})
}

func FuzzSearchPlayer(f *testing.F) {
	// Search uses log debug for queries, this makes it far easier to reason with test failures
	slog.SetLogLoggerLevel(slog.LevelDebug)

	f.Add("'; DELETE FROM players; --")

	f.Add("Dave")
	f.Add("John Major")
	f.Add("deleted=1")
	f.Add("deleted=t")
	f.Add("deleted=T")
	f.Add("deleted=true")
	f.Add("deleted=TrUe")

	f.Add("deleted=0")
	f.Add("deleted=f")
	f.Add("deleted=F")
	f.Add("deleted=false")
	f.Add("deleted=FaLsE")

	f.Add("rating>123")
	f.Add("rating>=124")
	f.Add("rating<1000")
	f.Add("rating<=999")
	f.Add("rating=999")
	f.Add("rating~999")

	f.Add("deviation=123")
	f.Add("name=Andy Burnham")
	f.Add("name~Andy Burnham")
	f.Add(`name~"Andy Burnham"`)

	f.Add(`deleted=false and (name_norm=greg OR name_norm="chas") and (rating>600 and rating<1600)`)
	f.Add(`(name_norm=greg OR name_norm="chas") and rating>600 and rating<1600`)

	// Here are some of the fun cases fuzzing raised that have since been fixed
	f.Add("(((((((((((((((((( 0")

	var isSetup bool
	var lock sync.Mutex
	var dbInstance db.Db

	setup := func(t *testing.T) db.Db {
		t.Helper()

		lock.Lock()
		defer lock.Unlock()

		if dbInstance == nil {
			dbInstance = testutils.GetDb(t)
		}

		if isSetup {
			return dbInstance
		}

		isSetup = true

		require.NoError(t, model.InsertPlayer(dbInstance, model.NewPlayer("Dave"+uuid.NewString())))
		require.NoError(t, model.InsertPlayer(dbInstance, model.NewPlayer("Chas"+uuid.NewString())))
		require.NoError(t, model.InsertPlayer(dbInstance, model.NewPlayer("Danny"+uuid.NewString())))
		require.NoError(t, model.InsertPlayer(dbInstance, model.NewPlayer("Greg"+uuid.NewString())))
		require.NoError(t, model.InsertPlayer(dbInstance, model.NewPlayer("Sophie"+uuid.NewString())))
		require.NoError(t, model.InsertPlayer(dbInstance, model.NewPlayer("Anne"+uuid.NewString())))
		require.NoError(t, model.InsertPlayer(dbInstance, model.NewPlayer("Charlotte"+uuid.NewString())))
		require.NoError(t, model.InsertPlayer(dbInstance, model.NewPlayer("Andy Burnham"+uuid.NewString())))

		return dbInstance
	}

	f.Fuzz(func(t *testing.T, query string) {
		t.Logf("Query is: %s", query)
		db := setup(t)

		players, err := search.SearchPlayers(db, query)
		if err == nil {
			require.NotNil(t, players)
			return
		}

		assertQueryResultAfterFuzz(t, err)
	})
}
