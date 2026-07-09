package main

import (
	chess_league "github.com/Nag-s-Head/chess-league/app"
)

func main() {
	app := chess_league.New()
	app.Theme.AppName = "Nag's Knights"
	app.Theme.VenueName = "The Nag's Head"
	app.Theme.PrimaryColour = "#ec003f"
	app.Theme.SecondaryColour = "#ffa1ad"
	app.Theme.TitleBarTextColour = "#ffffff"
	app.Run()
}
