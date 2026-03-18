package qrcode_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nag-s-Head/chess-league/db/model"
	qrcode "github.com/Nag-s-Head/chess-league/handlers/admin/qr_code"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	user := model.NewAdminUser(uuid.New().String(), uuid.New().String(), "0.0.0.0", "bob")
	req := httptest.NewRequest(http.MethodGet, "/admin/qr-code", nil)
	rr := httptest.NewRecorder()

	qrcode.Render(user)(rr, req)
	require.NotEmpty(t, rr.Body.Bytes())
}
