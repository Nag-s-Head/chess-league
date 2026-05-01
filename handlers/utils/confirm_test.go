package utils_test

import (
	"fmt"
	"io"
	"net/http/httptest"
	"net/url"
	"testing"

	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/Nag-s-Head/chess-league/handlers/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestConfirm(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	action := uuid.New().String()
	buttonValue := uuid.New().String()

	w := httptest.NewRecorder()
	utils.RenderConfirmationPage(w, action, buttonValue)

	body, err := io.ReadAll(w.Result().Body)
	require.NoError(t, err)
	require.NotEmpty(t, body)
	require.Contains(t, string(body), action)
	require.Contains(t, string(body), fmt.Sprintf(`value="%s"`, buttonValue))
}

func TestIsConfirmed(t *testing.T) {
	t.Run("Confirmed", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/", nil)
		r.Form = make(url.Values)
		r.Form.Set("confirm", "confirmed")
		require.True(t, utils.IsConfirmed(r))
	})

	t.Run("Not confirmed", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/", nil)
		r.Form = make(url.Values)
		r.Form.Set("confirm", "something else")
		require.False(t, utils.IsConfirmed(r))
	})

	t.Run("Missing", func(t *testing.T) {
		r := httptest.NewRequest("POST", "/", nil)
		r.Form = make(url.Values)
		require.False(t, utils.IsConfirmed(r))
	})
}
