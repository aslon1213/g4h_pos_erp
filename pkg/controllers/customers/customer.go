package customers

import (
	"context"
	"time"

	models "github.com/aslon1213/go-pos-erp/pkg/repository"
	"github.com/aslon1213/go-pos-erp/platform/cache"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type CustomersController struct {
	customersCollection *mongo.Collection
	salesCollection     *mongo.Collection
	bnplCollection      *mongo.Collection
	DB                  *mongo.Database
}

func New(db *mongo.Database, cache *cache.Cache) *CustomersController {
	return &CustomersController{
		customersCollection: db.Collection("customers"),
		salesCollection:     db.Collection("sales"),
		bnplCollection:      db.Collection("bnpl"),
		DB:                  db,
	}
}

type CustomerQuery struct {
	Name    string `query:"name"`
	Phone   string `query:"phone"`
	Address string `query:"address"`
	Page    int    `query:"page" default:"1"`
	Count   int    `query:"count" default:"10"`
}

// GetCustomers godoc
// @Security BearerAuth
// @Summary Get all customers
// @Description Get all customers from the database
// @Tags customers
// @Accept json
// @Produce json
// @Param name query string false "Customer name"
// @Param phone query string false "Customer phone"
// @Param address query string false "Customer address"
// @Param page query int false "Page number"
// @Param count query int false "Number of customers per page"
// @Success 200 {object} models.Output
// @Failure 500 {object} models.Output
// @Router /api/customers [get]
func (ctrl *CustomersController) GetCustomers(c *fiber.Ctx) error {
	var query CustomerQuery
	if err := c.QueryParser(&query); err != nil {
		log.Error().Err(err).Msg("Failed to parse customer query")
		return c.Status(fiber.StatusBadRequest).JSON(models.NewOutput([]interface{}{}, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusBadRequest,
		}))
	}

	// adjust the page and count
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Count <= 0 {
		query.Count = 10
	}

	// build filter
	filter := bson.M{}
	skip := (query.Page - 1) * query.Count
	limit := query.Count

	if query.Name != "" {
		filter["name"] = bson.M{"$regex": query.Name, "$options": "i"}
	}
	if query.Phone != "" {
		filter["phone"] = bson.M{"$regex": query.Phone, "$options": "i"}
	}
	if query.Address != "" {
		filter["address"] = bson.M{"$regex": query.Address, "$options": "i"}
	}

	log.Debug().Interface("filter", filter).Msg("Getting customers with filter")

	customers := make([]models.Customer, 0)
	cursor, err := ctrl.customersCollection.Find(context.Background(), filter, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
	if err != nil {
		log.Error().Err(err).Msg("Failed to find customers")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput([]interface{}{}, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &customers); err != nil {
		log.Error().Err(err).Msg("Failed to decode customers")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput([]interface{}{}, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}

	log.Debug().Int("count", len(customers)).Msg("Successfully retrieved customers")

	// query total customers number from the database
	total, err := ctrl.customersCollection.CountDocuments(context.Background(), filter)
	if err != nil {
		log.Error().Err(err).Msg("Failed to count customers")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput([]interface{}{}, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}
	total_pages := int(total) / query.Count
	if int(total)%query.Count != 0 {
		total_pages++
	}
	output := models.NewCustomerQueryOutput(customers, total_pages, query.Page, query.Count)
	return c.JSON(output)
}

// GetCustomerByID godoc
// @Security BearerAuth
// @Summary Get a customer by ID
// @Description Get a customer by its ID
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} models.Output
// @Failure 404 {object} models.Output
// @Failure 500 {object} models.Output
// @Router /api/customers/{id} [get]
func (ctrl *CustomersController) GetCustomerByID(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Debug().Str("id", id).Msg("Getting customer by ID")

	var customer models.Customer
	err := ctrl.customersCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Debug().Str("id", id).Msg("Customer not found")
			return c.Status(fiber.StatusNotFound).JSON(models.NewOutput([]interface{}{}, models.Error{
				Message: "Customer not found",
				Code:    fiber.StatusNotFound,
			}))
		}
		log.Error().Err(err).Str("id", id).Msg("Failed to find customer")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput([]interface{}{}, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}

	log.Debug().Str("id", id).Msg("Successfully retrieved customer")
	return c.JSON(models.NewOutput([]models.Customer{customer}))
}

// CreateCustomer godoc
// @Security BearerAuth
// @Summary Create a new customer
// @Description Create a new customer in the database
// @Tags customers
// @Accept json
// @Produce json
// @Param customer body models.CustomerBase true "Customer data"
// @Success 201 {object} models.Output
// @Failure 400 {object} models.Output
// @Failure 500 {object} models.Output
// @Router /api/customers [post]
func (ctrl *CustomersController) CreateCustomer(c *fiber.Ctx) error {
	log.Debug().Msg("Creating new customer")

	var customerBase models.CustomerBase
	if err := c.BodyParser(&customerBase); err != nil {
		log.Error().Err(err).Msg("Failed to parse customer data")
		return c.Status(fiber.StatusBadRequest).JSON(models.NewOutput([]interface{}{}, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusBadRequest,
		}))
	}

	// Validate required fields
	if customerBase.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewOutput([]interface{}{}, models.Error{
			Message: "Customer name is required",
			Code:    fiber.StatusBadRequest,
		}))
	}
	if customerBase.Phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewOutput([]interface{}{}, models.Error{
			Message: "Customer phone is required",
			Code:    fiber.StatusBadRequest,
		}))
	}

	// Check if customer with same phone already exists
	var existingCustomer models.Customer
	err := ctrl.customersCollection.FindOne(context.Background(), bson.M{"phone": customerBase.Phone}).Decode(&existingCustomer)
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(models.NewOutput([]interface{}{}, models.Error{
			Message: "Customer with this phone number already exists",
			Code:    fiber.StatusConflict,
		}))
	}

	// Create new customer
	customer := models.Customer{
		ID:              uuid.New().String(),
		CustomerBase:    customerBase,
		BNPLs:           []models.BNPL{},
		PurchaseHistory: []models.SalesSession{},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	_, err = ctrl.customersCollection.InsertOne(context.Background(), customer)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create customer")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput([]interface{}{}, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}

	log.Debug().Str("id", customer.ID).Msg("Successfully created customer")
	return c.Status(fiber.StatusCreated).JSON(models.NewOutput([]models.Customer{customer}))
}

