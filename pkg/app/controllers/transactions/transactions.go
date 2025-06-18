package transactions

import (
	"aslon1213/magazin_pos/pkg/app"
	models "aslon1213/magazin_pos/pkg/repository"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TransactionsController struct {
	collection *mongo.Collection
	// cache      cache
}

func NewTransactionsController(app *app.App) *TransactionsController {
	return &TransactionsController{
		collection: app.DB.Database(app.Config.DB.Database).Collection("transactions"),
	}
}

func (t *TransactionsController) GetTransactions(c *fiber.Ctx) error {
	panic("not implemented")
}

func (t *TransactionsController) GetTransactionByID(c *fiber.Ctx) error {

	panic("not implemented")
}

func (t *TransactionsController) CreateTransaction(c *fiber.Ctx) error {
	var transactionBase models.TransactionBase
	if err := c.BodyParser(&transactionBase); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusBadRequest,
		}))
	}

	transaction := models.NewTransaction(&transactionBase, models.InitiatorTypeOther)

	err := models.Finance.CreateTransaction(models.Finance{}, *transaction, t.collection)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}

	return c.JSON(models.NewOutput(transaction))
}

func (t *TransactionsController) UpdateTransaction(c *fiber.Ctx) error {
	idx := c.Params("id")
	amount := c.Query("amount", "")
	description := c.Query("description", "")
	typeOfTransaction := c.Query("type", "")

	query := bson.M{
		"_id": idx,
	}
	set := bson.M{
		"$set": bson.M{},
	}
	if amount != "" {
		set["$set"].(bson.M)["amount"] = amount
	}
	if description != "" {
		set["$set"].(bson.M)["description"] = description
	}
	if typeOfTransaction != "" {
		set["$set"].(bson.M)["type"] = typeOfTransaction
	}

	_, err := t.collection.UpdateOne(context.Background(), query, set)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}
	return c.JSON(models.NewOutput(fiber.Map{
		"message": "Transaction updated successfully",
	}))
}

func (t *TransactionsController) DeleteTransaction(c *fiber.Ctx) error {
	idx := c.Params("id")

	query := bson.M{
		"_id": idx,
	}

	_, err := t.collection.DeleteOne(context.Background(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}
	return c.JSON(models.NewOutput(fiber.Map{
		"message": "Transaction deleted successfully",
	}))
}
