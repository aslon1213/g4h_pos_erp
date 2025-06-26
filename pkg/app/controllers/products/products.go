package products

import (
	"context"

	models "aslon1213/magazin_pos/pkg/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PorductsController struct {
	ProductsCollection *mongo.Collection
}

func NewProductsController(db *mongo.Database) *PorductsController {
	return &PorductsController{
		ProductsCollection: db.Collection("products"),
	}
}

func (c *PorductsController) CreateProduct(ctx context.Context, product *models.Product) error {
	_, err := c.ProductsCollection.InsertOne(ctx, product)
	return err
}

func (c *PorductsController) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	product := &models.Product{}
	err := c.ProductsCollection.FindOne(ctx, bson.M{"_id": id}).Decode(product)
	return product, err
}

func (c *PorductsController) QueryProducts(ctx context.Context, branchID string) ([]*models.Product, error) {
	products := []*models.Product{}
	// cursor, err := c.ProductsCollection.Find(ctx, bson.M{"branch_id": branchID})
	// if err != nil {
	// 	return nil, err
	// }
	return products, nil
}
