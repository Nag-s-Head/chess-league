package model_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Nag-s-Head/chess-league/db/model"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/djpiper28/rpg-book/common/normalisation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestDeletePlayer(t *testing.T) {
	t.Parallel()

	t.Run("No games deleted", func(t *testing.T) {
		db := testutils.GetDb(t)
		defer db.Close()

		name := uuid.New().String()
		player := model.NewPlayer(name)
		require.NoError(t, model.InsertPlayer(db, player))

		admin := model.NewAdminUser("admin", uuid.New().String(), "127.0.0.1", "test-agent")
		tx, err := db.GetSqlxDb().BeginTxx(t.Context(), nil)
		require.NoError(t, err)
		require.NoError(t, model.InsertAdminUser(tx, *admin))
		require.NoError(t, tx.Commit())

		err = model.DeletePlayer(db, player.Id, admin.Id)
		require.NoError(t, err)

		deletedPlayer, err := model.GetPlayer(db, player.Id)
		require.NoError(t, err)
		require.True(t, deletedPlayer.Deleted)
		require.Equal(t, player.Id.String(), deletedPlayer.Name)
		require.Equal(t, normalisation.Normalise(player.Id.String()), deletedPlayer.NameNormalised)

		err = model.DeletePlayer(db, player.Id, admin.Id)
		require.Error(t, err)
		require.Contains(t, err.Error(), "already been deleted")
	})

	t.Run("Lots of games deleted", func(t *testing.T) {
		db := testutils.GetDb(t)
		defer db.Close()

		admin := model.NewAdminUser("admin", uuid.New().String(), "127.0.0.1", "test-agent")
		tx, err := db.GetSqlxDb().BeginTxx(t.Context(), nil)
		require.NoError(t, err)
		require.NoError(t, model.InsertAdminUser(tx, *admin))
		require.NoError(t, tx.Commit())

		p1 := model.NewPlayer(uuid.New().String())
		require.NoError(t, model.InsertPlayer(db, p1))
		p2 := model.NewPlayer(uuid.New().String())
		require.NoError(t, model.InsertPlayer(db, p2))

		for i := 0; i < 10; i++ {
			tx, err := db.GetSqlxDb().BeginTxx(t.Context(), nil)
			require.NoError(t, err)
			ikey, err := model.NextIKey(db)
			require.NoError(t, err)
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			_, _, _, err = model.CreateGame(tx, &p1, &p2, true, ikey, model.Score_Win, r)
			require.NoError(t, err)
			require.NoError(t, tx.Commit())

			p1, _ = model.GetPlayer(db, p1.Id)
			p2, _ = model.GetPlayer(db, p2.Id)
		}

		err = model.DeletePlayer(db, p1.Id, admin.Id)
		require.NoError(t, err)

		dp1, err := model.GetPlayer(db, p1.Id)
		require.NoError(t, err)
		require.True(t, dp1.Deleted)

		var count int
		err = db.GetSqlxDb().Get(&count, "SELECT count(*) FROM games WHERE (player_white=$1 OR player_black=$1) AND deleted=FALSE", p1.Id)
		require.NoError(t, err)
		require.Equal(t, 0, count)

		err = db.GetSqlxDb().Get(&count, "SELECT count(*) FROM games WHERE (player_white=$1 OR player_black=$1) AND deleted=TRUE", p1.Id)
		require.NoError(t, err)
		require.Equal(t, 10, count)
	})

	t.Run("Game replay required", func(t *testing.T) {
		db := testutils.GetDb(t)
		defer db.Close()

		admin := model.NewAdminUser("admin", uuid.New().String(), "127.0.0.1", "test-agent")
		tx, err := db.GetSqlxDb().BeginTxx(t.Context(), nil)
		defer tx.Rollback()

		require.NoError(t, err)
		require.NoError(t, model.InsertAdminUser(tx, *admin))
		require.NoError(t, tx.Commit())

		playerA := model.NewPlayer(uuid.New().String())
		require.NoError(t, model.InsertPlayer(db, playerA))
		deletedPlayer := model.NewPlayer(uuid.New().String())
		require.NoError(t, model.InsertPlayer(db, deletedPlayer))
		playerB := model.NewPlayer(uuid.New().String())
		require.NoError(t, model.InsertPlayer(db, playerB))

		createGame := func(p1, p2 *model.Player, score model.Score) model.Game {
			tx2, err := db.GetSqlxDb().BeginTxx(context.Background(), nil)
			require.NoError(t, err)

			defer tx2.Rollback()

			ikey, err := model.NextIKey(db)
			require.NoError(t, err)

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			g, _, _, err := model.CreateGame(tx2, p1, p2, true, ikey, score, r)
			require.NoError(t, err)
			require.NoError(t, tx2.Commit())

			*p1, _ = model.GetPlayer(db, p1.Id)
			*p2, _ = model.GetPlayer(db, p2.Id)
			return g
		}

		playerC := model.NewPlayer(uuid.New().String())
		require.NoError(t, model.InsertPlayer(db, playerC))
		playerD := model.NewPlayer(uuid.New().String())
		require.NoError(t, model.InsertPlayer(db, playerD))
		createGame(&playerC, &playerD, model.Score_Draw)

		g1 := createGame(&playerA, &deletedPlayer, model.Score_Win)
		time.Sleep(time.Millisecond)
		g2 := createGame(&deletedPlayer, &playerB, model.Score_Win)
		time.Sleep(time.Millisecond)
		g3 := createGame(&playerA, &playerB, model.Score_Win)

		ratingA_before := playerA.Liglicko2Rating

		err = model.DeletePlayer(db, deletedPlayer.Id, admin.Id)
		require.NoError(t, err)

		dpb, err := model.GetPlayer(db, deletedPlayer.Id)
		require.NoError(t, err)
		require.True(t, dpb.Deleted)

		g1After, err := model.GetGameWithDetails(db, g1.IKey)
		require.NoError(t, err)
		require.True(t, g1After.Deleted)

		g2After, err := model.GetGameWithDetails(db, g2.IKey)
		require.NoError(t, err)
		require.True(t, g2After.Deleted)

		g3After, err := model.GetGameWithDetails(db, g3.IKey)
		require.NoError(t, err)
		require.False(t, g3After.Deleted)

		pa_after, err := model.GetPlayer(db, playerA.Id)
		require.NoError(t, err)
		require.NotEqual(t, ratingA_before, pa_after.Liglicko2Rating)
	})
}

func TestDeleteNonExistentPlayer(t *testing.T) {
	t.Parallel()
	db := testutils.GetDb(t)
	defer db.Close()

	adminId := uuid.New()
	err := model.DeletePlayer(db, uuid.New(), adminId)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Cannot get player")
}
