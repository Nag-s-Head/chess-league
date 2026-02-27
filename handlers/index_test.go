package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	testutils "github.com/Nag-s-Head/chess-league/db/test_utils"
	"github.com/Nag-s-Head/chess-league/handlers"
	"github.com/stretchr/testify/require"
)

func TestIndex(t *testing.T) {
	db := testutils.GetDb(t)
	defer db.Close()

	r := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	w := httptest.NewRecorder()
	handlers.Index(db)(w, r)

	body, err := io.ReadAll(w.Result().Body)
	require.NoError(t, err)
	require.NotEmpty(t, body)
}
