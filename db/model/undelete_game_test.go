package model_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Nag-s-Head/chess-league/db"
	"github.com/Nag-s-Head/chess-league/db/model"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type undeleteFixture struct {
	db      db.Db
	admin   *model.AdminUser
	white   model.Player
	black   model.Player
	request *http.Request
}

func newUndeleteFixture(t *testing.T) *undeleteFixture {
	t.Helper()

	database := testutils.GetDb(t)
	t.Cleanup(database.Close)

	admin := model.NewAdminUser("Undelete test admin", uuid.NewString(), "127.0.0.1", "Go test")
	tx, err := database.GetSqlxDb().BeginTxx(t.Context(), nil)
	require.NoError(t, err)
	require.NoError(t, model.InsertAdminUser(tx, *admin))
	require.NoError(t, tx.Commit())

	white := model.NewPlayer("undelete-white-" + uuid.NewString())
	require.NoError(t, model.InsertPlayer(database, white))
	black := model.NewPlayer("undelete-black-" + uuid.NewString())
	require.NoError(t, model.InsertPlayer(database, black))

	request := httptest.NewRequest(http.MethodPost, "/", nil)
	request.RemoteAddr = "127.0.0.1"

	return &undeleteFixture{
		db:      database,
		admin:   admin,
		white:   white,
		black:   black,
		request: request,
	}
}

func (f *undeleteFixture) createGame(t *testing.T, score model.Score) int64 {
	t.Helper()

	ikey, err := model.NextIKey(f.db)
	require.NoError(t, err)
	tx, err := f.db.GetSqlxDb().BeginTxx(t.Context(), nil)
	require.NoError(t, err)
	_, _, _, err = model.CreateGame(tx, &f.white, &f.black, true, ikey, score, f.request)
	require.NoError(t, err)
	require.NoError(t, tx.Commit())
	return ikey
}

func requireSameRating(t *testing.T, expected, actual model.Player) {
	t.Helper()
	require.InDelta(t, expected.Liglicko2Rating, actual.Liglicko2Rating, 1e-9)
	require.InDelta(t, expected.Liglicko2Deviation, actual.Liglicko2Deviation, 1e-9)
	require.InDelta(t, expected.Liglicko2Volatility, actual.Liglicko2Volatility, 1e-9)
	require.InDelta(t, expected.Liglicko2At, actual.Liglicko2At, 1e-9)
}

func requireSameGameRatings(t *testing.T, expected, actual model.GameWithDetails) {
	t.Helper()
	require.Equal(t, expected.IKey, actual.IKey)
	require.Equal(t, expected.Score, actual.Score)
	require.Equal(t, expected.Deleted, actual.Deleted)
	require.InDelta(t, expected.Liglicko2White, actual.Liglicko2White, 1e-9)
	require.InDelta(t, expected.Liglicko2Black, actual.Liglicko2Black, 1e-9)
	require.InDelta(t, expected.Liglicko2WhiteOldRating, actual.Liglicko2WhiteOldRating, 1e-9)
	require.InDelta(t, expected.Liglicko2WhiteOldDeviation, actual.Liglicko2WhiteOldDeviation, 1e-9)
	require.InDelta(t, expected.Liglicko2WhiteOldVolatility, actual.Liglicko2WhiteOldVolatility, 1e-9)
	require.InDelta(t, expected.Liglicko2WhiteOldAt, actual.Liglicko2WhiteOldAt, 1e-9)
	require.InDelta(t, expected.Liglicko2BlackOldRating, actual.Liglicko2BlackOldRating, 1e-9)
	require.InDelta(t, expected.Liglicko2BlackOldDeviation, actual.Liglicko2BlackOldDeviation, 1e-9)
	require.InDelta(t, expected.Liglicko2BlackOldVolatility, actual.Liglicko2BlackOldVolatility, 1e-9)
	require.InDelta(t, expected.Liglicko2BlackOldAt, actual.Liglicko2BlackOldAt, 1e-9)
}

func TestUnDeleteGameReplaysRatingsAndAuditsAction(t *testing.T) {
	f := newUndeleteFixture(t)

	ikeys := []int64{
		f.createGame(t, model.Score_Win),
		f.createGame(t, model.Score_Loss),
		f.createGame(t, model.Score_Draw),
	}

	wantWhite, err := model.GetPlayer(f.db, f.white.Id)
	require.NoError(t, err)
	wantBlack, err := model.GetPlayer(f.db, f.black.Id)
	require.NoError(t, err)
	wantGames := make([]model.GameWithDetails, len(ikeys))
	for i, ikey := range ikeys {
		wantGames[i], err = model.GetGameWithDetails(f.db, ikey)
		require.NoError(t, err)
	}

	undeletedIKey := ikeys[1]
	require.NoError(t, model.DeleteGame(f.db, f.admin.Id, undeletedIKey))
	deletedGame, err := model.GetGameWithDetails(f.db, undeletedIKey)
	require.NoError(t, err)
	require.True(t, deletedGame.Deleted)

	require.NoError(t, model.UnDeleteGame(f.db, f.admin.Id, undeletedIKey))

	for i, ikey := range ikeys {
		got, err := model.GetGameWithDetails(f.db, ikey)
		require.NoError(t, err)
		requireSameGameRatings(t, wantGames[i], got)
	}
	gotWhite, err := model.GetPlayer(f.db, f.white.Id)
	require.NoError(t, err)
	requireSameRating(t, wantWhite, gotWhite)
	gotBlack, err := model.GetPlayer(f.db, f.black.Id)
	require.NoError(t, err)
	requireSameRating(t, wantBlack, gotBlack)

	var descriptions []string
	err = f.db.GetSqlxDb().Select(&descriptions,
		"SELECT operation_description FROM audit_logs WHERE done_by = $1", f.admin.Id)
	require.NoError(t, err)
	ikeyString := strconv.FormatInt(undeletedIKey, 10)
	require.ElementsMatch(t, []string{
		"Deleted game " + ikeyString,
		"Undeleted game " + ikeyString,
	}, descriptions)
}

