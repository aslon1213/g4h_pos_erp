package app

import (
	"aslon1213/magazin_pos/pkg/app/controllers/finance"
	"aslon1213/magazin_pos/pkg/app/controllers/sales"
	"aslon1213/magazin_pos/pkg/app/controllers/suppliers"
	"aslon1213/magazin_pos/pkg/app/controllers/transactions"
	"aslon1213/magazin_pos/pkg/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Controllers struct {
	Finance      *finance.FinanceController
	Suppliers    *suppliers.SuppliersController
	Transactions *transactions.TransactionsController
	Sales        *sales.SalesTransactionsController
}

func NewControllers(db *mongo.Database) *Controllers {
	log.Debug().Msg("Initializing new controllers")
	controllers := &Controllers{
		Finance:      finance.New(db),
		Suppliers:    suppliers.New(db),
		Transactions: transactions.New(db),
		Sales:        sales.New(db),
	}
	log.Debug().Msg("Controllers initialized successfully")
	return controllers
}

func SetupRoutes(app *fiber.App, controllers *Controllers) {
	log.Debug().Msg("Setting up routes")
	routes.SuppliersRoutes(app, controllers.Suppliers)
	log.Debug().Msg("Suppliers routes set up successfully")
	routes.TransactionsRoutes(app, controllers.Transactions)
	log.Debug().Msg("Transactions routes set up successfully")
	routes.FinanceRoutes(app, controllers.Finance)
	log.Debug().Msg("Finance routes set up successfully")
	routes.SalesRoutes(app, controllers.Sales)
	log.Debug().Msg("All routes set up successfully")
}
