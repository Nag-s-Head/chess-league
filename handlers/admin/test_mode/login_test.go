package testmode_test 

import (
	"testing"

	testmode "github.com/Nag-s-Head/chess-league/handlers/admin/test_mode"
	"github.com/stretchr/testify/require"
)

func TestRender(t *testing.T) {
	tpl, err := testmode.Login(nil, nil)
	require.NoError(t, err)
	require.NotNil(t, tpl)
}
