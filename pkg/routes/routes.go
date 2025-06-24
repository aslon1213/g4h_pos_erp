package routes

import (
	"aslon1213/magazin_pos/pkg/app/controllers/finance"
	journal_handlers "aslon1213/magazin_pos/pkg/app/controllers/journals"
	"aslon1213/magazin_pos/pkg/app/controllers/sales"
	"aslon1213/magazin_pos/pkg/app/controllers/suppliers"
	"aslon1213/magazin_pos/pkg/app/controllers/transactions"

	"github.com/gofiber/fiber/v2"
)

func SuppliersRoutes(router *fiber.App, suppliersController *suppliers.SuppliersController) {
	router.Get("/suppliers", suppliersController.GetSuppliers) // get all suppliers
	router.Get("/suppliers/:id", suppliersController.GetSupplierByID)
	router.Post("/suppliers", suppliersController.CreateSupplier)
	router.Put("/suppliers/:id", suppliersController.UpdateSupplier)
	router.Delete("/suppliers/:id", suppliersController.DeleteSupplier)
	router.Post("/suppliers/:branch_id/:supplier_id/transactions", suppliersController.NewTransaction)
}

func SalesRoutes(router *fiber.App, salesController *sales.SalesTransactionsController) {
	router.Post("/sales/:branch_id", salesController.CreateSalesTransaction)
}

func JournalsRoutes(router *fiber.App, journalsController *journal_handlers.JournalHandlers, operationsController *journal_handlers.OperationHandlers) {
	router.Get("/journals/:id", journalsController.GetJournalEntryByID)
	router.Get("/journals/branch/:branch_id", journalsController.QueryJournalEntries)
	router.Post("/journals", journalsController.NewJournalEntry)
	router.Post("/journals/:id/close", journalsController.CloseJournalEntry)
	router.Post("/journals/:id/reopen", journalsController.ReOpenJournalEntry)

	// operations
	router.Post("/journals/:id/operations", operationsController.NewOperationTransaction)
	router.Put("/journals/:id/operations/:operation_id", operationsController.UpdateOperationTransactionByID)
	router.Delete("/journals/:id/operations/:operation_id", operationsController.DeleteOperationTransactionByID)
	router.Get("/journals/:id/operations/:operation_id", operationsController.GetOperationTransactionByID)

}

func InternalExpensesRoutes(router *fiber.App) {

}

func FinanceRoutes(router *fiber.App, financeController *finance.FinanceController) {
	router.Get("/finance/branches", financeController.GetBranches)                            // get all branches
	router.Get("/finance/branch/id/:id", financeController.GetBranchByBranchID)               // get branch by id
	router.Get("/finance/branch/name/:branch_name", financeController.GetFinanceByBranchName) // get branch by name
	router.Get("/finance/id/:id", financeController.GetFinanceByID)                           // get finance by id
	router.Post("/finance", financeController.NewFinanceOfBranch)                             // create new finance of branch

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
