package model_test

import (
	"testing"

	"github.com/Nag-s-Head/chess-league/db/model"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestSetAndGetPlayers(t *testing.T) {
	t.Parallel()

	db := testutils.GetDb(t)
	defer db.Close()

	for range 10 {
		err := model.InsertPlayer(db, model.NewPlayer(uuid.New().String()))
		require.NoError(t, err)
	}

	var targetPlayers []uuid.UUID
	for range 10 {
		player := model.NewPlayer(uuid.New().String() + "-league-player")
		err := model.InsertPlayer(db, player)
		require.NoError(t, err)
		targetPlayers = append(targetPlayers, player.Id)
	}

	admin := model.NewAdminUser(uuid.New().String(), uuid.New().String(), "test", "test")
	require.NoError(t, db.DoTx(func(tx *sqlx.Tx) error {
		require.NoError(t, model.InsertAdminUser(tx, *admin))
		return nil
	}))

	require.NoError(t, model.SetLeaguePlayers(db, admin.Id, targetPlayers))

	require.NoError(t, db.DoTx(func(tx *sqlx.Tx) error {
		leaguePlayers, err := model.GetLeaguePlayers(tx)
		require.NoError(t, err)

		for _, expectedId := range targetPlayers {
			found := false
			for _, player := range leaguePlayers {
				if player.Id == expectedId {
					found = true
					break
				}
			}

			require.True(t, found, "Cannot find the target player %s in GetLeaguePlayers", expectedId)
		}

		require.Len(t, leaguePlayers, len(targetPlayers))

		return nil
	}))

	players, err := model.GetUiFriendlyLeaguePlayers(db)
	require.NoError(t, err)

	for _, expectedId := range targetPlayers {
		found := false
		for _, player := range players {
			if player.Id == expectedId {
				found = true
				require.True(t, player.InLeague)
				break
			}
		}

		require.True(t, found, "Cannot find the target player %s in GetUiFriendlyLeaguePlayers", expectedId)
	}
}
