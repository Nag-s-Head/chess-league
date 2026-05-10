package model_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func setupAdmin(t *testing.T, database *db.Db) model.AdminUser {
	admin := model.NewAdminUser("admin", uuid.New().String(), "password", "salt")
	tx, err := database.GetSqlxDb().BeginTxx(context.Background(), nil)
	require.NoError(t, err)
	defer tx.Rollback()

	err = model.InsertAdminUser(tx, *admin)
	require.NoError(t, err)
	require.NoError(t, tx.Commit())
	return *admin
}

func createTestGame(t *testing.T, database *db.Db, white, black *model.Player, score model.Score) int64 {
	tx, err := database.GetSqlxDb().BeginTxx(context.Background(), nil)
	require.NoError(t, err)
	defer tx.Rollback()

	ikey, err := model.NextIKey(database)
	require.NoError(t, err)

	req, _ := http.NewRequest("GET", "/", nil)
	req.RemoteAddr = "127.0.0.1:1234"

	// CreateGame updates the player objects passed to it
	_, _, _, err = model.CreateGame(tx, white, black, true, ikey, score, req)
	require.NoError(t, err)
	require.NoError(t, tx.Commit())
	return ikey
}

func TestMergePlayers_NoGames(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()
	admin := setupAdmin(t, db)

	dest := model.NewPlayer("John Smith " + uuid.New().String())
	target := model.NewPlayer("Jon Smith " + uuid.New().String())

	require.NoError(t, model.InsertPlayer(db, dest))
	require.NoError(t, model.InsertPlayer(db, target))

	// Should probably not fail if target has no games, but current implementation might
	err := model.MergePlayers(db, target.Id, dest.Id, admin.Id)
	require.NoError(t, err)

	// Check target is deleted
	p, err := model.GetPlayer(db, target.Id)
	require.NoError(t, err)
	require.True(t, p.Deleted)
	require.Equal(t, p.Id.String(), p.NameNormalised)
}

func TestMergePlayers_TargetPlayedDest(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()
	admin := setupAdmin(t, db)

	dest := model.NewPlayer("John Smith " + uuid.New().String())
	target := model.NewPlayer("Jon Smith " + uuid.New().String())

	require.NoError(t, model.InsertPlayer(db, dest))
	require.NoError(t, model.InsertPlayer(db, target))

	// Target vs Dest
	createTestGame(t, db, &target, &dest, model.Score_Win)

	err := model.MergePlayers(db, target.Id, dest.Id, admin.Id)
	require.NoError(t, err)

	// Check the game is deleted (self-play)
	games, err := model.GetGamesByPlayer(db, dest.Id)
	require.NoError(t, err)
	for _, g := range games {
		if g.PlayerWhite == dest.Id && g.PlayerBlack == dest.Id {
			require.True(t, g.Deleted, "Self-play game should be marked as deleted")
		}
	}
}

func TestMergePlayers_ManyGames(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()
	admin := setupAdmin(t, db)

	dest := model.NewPlayer("John Smith " + uuid.New().String())
	target := model.NewPlayer("Jon Smith " + uuid.New().String())
	other := model.NewPlayer("Other Player " + uuid.New().String())

	require.NoError(t, model.InsertPlayer(db, dest))
	require.NoError(t, model.InsertPlayer(db, target))
	require.NoError(t, model.InsertPlayer(db, other))

	// 10 games for dest
	for i := 0; i < 10; i++ {
		createTestGame(t, db, &dest, &other, model.Score_Win)
	}

	// 10 games for target
	for i := 0; i < 10; i++ {
		createTestGame(t, db, &target, &other, model.Score_Loss)
	}

	err := model.MergePlayers(db, target.Id, dest.Id, admin.Id)
	require.NoError(t, err)

	// Check total games for dest
	games, err := model.GetGamesByPlayer(db, dest.Id)
	require.NoError(t, err)

	// Should have 20 games (10 from dest, 10 from target)
	// Some might be deleted if we don't filter them in GetGamesByPlayer,
	// but GetGamesByPlayer doesn't seem to filter by deleted.

	activeGames := 0
	for _, g := range games {
		if !g.Deleted {
			activeGames++
		}
	}
	require.Equal(t, 20, activeGames)

	// Check target is deleted
	p, err := model.GetPlayer(db, target.Id)
	require.NoError(t, err)
	require.True(t, p.Deleted)
	require.Equal(t, p.Id.String(), p.NameNormalised)
}
