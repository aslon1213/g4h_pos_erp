package sales

import (
	models "aslon1213/magazin_pos/pkg/repository"
	"aslon1213/magazin_pos/platform/database"
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SalesTransactionsController struct {
	transactions *mongo.Collection
	finances     *mongo.Collection
}

func New(db *mongo.Database) *SalesTransactionsController {
	return &SalesTransactionsController{
		transactions: db.Collection("transactions"),
		finances:     db.Collection("finance"),
	}
}

// CreateSalesTransaction godoc
// @Summary Create a new sales transaction
// @Description Create a new sales transaction for a branch
// @Tags sales
// @Accept json
// @Produce json
// @Param branch_id path string true "Branch ID"
// @Param transaction body models.TransactionBase true "Transaction details"
// @Success 201 {object} models.Output
// @Failure 400 {object} models.Output
// @Failure 500 {object} models.Output
// @Router /sales/{branch_id} [post]
func (s *SalesTransactionsController) CreateSalesTransaction(c *fiber.Ctx) error {
	branch_id := c.Params("branch_id")
	transaction_base := models.TransactionBase{}
	if err := c.BodyParser(&transaction_base); err != nil {
		log.Error().Err(err).Msg("Failed to parse transaction base")
		return c.Status(fiber.StatusBadRequest).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusBadRequest,
		}))
	}
	// start a transaction
	ses, ctx, err := database.StartTransaction(c, s.transactions.Database().Client())
	if err != nil {
		log.Error().Err(err).Msg("Failed to start transaction")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}
	defer ses.EndSession(ctx)
	transaction, err := NewTransaction(ctx, transaction_base, branch_id, s.transactions, s.finances)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}

	// commit the transaction
	if err := ses.CommitTransaction(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}
	return c.Status(fiber.StatusCreated).JSON(models.NewOutput(transaction))
}

func NewTransaction(ctx context.Context, transaction_base models.TransactionBase, branch_id string, transactionsCollection *mongo.Collection, financesCollection *mongo.Collection) (*models.Transaction, error) {
	// sales transaction are always credits
	transaction_base.Type = models.TransactionTypeCredit
	log.Info().Str("Collection", transactionsCollection.Name()).Str("Branch ID", branch_id).Msgf("Creating Sales Transaction with Transaction base: %+v", transaction_base)
	// accepted types of sales methods: cash, bank, terminal, online_payment, online_transfer, cheque
	if err := models.ValidatePaymentMethod(transaction_base.PaymentMethod); err != nil {
		log.Error().Err(err).Msg("Invalid payment method")
		return nil, err
	}

	// create a new transaction
	transaction := models.NewTransaction(
		&transaction_base,
		models.InitiatorTypeSales,
		branch_id,
	)
	// insert the transaction into the database
	_, err := transactionsCollection.InsertOne(ctx, transaction)
	if err != nil {
		log.Error().Err(err).Msg("Failed to insert transaction")
		return nil, err
	}
	// update the finance of the branch
	if err := IncrementBalance(ctx, financesCollection, branch_id, transaction_base); err != nil {
		log.Error().Err(err).Msg("Failed to increment balance")
		return nil, err
	}

	log.Info().Msg("Transaction created successfully")
	return transaction, nil
}

// DeleteSalesTransaction godoc
// @Summary Delete a sales transaction
// @Description Delete a sales transaction by ID
// @Tags sales
// @Accept json
// @Produce json
// @Param transaction_id path string true "Transaction ID"
// @Success 200 {object} models.Output
// @Failure 500 {object} models.Output
// @Router /sales/{transaction_id} [delete]
func (s *SalesTransactionsController) DeleteSalesTransaction(c *fiber.Ctx) error {
	transaction_id := c.Params("transaction_id")

	// start a db transaction
	ses, ctx, err := database.StartTransaction(c, s.transactions.Database().Client())
	if err != nil {
		log.Error().Err(err).Msg("Failed to start transaction")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}
	defer ses.EndSession(ctx)
	//
	transaction, err := DeleteSalesTransaction(ctx, transaction_id, s.transactions, s.finances)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}

	// commit the transaction
	if err := ses.CommitTransaction(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}
	log.Info().Msg("Transaction deleted successfully")
	return c.JSON(models.NewOutput(transaction))
}

