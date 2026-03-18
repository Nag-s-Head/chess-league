package admin_test

import (
	"testing"

	"github.com/Nag-s-Head/chess-league/handlers/admin"
	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	tpl, err := admin.AdminIndex(nil, nil)
	require.NoError(t, err)
	require.NotNil(t, tpl)
}
