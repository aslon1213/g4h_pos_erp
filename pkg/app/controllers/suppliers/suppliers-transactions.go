package suppliers

import (
	models "aslon1213/magazin_pos/pkg/repository"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *SuppliersController) NewTransaction(c *fiber.Ctx) error {
	// Parse transaction data from request body
	var transactionBase models.TransactionBase
	if err := c.BodyParser(&transactionBase); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusBadRequest,
		}))
	}

	// Create new transaction
	transaction := models.NewTransaction(&transactionBase, models.InitiatorTypeSupplier)

	// Start a session and transaction
	session, err := s.DB.Client().StartSession()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}
	defer session.EndSession(context.Background())

	// Start transaction
	err = session.StartTransaction()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}

	// Insert transaction
	_, err = s.transactionsCollection.InsertOne(context.Background(), transaction)
	if err != nil {
		session.AbortTransaction(context.Background())
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}

	// Update supplier financial data
	update := bson.M{
		"$push": bson.M{
			"financial_data.transactions": transaction,
		},
		"$inc": bson.M{
			"financial_data.balance": func() float64 {
				if transaction.TransactionBase.Type == models.TransactionTypeCredit {
					return transaction.Amount
				}
				return -transaction.Amount
			}(),
			"financial_data.total_income": func() float64 {
				if transaction.TransactionBase.Type == models.TransactionTypeCredit {
					return transaction.Amount
				}
				return 0
			}(),
			"financial_data.total_expenses": func() float64 {
				if transaction.TransactionBase.Type == models.TransactionTypeDebit {
					return transaction.Amount
				}
				return 0
			}(),
		},
	}

	_, err = s.suppliersCollection.UpdateOne(context.Background(), bson.M{"_id": c.Params("id")}, update)
	if err != nil {
		session.AbortTransaction(context.Background())
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}

	// Commit transaction
	err = session.CommitTransaction(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}

	return c.JSON(models.NewOutput(transaction))
}
