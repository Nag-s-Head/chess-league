package rules_test

import (
	"testing"

	"github.com/Nag-s-Head/chess-league/handlers/rules"
	"github.com/stretchr/testify/require"
)

func TestPrivacy(t *testing.T) {
	tpl, err := rules.Render(false)
	require.NoError(t, err)
	require.NotNil(t, tpl)
}