func DeleteSalesTransaction(ctx context.Context, transactionID string, transactionsCollection *mongo.Collection, financesCollection *mongo.Collection) (models.Transaction, error) {
	// retrieve the transaction -> delete the transaction -> change the Balance Info of the branch
	transaction := models.Transaction{}
	err := transactionsCollection.FindOne(ctx, bson.M{"_id": transactionID}).Decode(&transaction)
	if err != nil {
		log.Error().Err(err).Msg("Failed to find transaction")
		return models.Transaction{}, err
	}
	// delete the transaction
	_, err = transactionsCollection.DeleteOne(ctx, bson.M{"_id": transactionID})
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete transaction")
		return models.Transaction{}, err
	}
	// change the Balance Info of the branch
	if err := DecrementBalance(ctx, financesCollection, transaction.BranchID, transaction); err != nil {
		log.Error().Err(err).Msg("Failed to decrement balance")
		return models.Transaction{}, err
	}
	return transaction, nil
}

func IncrementBalance(ctx context.Context, finance *mongo.Collection, branch_id string, transaction_base models.TransactionBase) error {
	filter := bson.M{
		"branch_id": branch_id,
	}
	update := bson.M{
		"$inc": bson.M{},
	}
	switch transaction_base.PaymentMethod {
	case models.PaymentMethodCash:
		update["$inc"].(bson.M)["finance.balance.cash"] = int(transaction_base.Amount)
	case models.PaymentMethodBank:
		update["$inc"].(bson.M)["finance.balance.bank"] = int(transaction_base.Amount)
	case models.PaymentMethodTerminal:
		update["$inc"].(bson.M)["finance.balance.terminal"] = int(transaction_base.Amount)
	case models.OnlineMobileAppPayment:
		update["$inc"].(bson.M)["finance.balance.mobile_apps"] = int(transaction_base.Amount)
	case models.OnlineTransfer:
		update["$inc"].(bson.M)["finance.balance.mobile_apps"] = int(transaction_base.Amount)
	}
	log.Info().Interface("update", update).Interface("filter", filter).Str("Collection", finance.Name()).Msg("Updating finance of the branch")
	result, err := finance.UpdateOne(
		ctx,
		filter,
		update,
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update finance")
		return err
	}
	if result.MatchedCount == 0 {
		log.Error().Msg("No finance found for the branch")
		return errors.New("no finance found for the branch")
	}
	return nil
}

func DecrementBalance(ctx context.Context, finance *mongo.Collection, branch_id string, transaction models.Transaction) error {
	filter := bson.M{
		"branch_id": branch_id,
	}
	update := bson.M{
		"$inc": bson.M{},
	}
	switch transaction.PaymentMethod {
	case models.PaymentMethodCash:
		update["$inc"].(bson.M)["finance.balance.cash"] = -int32(transaction.Amount)
	case models.PaymentMethodBank:
		update["$inc"].(bson.M)["finance.balance.bank"] = -int32(transaction.Amount)
	case models.PaymentMethodTerminal:
		update["$inc"].(bson.M)["finance.balance.terminal"] = -int32(transaction.Amount)
	case models.OnlineMobileAppPayment:
		update["$inc"].(bson.M)["finance.balance.mobile_apps"] = -int32(transaction.Amount)
	case models.OnlineTransfer:
		update["$inc"].(bson.M)["finance.balance.mobile_apps"] = -int32(transaction.Amount)
	}
	result, err := finance.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update finance")
		return err
	}
	if result.MatchedCount == 0 {
		log.Error().Msg("No finance found for the branch")
		return errors.New("no finance found for the branch")
	}
	return nil
}
