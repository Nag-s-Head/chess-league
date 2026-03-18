package auth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/stretchr/testify/require"
)

func TestLoginHandler(t *testing.T) {
	t.Run("Test Mode", func(t *testing.T) {
		isTestMode = true
		defer func() { isTestMode = os.Getenv("TEST_MODE") == "true" }()

		req := httptest.NewRequest(http.MethodGet, "/admin/login", nil)
		rr := httptest.NewRecorder()

		Login(rr, req)

		require.Equal(t, http.StatusTemporaryRedirect, rr.Code)
		require.Equal(t, "/admin/test-mode", rr.Header().Get("Location"))
	})

	t.Run("Production Mode", func(t *testing.T) {
		isTestMode = false
		defer func() { isTestMode = os.Getenv("TEST_MODE") == "true" }()

		req := httptest.NewRequest(http.MethodGet, "/admin/login", nil)
		rr := httptest.NewRecorder()

		Login(rr, req)

		require.Equal(t, http.StatusTemporaryRedirect, rr.Code)
		require.Contains(t, rr.Header().Get("Location"), "github.com/login/oauth/authorize")
	})
}

func TestCallbackHandler_NoCode(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	req := httptest.NewRequest(http.MethodGet, "/admin/auth/callback", nil)
	rr := httptest.NewRecorder()

	handler := Callback(db)
	handler(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}
