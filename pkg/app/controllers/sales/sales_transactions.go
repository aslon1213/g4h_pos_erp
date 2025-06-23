package sales

import (
	models "aslon1213/magazin_pos/pkg/repository"
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
	// sales transaction are always credits
	transaction_base.Type = models.TransactionTypeCredit
	log.Info().Msgf("Transaction base: %+v", transaction_base)
	// accepted types of sales methods: cash, bank, terminal, online_payment, online_transfer, cheque
	if err := models.ValidatePaymentMethod(transaction_base.PaymentMethod); err != nil {
		log.Error().Err(err).Msg("Invalid payment method")
		return c.Status(fiber.StatusBadRequest).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusBadRequest,
		}))
	}

	// start a transaction
	ses, err := s.transactions.Database().Client().StartSession()
	if err != nil {
		log.Error().Err(err).Msg("Failed to start session")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}
	defer ses.EndSession(context.Background())
	if err := ses.StartTransaction(); err != nil {
		log.Error().Err(err).Msg("Failed to start transaction")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}
	// create a new transaction
	transaction := models.NewTransaction(
		&transaction_base,
		models.InitiatorTypeSales,
		branch_id,
	)
	// insert the transaction into the database
	_, err = s.transactions.InsertOne(context.Background(), transaction)
	if err != nil {
		log.Error().Err(err).Msg("Failed to insert transaction")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}
	// update the finance of the branch
	if err := IncrementBalance(s.finances, branch_id, transaction_base); err != nil {
		log.Error().Err(err).Msg("Failed to increment balance")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}

	// commit the transaction
	if err := ses.CommitTransaction(context.Background()); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}
	log.Info().Msg("Transaction created successfully")
	return c.Status(fiber.StatusCreated).JSON(models.NewOutput(transaction))
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
	ses, err := s.transactions.Database().Client().StartSession()
	if err != nil {
		log.Error().Err(err).Msg("Failed to start session")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}
	defer ses.EndSession(context.Background())
	if err := ses.StartTransaction(); err != nil {
		log.Error().Err(err).Msg("Failed to start transaction")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}
	// retrieve the transaction -> delete the transaction -> change the Balance Info of the branch
	transaction := models.Transaction{}
	err = s.transactions.FindOne(context.Background(), bson.M{"_id": transaction_id}).Decode(&transaction)
	if err != nil {
		log.Error().Err(err).Msg("Failed to find transaction")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}
	// delete the transaction
	_, err = s.transactions.DeleteOne(context.Background(), bson.M{"_id": transaction_id})
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete transaction")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}
	// change the Balance Info of the branch
	if err := DecrementBalance(s.finances, transaction.BranchID, transaction); err != nil {
		log.Error().Err(err).Msg("Failed to decrement balance")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}
	// commit the transaction
	if err := ses.CommitTransaction(context.Background()); err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}
	log.Info().Msg("Transaction deleted successfully")
	return c.JSON(models.NewOutput(transaction))
}

func IncrementBalance(finance *mongo.Collection, branch_id string, transaction_base models.TransactionBase) error {
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
	log.Info().Interface("update", update).Interface("filter", filter).Msg("Updating finance of the branch")
	result, err := finance.UpdateOne(
		context.Background(),
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

func DecrementBalance(finance *mongo.Collection, branch_id string, transaction models.Transaction) error {
	filter := bson.M{
		"branch_id": branch_id,
	}
	update := bson.M{
		"$inc": bson.M{},
	}
	switch transaction.PaymentMethod {
	case models.PaymentMethodCash:
		update["$inc"].(bson.M)["finance.balance.cash"] = -transaction.Amount
	case models.PaymentMethodBank:
		update["$inc"].(bson.M)["finance.balance.bank"] = -transaction.Amount
	case models.PaymentMethodTerminal:
		update["$inc"].(bson.M)["finance.balance.terminal"] = -transaction.Amount
	case models.OnlineMobileAppPayment:
		update["$inc"].(bson.M)["finance.balance.mobile_apps"] = -transaction.Amount
	case models.OnlineTransfer:
		update["$inc"].(bson.M)["finance.balance.mobile_apps"] = -transaction.Amount
	}
	result, err := finance.UpdateOne(context.Background(), filter, update)
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
