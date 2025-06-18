package routes

import (
	"aslon1213/magazin_pos/pkg/app"
	"aslon1213/magazin_pos/pkg/app/controllers/suppliers"

	"github.com/gofiber/fiber/v2"
)

func CustomersRoutes(router *fiber.App) {

}

func SuppliersRoutes(router *fiber.App, app *app.App) {
	suppliersController := suppliers.NewSuppliersController(app)
	router.Get("/suppliers", suppliersController.GetSuppliers)
	router.Get("/suppliers/:id", suppliersController.GetSupplierByID)
	router.Post("/suppliers", suppliersController.CreateSupplier)
	router.Put("/suppliers/:id", suppliersController.UpdateSupplier)
	router.Delete("/suppliers/:id", suppliersController.DeleteSupplier)
	router.Post("/suppliers/:id/transactions", suppliersController.NewTransaction)
}

func SalesRoutes(router *fiber.App) {

}

func InternalExpensesRoutes(router *fiber.App) {

}

func FinanceRoutes(router *fiber.App) {

}

func TransactionsRoutes(router *fiber.App) {

}

func SetupRoutes(router *fiber.App, app *app.App) {
	CustomersRoutes(router)
	SuppliersRoutes(router, app)
	SalesRoutes(router)
	InternalExpensesRoutes(router)
	FinanceRoutes(router)
	TransactionsRoutes(router)
}
