package model_test

import (
	"testing"

	"github.com/Nag-s-Head/chess-league/db/model"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/djpiper28/rpg-book/common/normalisation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestDeletePlayer(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	// 1. Create a player
	name := uuid.New().String()
	player := model.NewPlayer(name)
	require.NoError(t, model.InsertPlayer(db, player))

	// 2. Create an admin user
	admin := model.NewAdminUser("admin", uuid.New().String(), "127.0.0.1", "test-agent")
	tx, err := db.GetSqlxDb().BeginTxx(t.Context(), nil)
	require.NoError(t, err)
	require.NoError(t, model.InsertAdminUser(tx, *admin))
	require.NoError(t, tx.Commit())

	// 3. Delete the player
	err = model.DeletePlayer(db, player.Id, admin.Id)
	require.NoError(t, err)

	// 4. Verify player is deleted and fields updated
	deletedPlayer, err := model.GetPlayer(db, player.Id)
	require.NoError(t, err)
	require.True(t, deletedPlayer.Deleted)
	require.Equal(t, player.Id.String(), deletedPlayer.Name)
	require.Equal(t, normalisation.Normalise(player.Id.String()), deletedPlayer.NameNormalised)

	// 5. Try to delete again - should fail
	err = model.DeletePlayer(db, player.Id, admin.Id)
	require.Error(t, err)
	require.Contains(t, err.Error(), "already been deleted")
}

func TestDeleteNonExistentPlayer(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	adminId := uuid.New()
	err := model.DeletePlayer(db, uuid.New(), adminId)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Cannot get player")
}
