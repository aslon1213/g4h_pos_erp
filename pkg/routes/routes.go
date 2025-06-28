package routes

import (
	"github.com/aslon1213/go-pos-erp/pkg/controllers/finance"
	journal_handlers "github.com/aslon1213/go-pos-erp/pkg/controllers/journals"
	"github.com/aslon1213/go-pos-erp/pkg/controllers/products"
	"github.com/aslon1213/go-pos-erp/pkg/controllers/sales"
	"github.com/aslon1213/go-pos-erp/pkg/controllers/suppliers"
	"github.com/aslon1213/go-pos-erp/pkg/controllers/transactions"

	"github.com/gofiber/fiber/v2"
)

func SuppliersRoutes(router *fiber.App, suppliersController *suppliers.SuppliersController) {
	router.Get("/suppliers", suppliersController.GetSuppliers)                                         // get all suppliers
	router.Get("/suppliers/:id", suppliersController.GetSupplierByID)                                  // get supplier by id
	router.Post("/suppliers", suppliersController.CreateSupplier)                                      // create supplier
	router.Put("/suppliers/:id", suppliersController.UpdateSupplier)                                   // update supplier
	router.Delete("/suppliers/:id", suppliersController.DeleteSupplier)                                // delete supplier
	router.Post("/suppliers/:branch_id/:supplier_id/transactions", suppliersController.NewTransaction) // create transaction
}

func SalesRoutes(router *fiber.App, salesController *sales.SalesTransactionsController) {
	router.Post("/sales/transactions/:branch_id", salesController.CreateSalesTransaction)        // create sales transaction
	router.Delete("/sales/transactions/:transaction_id", salesController.DeleteSalesTransaction) // delete sales transaction
	// sales session routes
	router.Post("/sales/session/branch/:branch_id", salesController.OpenSalesSession)          // open sales session
	router.Post("/sales/session/:session_id/product", salesController.AddProductItemToSession) // add product to session
	router.Post("/sales/session/:session_id/close", salesController.CloseSalesSession)         // close sales session
	router.Get("/sales/session/:session_id", salesController.GetSalesSession)                  // get sales session
	router.Delete("/sales/session/:session_id", salesController.DeleteSalesSession)            // delete sales session
	router.Get("/sales/session/branch/:branch_id", salesController.GetSalesSessionsOfBranch)   // get sales of session
	// router.Get("/sales/branch/:branch_id/sessions", salesController.GetSalesOfSession)

}

func ProductsRoutes(router *fiber.App, productsController *products.ProductsController) {

	router.Post("/products", productsController.CreateProduct)                        // create product
	router.Put("/products/:id", productsController.EditProduct)                       // edit product
	router.Delete("/products/:id", productsController.DeleteProduct)                  // delete product
	router.Get("/products/:id", productsController.GetProductByID)                    // get product by id
	router.Get("/products", productsController.QueryProducts)                         // query products
	router.Post("/products/:id/income", productsController.NewIncome)                 // create income
	router.Post("/products/transfer", productsController.NewTransfer)                 // create transfer
	router.Post("/products/:id/images", productsController.UploadProductImage)        // upload product image
	router.Delete("/products/:id/images/:key", productsController.DeleteProductImage) // delete product image
	router.Get("/products/:id/images", productsController.GetImagesOfProduct)         // get images of product
	router.Get("/products/images/:key", productsController.GetImage)                  // get image
}

func JournalsRoutes(router *fiber.App, journalsController *journal_handlers.JournalHandlers, operationsController *journal_handlers.OperationHandlers) {
	router.Get("/journals/:id", journalsController.GetJournalEntryByID)               // get journal entry by id
	router.Get("/journals/branch/:branch_id", journalsController.QueryJournalEntries) // query journal entries
	router.Post("/journals", journalsController.NewJournalEntry)                      // create journal entry
	router.Post("/journals/:id/close", journalsController.CloseJournalEntry)          // close journal entry
	router.Post("/journals/:id/reopen", journalsController.ReOpenJournalEntry)        // reopen journal entry

	// operations
	router.Post("/journals/:id/operations", operationsController.NewOperationTransaction)                        // create operation transaction
	router.Put("/journals/:id/operations/:operation_id", operationsController.UpdateOperationTransactionByID)    // update operation transaction by id
	router.Delete("/journals/:id/operations/:operation_id", operationsController.DeleteOperationTransactionByID) // delete operation transaction by id
	router.Get("/journals/:id/operations/:operation_id", operationsController.GetOperationTransactionByID)       // get operation transaction by id

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
	router.Get("/transactions/branch/:branch_id", transactionsController.GetTransactionsByQueryParams) // get transactions by query params
	router.Get("/transactions/:id", transactionsController.GetTransactionByID)                         // get transaction by id
	// router.Post("/transactions/:branch_id", transactionsController.Tra)
	router.Put("/transactions/:id", transactionsController.UpdateTransactionByID)            // update transaction by id
	router.Delete("/transactions/:id", transactionsController.DeleteTransactionByID)         // delete transaction by id
	router.Get("/transactions/docs/initiator_type", transactionsController.GetInitiatorType) // get initiator type
	router.Get("/transactions/docs/type", transactionsController.GetTransactionType)         // get transaction type
	router.Get("/transactions/docs/payment_method", transactionsController.GetPaymentMethod) // get payment method
}
