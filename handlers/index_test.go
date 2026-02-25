package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Nag-s-Head/chess-league/handlers"
	"github.com/stretchr/testify/require"
)

func TestIndex(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	w := httptest.NewRecorder()
	handlers.Index(w, r)

	body, err := io.ReadAll(w.Result().Body)
	require.NoError(t, err)
	require.NotEmpty(t, body)
}
