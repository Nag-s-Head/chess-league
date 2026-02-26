package submitgame_test

import (
	"testing"

	submitgame "github.com/Nag-s-Head/chess-league/handlers/submit_game"
	"github.com/stretchr/testify/require"
)

func TestSubmit(t *testing.T) {
	tpl, err := submitgame.Render()
	require.NoError(t, err)
	require.NotNil(t, tpl)
}