func TestUnDeleteLatestGameRestoresCurrentRatings(t *testing.T) {
	f := newUndeleteFixture(t)
	f.createGame(t, model.Score_Win)
	latestIKey := f.createGame(t, model.Score_Loss)

	wantWhite, err := model.GetPlayer(f.db, f.white.Id)
	require.NoError(t, err)
	wantBlack, err := model.GetPlayer(f.db, f.black.Id)
	require.NoError(t, err)
	wantGame, err := model.GetGameWithDetails(f.db, latestIKey)
	require.NoError(t, err)

	require.NoError(t, model.DeleteGame(f.db, f.admin.Id, latestIKey))
	require.NoError(t, model.UnDeleteGame(f.db, f.admin.Id, latestIKey))

	gotWhite, err := model.GetPlayer(f.db, f.white.Id)
	require.NoError(t, err)
	requireSameRating(t, wantWhite, gotWhite)
	gotBlack, err := model.GetPlayer(f.db, f.black.Id)
	require.NoError(t, err)
	requireSameRating(t, wantBlack, gotBlack)
	gotGame, err := model.GetGameWithDetails(f.db, latestIKey)
	require.NoError(t, err)
	requireSameGameRatings(t, wantGame, gotGame)
}

func TestUnDeleteGameReplaysAllSubsequentGames(t *testing.T) {
	f := newUndeleteFixture(t)

	var ikeys []int64
	for i := 0; i < 10; i++ {
		score := model.Score_Win
		if i%2 == 1 {
			score = model.Score_Loss
		}
		ikeys = append(ikeys, f.createGame(t, score))
	}

	wantGames := make([]model.GameWithDetails, len(ikeys))
	for i, ikey := range ikeys {
		game, err := model.GetGameWithDetails(f.db, ikey)
		require.NoError(t, err)
		wantGames[i] = game
	}
	wantWhite, err := model.GetPlayer(f.db, f.white.Id)
	require.NoError(t, err)
	wantBlack, err := model.GetPlayer(f.db, f.black.Id)
	require.NoError(t, err)

	deletedIKey := ikeys[2]
	require.NoError(t, model.DeleteGame(f.db, f.admin.Id, deletedIKey))
	require.NoError(t, model.UnDeleteGame(f.db, f.admin.Id, deletedIKey))

	for i, ikey := range ikeys {
		got, err := model.GetGameWithDetails(f.db, ikey)
		require.NoError(t, err)
		requireSameGameRatings(t, wantGames[i], got)
	}
	gotWhite, err := model.GetPlayer(f.db, f.white.Id)
	require.NoError(t, err)
	requireSameRating(t, wantWhite, gotWhite)
	gotBlack, err := model.GetPlayer(f.db, f.black.Id)
	require.NoError(t, err)
	requireSameRating(t, wantBlack, gotBlack)
}

func TestUnDeleteNonExistentGameFails(t *testing.T) {
	f := newUndeleteFixture(t)

	err := model.UnDeleteGame(f.db, f.admin.Id, 999999999)
	require.Error(t, err)
}

func TestUnDeleteGameThatIsNotDeletedFails(t *testing.T) {
	f := newUndeleteFixture(t)
	ikey := f.createGame(t, model.Score_Win)

	err := model.UnDeleteGame(f.db, f.admin.Id, ikey)
	require.Error(t, err)
	game, err := model.GetGameWithDetails(f.db, ikey)
	require.NoError(t, err)
	require.False(t, game.Deleted)
}

func TestUnDeleteGameTwiceFails(t *testing.T) {
	f := newUndeleteFixture(t)
	ikey := f.createGame(t, model.Score_Win)

	require.NoError(t, model.DeleteGame(f.db, f.admin.Id, ikey))
	require.NoError(t, model.UnDeleteGame(f.db, f.admin.Id, ikey))
	err := model.UnDeleteGame(f.db, f.admin.Id, ikey)
	require.Error(t, err)

	game, err := model.GetGameWithDetails(f.db, ikey)
	require.NoError(t, err)
	require.False(t, game.Deleted)
}
