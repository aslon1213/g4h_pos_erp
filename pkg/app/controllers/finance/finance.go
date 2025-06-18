package finance

import (
	"aslon1213/magazin_pos/pkg/app"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type FinanceController struct {
	collection *mongo.Collection
}

func NewFinanceController(app *app.App) *FinanceController {
	return &FinanceController{
		collection: app.DB.Database(app.Config.DB.Database).Collection("finance"),
	}
}

func (f *FinanceController) GetFinance(c *fiber.Ctx) error {
	panic("not implemented")
	return c.SendString("Hello, World!")
}
