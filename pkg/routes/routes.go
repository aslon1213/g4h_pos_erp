package routes

import (
	"github.com/aslon1213/go-pos-erp/pkg/controllers/auth"
	"github.com/aslon1213/go-pos-erp/pkg/controllers/finance"
	journal_handlers "github.com/aslon1213/go-pos-erp/pkg/controllers/journals"
	"github.com/aslon1213/go-pos-erp/pkg/controllers/products"
	"github.com/aslon1213/go-pos-erp/pkg/controllers/sales"
	"github.com/aslon1213/go-pos-erp/pkg/controllers/suppliers"
	"github.com/aslon1213/go-pos-erp/pkg/controllers/transactions"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(router *fiber.App, authController *auth.AuthControllers) {

	router.Post("/auth/login", authController.Login)
	router.Post("/auth/register", authController.Register)
	router.Get("api/auth/me", authController.InfoMe)
}

func SuppliersRoutes(router *fiber.App, suppliersController *suppliers.SuppliersController) {
	api := router.Group("/api")
	api.Get("/suppliers", suppliersController.GetSuppliers)                                         // get all suppliers
	api.Get("/suppliers/:id", suppliersController.GetSupplierByID)                                  // get supplier by id
	api.Post("/suppliers", suppliersController.CreateSupplier)                                      // create supplier
	api.Put("/suppliers/:id", suppliersController.UpdateSupplier)                                   // update supplier
	api.Delete("/suppliers/:id", suppliersController.DeleteSupplier)                                // delete supplier
	api.Post("/suppliers/:branch_id/:supplier_id/transactions", suppliersController.NewTransaction) // create transaction
}

func SalesRoutes(router *fiber.App, salesController *sales.SalesTransactionsController) {
	api := router.Group("/api")
	api.Post("/sales/transactions/:branch_id", salesController.CreateSalesTransaction)        // create sales transaction
	api.Delete("/sales/transactions/:transaction_id", salesController.DeleteSalesTransaction) // delete sales transaction
	// sales session routes
	api.Post("/sales/session/branch/:branch_id", salesController.OpenSalesSession)          // open sales session
	api.Post("/sales/session/:session_id/product", salesController.AddProductItemToSession) // add product to session
	api.Post("/sales/session/:session_id/close", salesController.CloseSalesSession)         // close sales session
	api.Get("/sales/session/:session_id", salesController.GetSalesSession)                  // get sales session
	api.Delete("/sales/session/:session_id", salesController.DeleteSalesSession)            // delete sales session
	api.Get("/sales/session/branch/:branch_id", salesController.GetSalesSessionsOfBranch)   // get sales of session
	// router.Get("/sales/branch/:branch_id/sessions", salesController.GetSalesOfSession)

}

func ProductsRoutes(router *fiber.App, productsController *products.ProductsController) {
	api := router.Group("/api")

	api.Post("/products", productsController.CreateProduct)                        // create product
	api.Put("/products/:id", productsController.EditProduct)                       // edit product
	api.Delete("/products/:id", productsController.DeleteProduct)                  // delete product
	api.Get("/products/:id", productsController.GetProductByID)                    // get product by id
	api.Get("/products", productsController.QueryProducts)                         // query products
	api.Post("/products/:id/income", productsController.NewIncome)                 // create income
	api.Post("/products/transfer", productsController.NewTransfer)                 // create transfer
	api.Post("/products/:id/images", productsController.UploadProductImage)        // upload product image
	api.Delete("/products/:id/images/:key", productsController.DeleteProductImage) // delete product image
	api.Get("/products/:id/images", productsController.GetImagesOfProduct)         // get images of product
	api.Get("/products/images/:key", productsController.GetImage)                  // get image
}

func JournalsRoutes(router *fiber.App, journalsController *journal_handlers.JournalHandlers, operationsController *journal_handlers.OperationHandlers) {
	api := router.Group("/api")
	api.Get("/journals/:id", journalsController.GetJournalEntryByID)               // get journal entry by id
	api.Get("/journals/branch/:branch_id", journalsController.QueryJournalEntries) // query journal entries
	api.Post("/journals", journalsController.NewJournalEntry)                      // create journal entry
	api.Post("/journals/:id/close", journalsController.CloseJournalEntry)          // close journal entry
	api.Post("/journals/:id/reopen", journalsController.ReOpenJournalEntry)        // reopen journal entry

	// operations
	api.Post("/journals/:id/operations", operationsController.NewOperationTransaction)                        // create operation transaction
	api.Put("/journals/:id/operations/:operation_id", operationsController.UpdateOperationTransactionByID)    // update operation transaction by id
	api.Delete("/journals/:id/operations/:operation_id", operationsController.DeleteOperationTransactionByID) // delete operation transaction by id
	api.Get("/journals/:id/operations/:operation_id", operationsController.GetOperationTransactionByID)       // get operation transaction by id

}

func InternalExpensesRoutes(router *fiber.App) {

}

func FinanceRoutes(router *fiber.App, financeController *finance.FinanceController) {
	api := router.Group("/api")
	api.Get("/finance/branches", financeController.GetBranches)                            // get all branches
	api.Get("/finance/branch/id/:id", financeController.GetBranchByBranchID)               // get branch by id
	api.Get("/finance/branch/name/:branch_name", financeController.GetFinanceByBranchName) // get branch by name
	api.Get("/finance/id/:id", financeController.GetFinanceByID)                           // get finance by id
	api.Post("/finance", financeController.NewFinanceOfBranch)                             // create new finance of branch

}

func TransactionsRoutes(router *fiber.App, transactionsController *transactions.TransactionsController) {
	api := router.Group("/api")
	api.Get("/transactions/branch/:branch_id", transactionsController.GetTransactionsByQueryParams) // get transactions by query params
	api.Get("/transactions/:id", transactionsController.GetTransactionByID)                         // get transaction by id
	// router.Post("/transactions/:branch_id", transactionsController.Tra)
	api.Put("/transactions/:id", transactionsController.UpdateTransactionByID)            // update transaction by id
	api.Delete("/transactions/:id", transactionsController.DeleteTransactionByID)         // delete transaction by id
	api.Get("/transactions/docs/initiator_type", transactionsController.GetInitiatorType) // get initiator type
	api.Get("/transactions/docs/type", transactionsController.GetTransactionType)         // get transaction type
	api.Get("/transactions/docs/payment_method", transactionsController.GetPaymentMethod) // get payment method
}
