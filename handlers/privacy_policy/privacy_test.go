package privacypolicy_test

import (
	"testing"

	privacypolicy "github.com/Nag-s-Head/chess-league/handlers/privacy_policy"
	"github.com/stretchr/testify/require"
)

func TestPrivacy(t *testing.T) {
	tpl, err := privacypolicy.Render()
	require.NoError(t, err)
	require.NotNil(t, tpl)
}
