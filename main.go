package main

import (
	"github.com/Nag-s-Head/chess-league/chess_league"
)

func main() {
	app := chess_league.New()
	app.Theme.AppName = "Nag's Knights"
	app.Theme.VenueName = "The Nags' Head"
	app.Theme.PrimaryColour = "#ec003f"
	app.Theme.SecondaryColour = "#ffa1ad"
	app.Run()
}
