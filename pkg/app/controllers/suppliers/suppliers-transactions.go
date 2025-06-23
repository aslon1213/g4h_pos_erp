package suppliers

import (
	models "aslon1213/magazin_pos/pkg/repository"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// NewTransaction godoc
// @Summary Create a new transaction for a supplier
// @Description Create a new transaction for a supplier and update financial records
// @Tags suppliers, transactions
// @Accept json
// @Produce json
// @Param branch_id path string true "Branch ID"
// @Param supplier_id path string true "Supplier ID"
// @Param transaction body models.TransactionBase true "Transaction data"
// @Success 200 {object} models.Output
// @Failure 400 {object} models.Output
// @Failure 500 {object} models.Output
// @Router /suppliers/{branch_id}/{supplier_id}/transactions [post]
func (s *SuppliersController) NewTransaction(c *fiber.Ctx) error {
	// if suppliers gets money, that is when we payed money to them,
	// so we need to add money to the supplier's balance and decrease the debt of the branch also decrease the balance of the branch
	// if suppliers pays money, that is when they have given us some product,
	// service or something valuable, so we need to decrease the balance of the supplier and increase the debt of the branch

	branch_id := c.Params("branch_id")

	supplier_id := c.Params("supplier_id")
	log.Info().Str("branch_id", branch_id).Str("supplier_id", supplier_id).Msg("Starting new transaction")

	// Parse transaction data from request body
	var transactionBase models.TransactionBase
	if err := c.BodyParser(&transactionBase); err != nil {
		log.Error().Err(err).Msg("Failed to parse transaction data")
		return c.Status(fiber.StatusBadRequest).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusBadRequest,
		}))
	}

	// Create new transaction
	transaction := models.NewTransaction(&transactionBase, models.InitiatorTypeSupplier, branch_id)
	log.Info().Interface("transaction", transaction).Msg("Created new transaction")

	// Start a session
	session, err := s.DB.Client().StartSession()
	if err != nil {
		log.Error().Err(err).Msg("Failed to start session")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}
	defer session.EndSession(context.Background())

	// Start transaction
	err = session.StartTransaction()
	if err != nil {
		log.Error().Err(err).Msg("Failed to start transaction")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}

	// Insert transaction
	_, err = s.transactionsCollection.InsertOne(context.Background(), transaction)
	if err != nil {
		log.Error().Err(err).Msg("Failed to insert transaction")
		session.AbortTransaction(context.Background())
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}
	log.Info().Msg("Transaction inserted successfully")

	// Update supplier financial data
	update := bson.M{
		"$push": bson.M{
			"financial_data.transactions": transaction,
		},
		"$inc": bson.M{
			"financial_data.balance": func() int32 {
				if transaction.TransactionBase.Type == models.TransactionTypeCredit {
					return int32(transaction.Amount)
				}
				return -int32(transaction.Amount)
			}(),
			"financial_data.total_income": func() int32 {
				if transaction.TransactionBase.Type == models.TransactionTypeCredit {
					return int32(transaction.Amount)
				}
				return 0
			}(),
			"financial_data.total_expenses": func() int32 {
				if transaction.TransactionBase.Type == models.TransactionTypeDebit {
					return int32(transaction.Amount)
				}
				return 0
			}(),
		},
	}

	_, err = s.suppliersCollection.UpdateOne(context.Background(), bson.M{"_id": supplier_id}, update)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update supplier financial data")
		session.AbortTransaction(context.Background())
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}
	log.Info().Msg("Supplier financial data updated successfully")

	// change finance of branch

	// update finance of branch
	filter := bson.M{
		"branch_id": branch_id,
	}
	update_bson := bson.M{
		"finance.debt": func() int32 {
			switch transaction.TransactionBase.Type {
			case models.TransactionTypeDebit:
				return int32(transaction.Amount)
			case models.TransactionTypeCredit:
				return -int32(transaction.Amount)
			}
			return 0
		}(),
	}

	switch transaction.TransactionBase.PaymentMethod {
	case models.PaymentMethodCash:
		update_bson["finance.balance.cash"] = func() int32 {
			if transaction.TransactionBase.Type == models.TransactionTypeCredit {
				return -int32(transaction.Amount)
			}
			return 0
		}()
	case models.PaymentMethodBank, models.OnlineMobileAppPayment:
		update_bson["finance.balance.bank"] = func() int32 {
			if transaction.TransactionBase.Type == models.TransactionTypeCredit {
				return -int32(transaction.Amount)
			}
			return 0
		}()
	case models.OnlineTransfer:
		update_bson["finance.balance.mobile_apps"] = func() int32 {
			if transaction.TransactionBase.Type == models.TransactionTypeCredit {
				return -int32(transaction.Amount)
			}
			return 0
		}()
	}

	update = bson.M{
		"$inc": update_bson,
	}

	res, err := s.financeCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update branch finance")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}
	log.Info().Interface("res", res).Msg("Branch finance updated successfully")

	// Commit transaction
	err = session.CommitTransaction(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}
	log.Info().Msg("Transaction committed successfully")

	return c.Status(fiber.StatusCreated).JSON(models.NewOutput(transaction))
}
