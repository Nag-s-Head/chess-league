package search_test

import (
	"log/slog"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
	"github.com/Nag-s-Head/chess-league/db/search"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func setupGames(t *testing.T, db db.Db) {
	t.Helper()

	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Dariuz"+uuid.NewString())))
	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Chas"+uuid.NewString())))
	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Danny"+uuid.NewString())))
	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Greg"+uuid.NewString())))
	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Sophie"+uuid.NewString())))
	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Anne"+uuid.NewString())))
	require.NoError(t, model.InsertPlayer(db, model.NewPlayer("Charlotte"+uuid.NewString())))

	andy := model.NewPlayer("Andy Burnham" + uuid.NewString())
	require.NoError(t, model.InsertPlayer(db, andy))

	beryl := model.NewPlayer("Beryl" + uuid.NewString())
	beryl.Deleted = true
	beryl.Liglicko2Rating = 3900
	require.NoError(t, model.InsertPlayer(db, beryl))

	require.NoError(t, db.DoTx(func(tx *sqlx.Tx) error {
		ikey, err := model.NextIKey(db)
		require.NoError(t, err)

		r := httptest.NewRequest("GET", "/mocked-url", nil)
		_, _, _, err = model.CreateGame(tx, &andy, &beryl, true, ikey, model.Score_Win, r)
		require.NoError(t, err)
		return nil
	}))

	require.NoError(t, db.DoTx(func(tx *sqlx.Tx) error {
		ikey, err := model.NextIKey(db)
		require.NoError(t, err)

		r := httptest.NewRequest("GET", "/mocked-url", nil)
		_, _, _, err = model.CreateGame(tx, &andy, &beryl, false, ikey, model.Score_Draw, r)
		require.NoError(t, err)
		return nil
	}))

	ikey, err := model.NextIKey(db)
	require.NoError(t, err)
	require.NoError(t, db.DoTx(func(tx *sqlx.Tx) error {
		r := httptest.NewRequest("GET", "/mocked-url", nil)
		_, _, _, err = model.CreateGame(tx, &andy, &beryl, false, ikey, model.Score_Loss, r)
		require.NoError(t, err)

		return nil
	}))

	admin := model.NewAdminUser(uuid.NewString(), uuid.NewString(), "ip", "user agent")

	require.NoError(t, db.DoTx(func(tx *sqlx.Tx) error {
		require.NoError(t, model.InsertAdminUser(tx, *admin))
		return nil
	}))
	require.NoError(t, model.DeleteGame(db, admin.Id, ikey))
}

func TestGamesSearch(t *testing.T) {
	t.Parallel()

	db := testutils.GetDb(t)
	defer db.Close()

	setupGames(t, db)

	t.Run("Default Search", func(t *testing.T) {
		games, err := search.SearchGames(db, "andy")
		require.NoError(t, err)
		require.True(t, len(games) >= 1)

		games2, err := search.SearchGames(db, "beryl")
		require.NoError(t, err)
		require.True(t, len(games2) >= 1)

		require.Contains(t, games, games2[0])
	})

	t.Run("Any Player Search", func(t *testing.T) {
		games, err := search.SearchGames(db, "any_player:andy")
		require.NoError(t, err)
		require.True(t, len(games) >= 1)
	})

	t.Run("Deleted Search", func(t *testing.T) {
		games, err := search.SearchGames(db, "deleted=true")
		require.NoError(t, err)
		require.True(t, len(games) >= 1)
		require.True(t, games[0].Deleted)

		games, err = search.SearchGames(db, "deleted=false")
		require.NoError(t, err)
		require.True(t, len(games) >= 1)
		require.False(t, games[0].Deleted)
	})

	t.Run("White Player Search", func(t *testing.T) {
		games, err := search.SearchGames(db, "white_player:andy")
		require.NoError(t, err)
		require.True(t, len(games) >= 1)

		for _, game := range games {
			player, err := model.GetPlayer(db, game.PlayerWhite)
			require.NoError(t, err)
			require.Contains(t, player.NameNormalised, "andy")
		}
	})

	t.Run("Black Player Search", func(t *testing.T) {
		games, err := search.SearchGames(db, "black_player:andy")
		require.NoError(t, err)
		require.True(t, len(games) >= 1)

		for _, game := range games {
			player, err := model.GetPlayer(db, game.PlayerBlack)
			require.NoError(t, err)
			require.Contains(t, player.NameNormalised, "andy")
		}
	})

	t.Run("Score Search", func(t *testing.T) {
		games, err := search.SearchGames(db, "score=1-0")
		require.NoError(t, err)
		require.True(t, len(games) >= 1)

		games, err = search.SearchGames(db, "score=0-1")
		require.NoError(t, err)
		require.True(t, len(games) >= 1)

		games, err = search.SearchGames(db, "score=1/2-1/2")
		require.NoError(t, err)
		require.True(t, len(games) >= 1)
	})
}

func FuzzSearchGame(f *testing.F) {
	// Search uses log debug for queries, this makes it far easier to reason with test failures
	slog.SetLogLoggerLevel(slog.LevelDebug)

	f.Add("'; DELETE FROM games; --")

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

	f.Add("player_white=Andy Burnham")
	f.Add("player_white~Andy Burnham")
	f.Add(`player_white~"Andy Burnham"`)

	f.Add(`deleted=false and (any_player=greg OR any_player="chas") and (liglicko2_white>600 and liglicko2_black<1600)`)

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
		setupGames(t, dbInstance)

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
