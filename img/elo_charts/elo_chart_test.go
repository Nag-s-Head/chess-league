package elo_charts_test

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"testing"

	elo_charts "github.com/Nag-s-Head/chess-league/img/elo_charts"
	"github.com/stretchr/testify/require"
)

func doTest(t *testing.T, params elo_charts.Params) {
	t.Helper()

	img, err := elo_charts.Render(params)
	require.NoError(t, err)

	_, _, err = image.Decode(bytes.NewBuffer(img))
	require.NoError(t, err)

	filename := fmt.Sprintf("%s_test.png", t.Name())
	err = os.WriteFile(filename, img, 0666)
	require.NoError(t, err)

	t.Logf("Saved as %s", filename)
}

func TestRenderEloNormal(t *testing.T) {
	t.Parallel()

	doTest(t, elo_charts.Params{
		EndElo: 1500,
		Changes: []elo_charts.EloChange{
			{
				Delta: 20,
			},
			{
				Delta: 15,
			},
			{
				Delta: 0,
			},
			{
				Delta: -5,
			},
			{
				Delta: -20,
			},
			{
				Delta: 20,
			},
		},
	})
}

func TestOneEntry(t *testing.T) {
	t.Parallel()

	doTest(t, elo_charts.Params{
		EndElo: 1500,
		Changes: []elo_charts.EloChange{
			{
				Delta: 20,
			},
		},
	})
}

func TestNoChanges(t *testing.T) {
	t.Parallel()

	doTest(t, elo_charts.Params{
		EndElo:  1500,
		Changes: []elo_charts.EloChange{},
	})
}
