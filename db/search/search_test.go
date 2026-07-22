package search_test

import (
	"log/slog"
	"strings"
	"testing"

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

	// Search uses log debug for queries, this makes it far easier to reason with test failures
	slog.SetLogLoggerLevel(slog.LevelDebug)

	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Dave"+uuid.NewString())))
	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Chas"+uuid.NewString())))
	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Danny"+uuid.NewString())))
	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Greg"+uuid.NewString())))
	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Sophie"+uuid.NewString())))
	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Anne"+uuid.NewString())))
	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Charlotte"+uuid.NewString())))

	beryl := model.NewPlayer("Beryl" + uuid.NewString())
	beryl.Deleted = true
	beryl.Liglicko2Rating = 1900
	require.NoError(t, model.InsertPlayer(db, beryl))

	t.Run("Default Search", func(t *testing.T) {
		players, err := search.SearchPlayers(db, "dave")
		require.NoError(t, err)
		require.NotEmpty(t, players)
		require.Contains(t, players[0].Name, "Dave")
	})

	t.Run("Default Search Fuzzy", func(t *testing.T) {
		players, err := search.SearchPlayers(db, "dav")
		require.NoError(t, err)
		require.NotEmpty(t, players)
		require.Contains(t, players[0].Name, "Dave")
	})

	t.Run("Test Rating Alias", func(t *testing.T) {
		players, err := search.SearchPlayers(db, "rating>=1900")
		require.NoError(t, err)
		require.NotEmpty(t, players)
		require.Contains(t, players[0].Name, "Beryl")
	})

	t.Run("Test Deleted Column", func(t *testing.T) {
		players, err := search.SearchPlayers(db, "deleted=true")
		require.NoError(t, err)
		require.NotEmpty(t, players)
		require.Contains(t, players[0].Name, "Beryl")
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
