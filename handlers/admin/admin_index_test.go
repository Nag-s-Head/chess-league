package admin_test

import (
	"testing"

	"github.com/Nag-s-Head/chess-league/db/model"
	"github.com/Nag-s-Head/chess-league/handlers/admin"
	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	tpl, err := admin.AdminIndex(nil, nil, model.NewAdminUser("bob", "bob", "0.0.0.0", "bob"))
	require.NoError(t, err)
	require.NotNil(t, tpl)
}
