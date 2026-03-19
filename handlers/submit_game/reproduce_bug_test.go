package submitgame_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/Nag-s-Head/chess-league/db/model"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/Nag-s-Head/chess-league/handlers/rules"
	submitgame "github.com/Nag-s-Head/chess-league/handlers/submit_game"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSubmitterBug(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	whiteName := "WhitePlayer-" + uuid.New().String()
	blackName := "BlackPlayer-" + uuid.New().String()

	err := model.InsertPlayer(db, model.NewPlayer(whiteName))
	require.NoError(t, err)
	err = model.InsertPlayer(db, model.NewPlayer(blackName))
	require.NoError(t, err)

	players, err := model.SearchPlayerByName(db, whiteName)
	require.NoError(t, err)
	require.Len(t, players, 1)
	whitePlayer := players[0]

	players, err = model.SearchPlayerByName(db, blackName)
	require.NoError(t, err)
	require.Len(t, players, 1)
	blackPlayer := players[0]

	// Scenario: Black player submits the game
	form := url.Values{}
	// These are from the first step
	form.Set("player-name", blackName)
	form.Set("played-as", "black")
	form.Set("other-player-name", whiteName)
	form.Set("winner", "black")

	// These are from the second step (confirmation)
	form.Set("white-player-name", whiteName)
	form.Set("black-player-name", blackName)
	form.Set("submit-type", "final")

	ikey, err := model.NextIKey(db)
	require.NoError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/submit-game/submit", strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.AddCookie(&http.Cookie{
		Name:  submitgame.MagicNumberCookie,
		Value: os.Getenv(submitgame.MagicNumberEnvVar),
	})
	r.AddCookie(&http.Cookie{
		Name:  rules.RulesVersionCookie,
		Value: rules.CurrentRulesVersion,
	})
	r.AddCookie(&http.Cookie{
		Name:  submitgame.IKeyCookie,
		Value: fmt.Sprintf("%d", ikey),
	})

	w := httptest.NewRecorder()
	err = submitgame.DoSubmit(db, w, r)
	require.NoError(t, err)

	// Check the database for the game
	var games []model.Game
	err = db.GetSqlxDb().Select(&games, "SELECT * FROM games WHERE ikey=$1", ikey)
	require.NoError(t, err)
	require.Len(t, games, 1)

	game := games[0]
	t.Logf("Game submitter: %s", game.Submitter)
	t.Logf("White player: %s", whitePlayer.Id)
	t.Logf("Black player: %s", blackPlayer.Id)

	// This is what should happen: Submitter should be Black player
	require.Equal(t, blackPlayer.Id.String(), game.Submitter.String(), "Submitter should be the black player")
}
