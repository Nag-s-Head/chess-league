package model_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nag-s-Head/chess-league/db/model"
	"github.com/stretchr/testify/require"
)

func TestGetRemoteAddr(t *testing.T) {
	t.Run("X-Forwarded-For present", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("X-Forwarded-For", "1.2.3.4")
		r.RemoteAddr = "127.0.0.1:12345"

		require.Equal(t, "1.2.3.4", model.GetRemoteAddr(r))
	})

	t.Run("X-Forwarded-For case insensitivity", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("x-forwarded-for", "1.2.3.4")
		r.RemoteAddr = "127.0.0.1:12345"

		require.Equal(t, "1.2.3.4", model.GetRemoteAddr(r))
	})

	t.Run("X-Forwarded-For list", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("X-Forwarded-For", "1.2.3.4, 5.6.7.8")
		r.RemoteAddr = "127.0.0.1:12345"

		require.Equal(t, "1.2.3.4", model.GetRemoteAddr(r))
	})

	t.Run("X-Forwarded-For absent, RemoteAddr with port", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.RemoteAddr = "127.0.0.1:12345"

		require.Equal(t, "127.0.0.1", model.GetRemoteAddr(r))
	})

	t.Run("X-Forwarded-For absent, RemoteAddr without port", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.RemoteAddr = "127.0.0.1"

		require.Equal(t, "127.0.0.1", model.GetRemoteAddr(r))
	})

	t.Run("X-Forwarded-For with spaces", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("X-Forwarded-For", "  1.2.3.4  , 5.6.7.8 ")

		require.Equal(t, "1.2.3.4", model.GetRemoteAddr(r))
	})
}
