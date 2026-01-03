package main

import "github.com/aslon1213/g4h_pos_erp/pkg/app"

// @title G4H ERP/POS API
// @version 1.0
// @description This is a ERP/POS API for G4H.
// @termsOfService http://swagger.io/terms/

// @contact.name G4H ERP/POS API
// @contact.url https://github.com/aslon1213/g4h_pos_erp
// @contact.email hamidovaslon1@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host g4h.com
func main() {

	app := app.New()

	app.Logger.Info().Msg("Starting server...")

	app.Run()
}
