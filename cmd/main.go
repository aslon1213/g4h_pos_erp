package main

import "aslon1213/magazin_pos/pkg/app"

func main() {

	app := app.NewApp()

	app.Logger.Info().Msg("Starting server...")

	app.Run()
}
