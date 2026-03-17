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

func TestForm(t *testing.T) {
	tpl, err := submitgame.Render()
	require.NoError(t, err)
	require.NotNil(t, tpl)
}

func TestSubmit(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	name := uuid.New().String()
	err := model.InsertPlayer(db, model.NewPlayer(name))
	require.NoError(t, err)

	t.Run("No magic number should error", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/mocked-url", strings.NewReader(""))
		w := httptest.NewRecorder()
		err = submitgame.DoSubmit(db, w, r)
		require.Error(t, err)
	})

	t.Run("No rules agreement should error", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/mocked-url", strings.NewReader(""))
		r.AddCookie(&http.Cookie{
			Name:  submitgame.MagicNumberCookie,
			Value: os.Getenv(submitgame.MagicNumberEnvVar),
		})
		w := httptest.NewRecorder()
		err = submitgame.DoSubmit(db, w, r)
		require.Error(t, err)
		require.Contains(t, err.Error(), "You must agree to the rules")
	})

	t.Run("No form data should error", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/mocked-url", strings.NewReader(""))
		r.AddCookie(&http.Cookie{
			Name:  submitgame.MagicNumberCookie,
			Value: os.Getenv(submitgame.MagicNumberEnvVar),
		})
		r.AddCookie(&http.Cookie{
			Name:  rules.RulesVersionCookie,
			Value: rules.CurrentRulesVersion,
		})
		r.AddCookie(&http.Cookie{
			Name:  rules.RulesVersionCookie,
			Value: rules.CurrentRulesVersion,
		})
		w := httptest.NewRecorder()
		err = submitgame.DoSubmit(db, w, r)
		require.Error(t, err)
	})

	t.Run("Test render of player lookup", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/mocked-url?player-name=%s&played-as=white&other-player-name=not_found&winner=white", name), strings.NewReader(""))
		r.AddCookie(&http.Cookie{
			Name:  submitgame.MagicNumberCookie,
			Value: os.Getenv(submitgame.MagicNumberEnvVar),
		})
		r.AddCookie(&http.Cookie{
			Name:  rules.RulesVersionCookie,
			Value: rules.CurrentRulesVersion,
		})
		r.AddCookie(&http.Cookie{
			Name:  rules.RulesVersionCookie,
			Value: rules.CurrentRulesVersion,
		})

		w := httptest.NewRecorder()
		err = submitgame.DoSubmit(db, w, r)
		require.NoError(t, err)
	})

	t.Run("Test render of player lookup, no magic number", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/mocked-url?player-name=%s&played-as=white&other-player-name=not_found&winner=white", name), strings.NewReader(""))

		w := httptest.NewRecorder()
		err = submitgame.DoSubmit(db, w, r)
		require.Error(t, err)
	})

	t.Run("Test final submit success", func(t *testing.T) {
		whiteName := uuid.New().String()
		blackName := uuid.New().String()

		form := url.Values{}
		form.Set("white-player-name", whiteName)
		form.Set("black-player-name", blackName)
		form.Set("winner", "white")
		form.Set("submit-type", "final")

		ikey, err := model.NextIKey(db)
		require.NoError(t, err)

		r := httptest.NewRequest(http.MethodPost, "/mocked-url", strings.NewReader(form.Encode()))
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
		require.Contains(t, w.Body.String(), "Game submitted successfully")
	})

	t.Run("Test final submit draw", func(t *testing.T) {
		whiteName := uuid.New().String()
		blackName := uuid.New().String()

		form := url.Values{}
		form.Set("white-player-name", whiteName)
		form.Set("black-player-name", blackName)
		form.Set("winner", "draw")
		form.Set("submit-type", "final")

		ikey, err := model.NextIKey(db)
		require.NoError(t, err)

		r := httptest.NewRequest(http.MethodPost, "/mocked-url", strings.NewReader(form.Encode()))
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
		require.Contains(t, w.Body.String(), "Game submitted successfully")
		require.Contains(t, w.Body.String(), "+0")
	})

	t.Run("Test final submit missing ikey", func(t *testing.T) {
		whiteName := uuid.New().String()
		blackName := uuid.New().String()

		form := url.Values{}
		form.Set("white-player-name", whiteName)
		form.Set("black-player-name", blackName)
		form.Set("winner", "white")
		form.Set("submit-type", "final")

		r := httptest.NewRequest(http.MethodPost, "/mocked-url", strings.NewReader(form.Encode()))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		r.AddCookie(&http.Cookie{
			Name:  submitgame.MagicNumberCookie,
			Value: os.Getenv(submitgame.MagicNumberEnvVar),
		})
		r.AddCookie(&http.Cookie{
			Name:  rules.RulesVersionCookie,
			Value: rules.CurrentRulesVersion,
		})

		w := httptest.NewRecorder()
		err = submitgame.DoSubmit(db, w, r)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Could not find ikey cookie")
	})

	t.Run("Test final submit invalid ikey", func(t *testing.T) {
		whiteName := uuid.New().String()
		blackName := uuid.New().String()

		form := url.Values{}
		form.Set("white-player-name", whiteName)
		form.Set("black-player-name", blackName)
		form.Set("winner", "white")
		form.Set("submit-type", "final")

		r := httptest.NewRequest(http.MethodPost, "/mocked-url", strings.NewReader(form.Encode()))
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
			Value: "not-a-number",
		})

		w := httptest.NewRecorder()
		err = submitgame.DoSubmit(db, w, r)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Could not read ikey cookie")
	})

	t.Run("Test final submit invalid winner", func(t *testing.T) {
		whiteName := uuid.New().String()
		blackName := uuid.New().String()

		form := url.Values{}
		form.Set("white-player-name", whiteName)
		form.Set("black-player-name", blackName)
		form.Set("winner", "invalid")
		form.Set("submit-type", "final")

		ikey, err := model.NextIKey(db)
		require.NoError(t, err)

		r := httptest.NewRequest(http.MethodPost, "/mocked-url", strings.NewReader(form.Encode()))
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
			Name:  rules.RulesVersionCookie,
			Value: rules.CurrentRulesVersion,
		})
		r.AddCookie(&http.Cookie{
			Name:  submitgame.IKeyCookie,
			Value: fmt.Sprintf("%d", ikey),
		})

		w := httptest.NewRecorder()
		err = submitgame.DoSubmit(db, w, r)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Invalid winner")
	})

	t.Run("Test final submit same players", func(t *testing.T) {
		whiteName := uuid.New().String()

		form := url.Values{}
		form.Set("white-player-name", whiteName)
		form.Set("black-player-name", whiteName)
		form.Set("winner", "white")
		form.Set("submit-type", "final")

		ikey, err := model.NextIKey(db)
		require.NoError(t, err)

		r := httptest.NewRequest(http.MethodPost, "/mocked-url", strings.NewReader(form.Encode()))
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
			Name:  rules.RulesVersionCookie,
			Value: rules.CurrentRulesVersion,
		})
		r.AddCookie(&http.Cookie{
			Name:  submitgame.IKeyCookie,
			Value: fmt.Sprintf("%d", ikey),
		})

		w := httptest.NewRecorder()
		err = submitgame.DoSubmit(db, w, r)
		require.Error(t, err)
		require.Contains(t, err.Error(), "Both players are the same")
	})

	t.Run("Empty magic cookie but valid URL param should work", func(t *testing.T) {
		url := fmt.Sprintf("/mocked-url?%s=%s&player-name=%s&played-as=white&other-player-name=not_found&winner=white",
			submitgame.MagicNumberParam, os.Getenv(submitgame.MagicNumberEnvVar), name)
		r := httptest.NewRequest(http.MethodGet, url, strings.NewReader(""))
		r.AddCookie(&http.Cookie{
			Name:  submitgame.MagicNumberCookie,
			Value: "",
		})
		r.AddCookie(&http.Cookie{
			Name:  rules.RulesVersionCookie,
			Value: rules.CurrentRulesVersion,
		})

		w := httptest.NewRecorder()
		err = submitgame.DoSubmit(db, w, r)
		require.NoError(t, err)
	})
}
