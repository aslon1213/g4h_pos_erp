package internalexpenses

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type InternalExpensesController struct {
	collection *mongo.Collection
}

func New(db *mongo.Database) *InternalExpensesController {
	return &InternalExpensesController{
		collection: db.Collection("internal_expenses"),
	}
}

func (i *InternalExpensesController) GetInternalExpenses(c *fiber.Ctx) error {
	panic("not implemented")
	return c.SendString("Hello, World!")
}

func (i *InternalExpensesController) GetInternalExpense(c *fiber.Ctx) error {
	panic("not implemented")
	return c.SendString("Hello, World!")
}

func (i *InternalExpensesController) CreateInternalExpense(c *fiber.Ctx) error {
	panic("not implemented")
	return c.SendString("Hello, World!")
}
