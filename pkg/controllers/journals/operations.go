package journal_handlers

import (
	"context"

	"github.com/aslon1213/go-pos-erp/pkg/controllers/sales"
	"github.com/aslon1213/go-pos-erp/pkg/controllers/suppliers"
	"github.com/aslon1213/go-pos-erp/pkg/middleware"
	models "github.com/aslon1213/go-pos-erp/pkg/repository"
	"github.com/aslon1213/go-pos-erp/platform/cache"
	"github.com/aslon1213/go-pos-erp/platform/database"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OperationHandlers struct {
	ctx                    context.Context
	OperationsCollection   *mongo.Collection
	FinancesCollection     *mongo.Collection
	JournalsCollection     *mongo.Collection
	SuppliersCollections   *mongo.Collection
	TransactionsCollection *mongo.Collection
	ActivitiesCollection   *mongo.Collection
	RedisClient            *cache.Cache
}

func NewOperationsHandler(db *mongo.Database, cache *cache.Cache) *OperationHandlers {
	ctx := context.Background()
	operationsCollection := db.Collection("operations")
	journalsCollection := db.Collection("journals")
	suppliersCollections := db.Collection("suppliers")
	financesCollection := db.Collection("finance")
	transactionsCollection := db.Collection("transactions")
	activitiesCollection := db.Collection("activities")
	return &OperationHandlers{
		ctx:                    ctx,
		OperationsCollection:   operationsCollection,
		JournalsCollection:     journalsCollection,
		FinancesCollection:     financesCollection,
		SuppliersCollections:   suppliersCollections,
		TransactionsCollection: transactionsCollection,
		RedisClient:            cache,
		ActivitiesCollection:   activitiesCollection,
	}
}

// NewOperationTransaction godoc
// @Security BearerAuth
// @Summary Create a new operation transaction
// @Description Create a new transaction and update the journal
// @Tags journals/operations
// @Accept json
// @Produce json
// @Param transaction body models.JournalOperationInput true "Transaction data"
// @Param journal_id path string true "Journal ID"
// @Success 200 {object} models.Output
// @Failure 400 {object} models.Output
// @Failure 500 {object} models.Output
// @Router /api/journals/{journal_id}/operations [post]
func (o *OperationHandlers) NewOperationTransaction(c *fiber.Ctx) error {
	transaction := models.JournalOperationInput{}
	if err := c.BodyParser(&transaction); err != nil {
		log.Error().Err(err).Msg("Failed to parse transaction data")
		return c.Status(fiber.StatusBadRequest).JSON(models.NewOutput([]interface{}{}, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusBadRequest,
		}))
	}

	log.Info().Interface("transaction", transaction).Msg("Transaction data")

	// transaction_created := &models.Transaction{}
	// start a new session and transaction
	ses, ctx, err := database.StartTransaction(o.JournalsCollection.Database().Client())
	if err != nil {
		log.Error().Err(err).Msg("Failed to start transaction")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput([]interface{}{}, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}
	defer ses.EndSession(ctx)
	// log activity

	ids := []string{}
	log.Info().Msg("Fetching journal by ID")
	journal, err := FetchJournalByID(ctx, c, true, o.JournalsCollection)
	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch journal by ID")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput([]interface{}{}, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}

	if !transaction.SupplierTransaction {
		transaction.TransactionBase.Type = models.TransactionTypeCredit
		log.Info().Msg("Creating new sales transaction")
		sales_transaction, err := sales.NewTransaction(ctx, transaction.TransactionBase, journal.Branch.ID, o.TransactionsCollection, o.FinancesCollection)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create sales transaction")

			return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput([]interface{}{}, models.Error{
				Message: err.Error(),
				Code:    fiber.StatusInternalServerError,
			}))
		}
		// idx := sales_transaction.ID
		ids = append(ids, sales_transaction.ID)
	} else {
		transaction.TransactionBase.Type = models.TransactionTypeCredit
		if transaction.SupplierID == "" {
			log.Error().Msg("Supplier ID is required")

			return c.Status(fiber.StatusBadRequest).JSON(models.NewOutput([]interface{}{}, models.Error{
				Message: "Supplier ID is required",
				Code:    fiber.StatusBadRequest,
			}))
		}
		log.Info().Msg("Creating new supplier transaction")
		supplier_transaction, err := suppliers.NewSupplierTransaction(ctx, transaction.TransactionBase, transaction.SupplierID, journal.Branch.ID, o.TransactionsCollection, o.JournalsCollection, o.SuppliersCollections)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create supplier transaction")

			return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput([]interface{}{}, models.Error{
				Message: err.Error(),
				Code:    fiber.StatusInternalServerError,
			}))
		}
		ids = append(ids, supplier_transaction.ID)
	}

	log.Info().Msg("Updating journal total")
	_, err = o.JournalsCollection.UpdateByID(
		ctx,
		journal.ID,
		bson.M{
			"$inc": bson.M{
				"total": transaction.Amount,
			},
			"$push": bson.M{
				"operations": bson.M{
					"$each": ids,
				},
			},
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update journal total")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput([]interface{}{}, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}

	// commit the transaction
	if err := ses.CommitTransaction(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput([]interface{}{}, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}

	middleware.LogActivityWithCtx(c, middleware.ActivityTypeCreateTransaction, transaction, o.ActivitiesCollection)

	log.Info().Msg("Transaction created successfully")
	return c.Status(fiber.StatusCreated).JSON(models.NewOutput(journal))
}

// UpdateOperationTransactionByID godoc
// @Security BearerAuth
// @Summary Update an operation transaction
// @Description Update an operation transaction by ID
// @Tags journals/operations
// @Accept json
// @Produce json
// @Param id path string true "Operation ID"
// @Param journal_id path string true "Journal ID"
// @Success 200 {object} models.Output
// @Failure 400 {object} models.Output
// @Failure 500 {object} models.Output
// @Router /api/journals/{journal_id}/operations/{id} [put]
func (o *OperationHandlers) UpdateOperationTransactionByID(c *fiber.Ctx) error {
	// log activity
	middleware.LogActivityWithCtx(c, middleware.ActivityTypeEditOperation, c.Params("id"), o.ActivitiesCollection)

	panic("Not implemented")

}

// DeleteOperationTransactionByID godoc
// @Security BearerAuth
// @Summary Delete an operation transaction
// @Description Delete an operation transaction by ID
// @Tags journals/operations
// @Accept json
// @Produce json
// @Param id path string true "Operation ID"
// @Param journal_id path string true "Journal ID"
// @Success 200 {object} models.Output
// @Failure 400 {object} models.Output
// @Failure 500 {object} models.Output
// @Router /api/journals/{journal_id}/operations/{id} [delete]
func (o *OperationHandlers) DeleteOperationTransactionByID(c *fiber.Ctx) error {
	// log activity
	middleware.LogActivityWithCtx(c, middleware.ActivityTypeDeleteOperation, c.Params("id"), o.ActivitiesCollection)

	panic("Not implemented")
}

// GetOperationTransactionByID godoc
// @Security BearerAuth
// @Summary Get an operation transaction by ID
// @Description Get an operation transaction by ID
// @Tags journals/operations
// @Accept json
// @Produce json
// @Param id path string true "Operation ID"
// @Param journal_id path string true "Journal ID"
// @Success 200 {object} models.Output
// @Failure 400 {object} models.Output
// @Failure 500 {object} models.Output
// @Router /api/journals/{journal_id}/operations/{id} [get]
func (o *OperationHandlers) GetOperationTransactionByID(c *fiber.Ctx) error {
	// log activity

	panic("Not implemented")
}
