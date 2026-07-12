package chess_league_test

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	chess_league "github.com/Nag-s-Head/chess-league/app"
	"github.com/Nag-s-Head/chess-league/app/theme"
	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {
	t.Parallel()

	app := chess_league.New()
	app.Addr = "0.0.0.0:8081"

	app.Theme.AppName = "Chess League"
	app.Theme.VenueName = "Our Club"
	app.Theme.PrimaryColour = "#300090"
	app.Theme.SecondaryColour = "#300050"
	app.Theme.TitleBarTextColour = "#ffffff"
	app.Theme.AppIconType = theme.AppIconType_Png

	icon, err := os.ReadFile("../knight.png")
	if err != nil {
		slog.Error("Cannot read the file")
		os.Exit(1)
	}
	app.Theme.AppIcon = icon
	go app.Run()

	for range 10 {
		resp, err := http.Get(fmt.Sprintf("http://%s/", app.Addr))
		if err != nil {
			t.Log("Waiting for server to start...")
			time.Sleep(time.Second / 10)
			continue
		}
		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Contains(t, string(body), app.Theme.AppName)
		return
	}

	require.Fail(t, "Could not check the response of the server")
}
