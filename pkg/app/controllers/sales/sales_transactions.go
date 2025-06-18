package sales

import (
	"aslon1213/magazin_pos/pkg/app"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SalesTransactionsController struct {
	collection *mongo.Collection
}

func NewSalesTransactionsController(app *app.App) *SalesTransactionsController {
	return &SalesTransactionsController{
		collection: app.DB.Database(app.Config.DB.Database).Collection("sales_transactions"),
	}
}

func (s *SalesTransactionsController) GetSalesTransactions(c *fiber.Ctx) error {
	panic("not implemented")
	return c.SendString("Hello, World!")
}