// UpdateCustomer godoc
// @Security BearerAuth
// @Summary Update a customer
// @Description Update an existing customer in the database
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param customer body models.CustomerBase true "Customer data"
// @Success 200 {object} models.Output
// @Failure 400 {object} models.Output
// @Failure 404 {object} models.Output
// @Failure 500 {object} models.Output
// @Router /api/customers/{id} [put]
func (ctrl *CustomersController) UpdateCustomer(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Debug().Str("id", id).Msg("Updating customer")

	var customerBase *models.CustomerBase
	if err := c.BodyParser(&customerBase); err != nil {
		log.Error().Err(err).Msg("Failed to parse customer data")
		return c.Status(fiber.StatusBadRequest).JSON(models.NewOutput([]interface{}{}, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusBadRequest,
		}))
	}

	// Check if customer exists
	var existingCustomer models.Customer
	err := ctrl.customersCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&existingCustomer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(models.NewOutput([]interface{}{}, models.Error{
				Message: "Customer not found",
				Code:    fiber.StatusNotFound,
			}))
		}
		log.Error().Err(err).Str("id", id).Msg("Failed to find customer")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput([]interface{}{}, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}

	// Check if phone number is being changed and if it conflicts with another customer
	if customerBase.Phone != existingCustomer.Phone {
		var phoneConflict models.Customer
		err := ctrl.customersCollection.FindOne(context.Background(), bson.M{"phone": customerBase.Phone}).Decode(&phoneConflict)
		if err == nil {
			return c.Status(fiber.StatusConflict).JSON(models.NewOutput([]interface{}{}, models.Error{
				Message: "Another customer with this phone number already exists",
				Code:    fiber.StatusConflict,
			}))
		}
	}

	// Update customer
	// Build dynamic update document
	updateFields := bson.M{"updated_at": time.Now()}
	if customerBase.Name != "" {
		updateFields["name"] = customerBase.Name
	}
	if customerBase.Phone != "" {
		updateFields["phone"] = customerBase.Phone
	}
	if customerBase.Address != "" {
		updateFields["address"] = customerBase.Address
	}

	// update the created time
	updateFields["created_at"] = time.Now()

	update := bson.M{
		"$set": updateFields,
	}

	result, err := ctrl.customersCollection.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to update customer")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput([]interface{}{}, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}

	if result.ModifiedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(models.NewOutput([]interface{}{}, models.Error{
			Message: "Customer not found or no changes made",
			Code:    fiber.StatusNotFound,
		}))
	}

	// Fetch updated customer
	var updatedCustomer models.Customer
	err = ctrl.customersCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&updatedCustomer)
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to fetch updated customer")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput([]interface{}{}, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}

	log.Debug().Str("id", id).Msg("Successfully updated customer")
	return c.JSON(models.NewOutput([]models.Customer{updatedCustomer}))
}

// DeleteCustomer godoc
// @Security BearerAuth
// @Summary Delete a customer
// @Description Delete a customer from the database
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} models.Output
// @Failure 404 {object} models.Output
// @Failure 500 {object} models.Output
// @Router /api/customers/{id} [delete]
func (ctrl *CustomersController) DeleteCustomer(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Debug().Str("id", id).Msg("Deleting customer")

	// Check if customer exists
	var existingCustomer models.Customer
	err := ctrl.customersCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&existingCustomer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Debug().Str("id", id).Msg("Customer not found")
			return c.Status(fiber.StatusNotFound).JSON(models.NewOutput([]interface{}{}, models.Error{
				Message: "Customer not found",
				Code:    fiber.StatusNotFound,
			}))
		}
		log.Error().Err(err).Str("id", id).Msg("Failed to find customer")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput([]interface{}{}, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}

	// Check if customer has active BNPLs
	if len(existingCustomer.BNPLs) > 0 {
		log.Debug().Str("id", id).Msg("Customer has active BNPL transactions")
		return c.Status(fiber.StatusBadRequest).JSON(models.NewOutput([]interface{}{}, models.Error{
			Message: "Cannot delete customer with active BNPL transactions",
			Code:    fiber.StatusBadRequest,
		}))
	}

	// Delete customer
	result, err := ctrl.customersCollection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		log.Error().Err(err).Str("id", id).Msg("Failed to delete customer")
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewOutput([]interface{}{}, models.Error{
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		}))
	}

	if result.DeletedCount == 0 {
		log.Debug().Str("id", id).Msg("Customer not found")
		return c.Status(fiber.StatusNotFound).JSON(models.NewOutput([]interface{}{}, models.Error{
			Message: "Customer not found",
			Code:    fiber.StatusNotFound,
		}))
	}

	log.Debug().Str("id", id).Msg("Successfully deleted customer")
	return c.JSON(models.NewOutput(map[string]interface{}{
		"message": "Customer deleted successfully",
		"id":      id,
	}))
}
