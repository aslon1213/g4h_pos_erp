package internalexpenses

import (
	"aslon1213/magazin_pos/pkg/app"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type InternalExpensesController struct {
	collection *mongo.Collection
}

func NewInternalExpensesController(app *app.App) *InternalExpensesController {
	return &InternalExpensesController{
		collection: app.DB.Database(app.Config.DB.Database).Collection("internal_expenses"),
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
