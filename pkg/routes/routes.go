package routes

import (
	"aslon1213/magazin_pos/pkg/app/controllers/finance"
	"aslon1213/magazin_pos/pkg/app/controllers/sales"
	"aslon1213/magazin_pos/pkg/app/controllers/suppliers"
	"aslon1213/magazin_pos/pkg/app/controllers/transactions"

	"github.com/gofiber/fiber/v2"
)

func SuppliersRoutes(router *fiber.App, suppliersController *suppliers.SuppliersController) {
	router.Get("/suppliers", suppliersController.GetSuppliers)
	router.Get("/suppliers/:id", suppliersController.GetSupplierByID)
	router.Post("/suppliers", suppliersController.CreateSupplier)
	router.Put("/suppliers/:id", suppliersController.UpdateSupplier)
	router.Delete("/suppliers/:id", suppliersController.DeleteSupplier)
	router.Post("/suppliers/:branch_id/:supplier_id/transactions", suppliersController.NewTransaction)
}

func SalesRoutes(router *fiber.App, salesController *sales.SalesTransactionsController) {
	router.Post("/sales/:branch_id", salesController.CreateSalesTransaction)
}

func InternalExpensesRoutes(router *fiber.App) {

}

func FinanceRoutes(router *fiber.App, financeController *finance.FinanceController) {
	router.Get("/finance", financeController.GetBranches)
	router.Get("/finance/id/:id", financeController.GetBranchByBranchID)
	router.Get("/finance/name/:branch_name", financeController.GetFinanceByBranchName)
	router.Post("/finance", financeController.NewFinanceOfBranch)

}

func TransactionsRoutes(router *fiber.App, transactionsController *transactions.TransactionsController) {
	router.Get("/transactions/branch/:branch_id", transactionsController.GetTransactionsByQueryParams)
	router.Get("/transactions/:id", transactionsController.GetTransactionByID)
	// router.Post("/transactions/:branch_id", transactionsController.Tra)
	router.Put("/transactions/:id", transactionsController.UpdateTransactionByID)
	router.Delete("/transactions/:id", transactionsController.DeleteTransactionByID)
	router.Get("/transactions/docs/initiator_type", transactionsController.GetInitiatorType)
	router.Get("/transactions/docs/type", transactionsController.GetTransactionType)
	router.Get("/transactions/docs/payment_method", transactionsController.GetPaymentMethod)

}
