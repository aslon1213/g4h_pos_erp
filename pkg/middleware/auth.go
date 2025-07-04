package middleware

import (
	models "github.com/aslon1213/go-pos-erp/pkg/repository"
	pasetoware "github.com/gofiber/contrib/paseto"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Middlewares struct {
	UserCollection       *mongo.Collection
	ActivitiesCollection *mongo.Collection
}

func New(db *mongo.Database) *Middlewares {
	return &Middlewares{
		UserCollection:       db.Collection("users"),
		ActivitiesCollection: db.Collection("activities"),
	}
}

func (m *Middlewares) AuthMiddleware(c *fiber.Ctx) error {

	values := c.Locals(
		pasetoware.DefaultContextKey,
	).(string)

	// got username
	username := values

	// retreive the user
	user := &models.User{}
	err := m.UserCollection.FindOne(c.Context(), bson.M{"username": username}).Decode(user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to find user")
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	// add the user to the context
	c.Locals("user", user.Username)

	return c.Next()
}
