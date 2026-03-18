package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nag-s-Head/chess-league/db/model"
	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/Nag-s-Head/chess-league/handlers/admin/auth"
	"github.com/stretchr/testify/require"
)

func TestLogout(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	t.Run("No cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/logout", nil)
		rr := httptest.NewRecorder()

		handler := auth.Logout(db)
		handler(rr, req)

		require.Equal(t, http.StatusTemporaryRedirect, rr.Code)
		require.Equal(t, "/", rr.Header().Get("Location"))
	})

	t.Run("Valid cookie", func(t *testing.T) {
		// Setup user
		user, err := model.AdminLogin(db, "Logout Tester", "github-logout", "127.0.0.1", "UA")
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/admin/logout", nil)
		req.AddCookie(&http.Cookie{
			Name:  auth.AuthCookie,
			Value: user.SessionKey,
		})
		rr := httptest.NewRecorder()

		handler := auth.Logout(db)
		handler(rr, req)

		require.Equal(t, http.StatusTemporaryRedirect, rr.Code)
		require.Equal(t, "/", rr.Header().Get("Location"))

		// Verify user is logged out in DB
		dbUser, err := model.AdminLogin(db, user.Name, user.OauthId, user.LastIp, user.LastUserAgent)
		require.NoError(t, err)
		require.NotEqual(t, user.SessionKey, dbUser.SessionKey)

		// Verify the old session key is gone
		_, err = model.AdminGetFromSessionKey(db, user.SessionKey)
		require.Error(t, err)
	})

	t.Run("Invalid session key in cookie", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/logout", nil)
		req.AddCookie(&http.Cookie{
			Name:  auth.AuthCookie,
			Value: "invalid-key",
		})
		rr := httptest.NewRecorder()

		handler := auth.Logout(db)
		handler(rr, req)

		require.Equal(t, http.StatusTemporaryRedirect, rr.Code)
		require.Equal(t, "/", rr.Header().Get("Location"))
	})
}
