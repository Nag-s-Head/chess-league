package adminutils_test

import (
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	adminutils "github.com/Nag-s-Head/chess-league/handlers/admin/admin_utils"
	"github.com/stretchr/testify/require"
)

func TestRenderError(t *testing.T) {
	rr := httptest.NewRecorder()

	testError := errors.New("TEST ERROR")
	adminutils.RenderError(rr, testError)

	res := rr.Result()
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.NotEmpty(t, body)
	require.Contains(t, string(body), fmt.Sprintf("%s", testError))
}
