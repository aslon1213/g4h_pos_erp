package suppliers

import (
	"aslon1213/magazin_pos/pkg/app"
	models "aslon1213/magazin_pos/pkg/repository"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SuppliersController struct {
	suppliersCollection    *mongo.Collection
	transactionsCollection *mongo.Collection
	DB                     *mongo.Database
}

func NewSuppliersController(app *app.App) *SuppliersController {
	return &SuppliersController{
		suppliersCollection:    app.DB.Database(app.Config.DB.Database).Collection("companies"),
		transactionsCollection: app.DB.Database(app.Config.DB.Database).Collection("transactions"),
		DB:                     app.DB.Database(app.Config.DB.Database),
	}
}

func (s *SuppliersController) GetSuppliers(c *fiber.Ctx) error {
	log.Debug().Msg("Getting all suppliers")

	var suppliers []models.Supplier
	cursor, err := s.suppliersCollection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to find suppliers")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}

	if err := cursor.All(context.Background(), &suppliers); err != nil {
		log.Error().Err(err).Msg("Failed to decode suppliers")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}

	log.Debug().Int("count", len(suppliers)).Msg("Successfully retrieved suppliers")
	return c.JSON(models.NewOutput(suppliers))
}

func (s *SuppliersController) GetSupplierByID(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Debug().Str("id", id).Msg("Getting supplier by ID")

	var supplier models.Supplier
	err := s.suppliersCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&supplier)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Debug().Str("id", id).Msg("Supplier not found")
			return c.Status(fiber.StatusNotFound).JSON(models.NewOutput(nil, models.Error{
				Message: "Supplier not found",
				Code:    fiber.StatusNotFound,
			}))
		}
		log.Error().Err(err).Str("id", id).Msg("Failed to find supplier")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}

	log.Debug().Str("id", id).Msg("Successfully retrieved supplier")
	return c.JSON(models.NewOutput(supplier))
}

func (s *SuppliersController) CreateSupplier(c *fiber.Ctx) error {
	log.Debug().Msg("Creating new supplier")

	var supplierBase models.SupplierBase
	if err := c.BodyParser(&supplierBase); err != nil {
		log.Error().Err(err).Msg("Failed to parse supplier data")
		return c.Status(fiber.StatusBadRequest).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusBadRequest,
		}))
	}

	now := time.Now()
	supplier := models.Supplier{
		SupplierBase: supplierBase,
		ID:           uuid.New().String(),
		FinancialData: models.FinancialData{
			Balance:       0,
			Transactions:  []models.Transaction{},
			TotalIncome:   0,
			TotalExpenses: 0,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err := s.suppliersCollection.InsertOne(context.Background(), supplier)
	if err != nil {
		log.Error().Err(err).Str("id", supplier.ID).Msg("Failed to insert supplier")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}

	log.Debug().Str("id", supplier.ID).Msg("Successfully created supplier")
	return c.Status(fiber.StatusCreated).JSON(models.NewOutput(supplier))
}

func (s *SuppliersController) UpdateSupplier(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Debug().Str("id", id).Msg("Updating supplier")

	var supplierBase models.SupplierBase
	if err := c.BodyParser(&supplierBase); err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to parse supplier update data")
		return c.Status(fiber.StatusBadRequest).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusBadRequest,
		}))
	}

	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}
	if supplierBase.Email != "" {
		update["$set"].(bson.M)["email"] = supplierBase.Email
	}
	if supplierBase.INN != "" {
		update["$set"].(bson.M)["inn"] = supplierBase.INN
	}
	if supplierBase.Name != "" {
		update["$set"].(bson.M)["name"] = supplierBase.Name
	}
	if supplierBase.Address != "" {
		update["$set"].(bson.M)["address"] = supplierBase.Address
	}
	if supplierBase.Phone != "" {
		update["$set"].(bson.M)["phone"] = supplierBase.Phone
	}
	if supplierBase.Notes != "" {
		update["$set"].(bson.M)["notes"] = supplierBase.Notes
	}
	if supplierBase.Branch != "" {
		update["$set"].(bson.M)["branch"] = supplierBase.Branch
	}

	result, err := s.suppliersCollection.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to update supplier")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}

	if result.MatchedCount == 0 {
		log.Debug().Str("id", id).Msg("Supplier not found for update")
		return c.Status(fiber.StatusNotFound).JSON(models.NewOutput(nil, models.Error{
			Message: "Supplier not found",
			Code:    fiber.StatusNotFound,
		}))
	}

	log.Debug().Str("id", id).Msg("Successfully updated supplier")
	return c.JSON(models.NewOutput(fiber.Map{"message": "Supplier updated successfully"}))
}

func (s *SuppliersController) DeleteSupplier(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Debug().Str("id", id).Msg("Deleting supplier")

	result, err := s.suppliersCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to delete supplier")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput(nil, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}

	if result.DeletedCount == 0 {
		log.Debug().Str("id", id).Msg("Supplier not found for deletion")
		return c.Status(fiber.StatusNotFound).JSON(models.NewOutput(nil, models.Error{
			Message: "Supplier not found",
			Code:    fiber.StatusNotFound,
		}))
	}

	log.Debug().Str("id", id).Msg("Successfully deleted supplier")
	return c.JSON(models.NewOutput(fiber.Map{"message": "Supplier deleted successfully"}))
}
