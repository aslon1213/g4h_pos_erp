package app

import (
	"github.com/aslon1213/go-pos-erp/pkg/app/controllers/finance"
	journal_handlers "github.com/aslon1213/go-pos-erp/pkg/app/controllers/journals"
	"github.com/aslon1213/go-pos-erp/pkg/app/controllers/products"
	"github.com/aslon1213/go-pos-erp/pkg/app/controllers/sales"
	"github.com/aslon1213/go-pos-erp/pkg/app/controllers/suppliers"
	"github.com/aslon1213/go-pos-erp/pkg/app/controllers/transactions"
	"github.com/aslon1213/go-pos-erp/pkg/routes"
	"github.com/aslon1213/go-pos-erp/platform/cache"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Controllers struct {
	Finance      *finance.FinanceController
	Suppliers    *suppliers.SuppliersController
	Transactions *transactions.TransactionsController
	Sales        *sales.SalesTransactionsController
	Journals     *journal_handlers.JournalHandlers
	Operations   *journal_handlers.OperationHandlers
	Products     *products.ProductsController
}

func NewControllers(db *mongo.Database, cache *cache.Cache) *Controllers {
	log.Debug().Msg("Initializing new controllers")
	controllers := &Controllers{
		Finance:      finance.New(db),
		Suppliers:    suppliers.New(db),
		Transactions: transactions.New(db),
		Sales:        sales.New(db, cache),
		Journals:     journal_handlers.New(db, cache),
		Operations:   journal_handlers.NewOperationsHandler(db, cache),
		Products:     products.NewProductsController(db),
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
	log.Debug().Msg("Sales routes set up successfully")
	routes.JournalsRoutes(app, controllers.Journals, controllers.Operations)
	log.Debug().Msg("Journals routes set up successfully")
	routes.ProductsRoutes(app, controllers.Products)
	log.Debug().Msg("Products routes set up successfully")
	log.Debug().Msg("All routes set up successfully")
}
